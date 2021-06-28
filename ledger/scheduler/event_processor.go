package scheduler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"

	stdHttp "net/http"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-http/pkg/http"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type SchedulerEvent struct {
	Name   string
	Params map[string]interface{}
}
type EventHandler func(event SchedulerEvent) error
type EventServer struct {
	Handlers map[string]EventHandler
}

type UpdateOraclesRequest struct {
	NebulaKey string                `mapstructure:"nebula_key"`
	ChainType account.ChainType     `mapstructure:"chain_type"`
	Sender    account.OraclesPubKey `mapstructure:"sender"`
	RoundId   int64                 `mapstructure:"round_id"`
	IsSender  bool                  `mapstructure:"is_sender"`
}

func PublishMessage(topic string, rawmsg interface{}) {
	tosend, err := json.Marshal(rawmsg)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	msg := message.NewMessage(watermill.NewUUID(), tosend)
	if err := EventBus.Publish(topic, msg); err != nil {
		zap.L().Error(err.Error())
	}
}

func NewEventServer() *EventServer {
	server := &EventServer{
		Handlers: map[string]EventHandler{},
	}
	server.SetHandler("handle_block", HandleBlock)
	server.SetHandler("update_oracles", UpdateOraclesHandler)
	return server
}
func (s *EventServer) SetHandler(name string, handler EventHandler) {
	s.Handlers[name] = handler
}
func (s *EventServer) handleEvent(event SchedulerEvent) {
	handler, ok := s.Handlers[event.Name]
	if ok {
		err := handler(event)
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
	}
}
func (s *EventServer) Serve(messages <-chan *message.Message) {
	go startHttpListener()
	for msg := range messages {
		var wg sync.WaitGroup
		lmsg := *msg
		wg.Add(1)
		go func(mwg *sync.WaitGroup, m message.Message) {
			zap.L().Sugar().Debug("Receive message", string(m.Payload))
			defer wg.Done()
			event := SchedulerEvent{}
			err := json.Unmarshal(msg.Payload, &event)
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
			s.handleEvent(event)
		}(&wg, lmsg)

		msg.Ack()
	}
}

func HandleBlock(event SchedulerEvent) error {
	customEvent := struct {
		Height int64 `mapstructure:"height"`
	}{}
	err := mapstructure.Decode(event.Params, &customEvent)
	if err != nil {
		return err
	}
	zap.L().Sugar().Debug("Calling processor", customEvent)
	GlobalScheduler.Process(int64(customEvent.Height))
	return nil
}

func UpdateOraclesHandler(event SchedulerEvent) error {
	Event := UpdateOraclesRequest{}
	err := mapstructure.Decode(event.Params, &Event)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	nebulaId, err := account.StringToNebulaId(Event.NebulaKey, Event.ChainType)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	success := false
	attempts := 4
	for {
		if success || attempts == 0 {
			break
		}
		err = GlobalScheduler.signOraclesByNebula(Event.RoundId, nebulaId, Event.ChainType, Event.Sender)
		time.Sleep(time.Second * 5)
		if err != nil {
			zap.L().Error(err.Error())
			attempts -= 1
			continue
		}

		if Event.IsSender {
			err = GlobalScheduler.sendOraclesToNebula(nebulaId, Event.ChainType, Event.RoundId)
			if err != nil {
				attempts -= 1
				continue
			}
		}
		success = true
	}

	//GlobalScheduler.Process(int64(customEvent.Height))
	return nil
}

func startHttpListener() {
	logger := watermill.NewStdLogger(true, true)
	//channelPublisher, err := EventBus.Subscribe(EventBus, "ledger.events")
	httpSubscriber, err := http.NewSubscriber(
		"127.0.0.1:3501",
		http.SubscriberConfig{
			UnmarshalMessageFunc: func(topic string, request *stdHttp.Request) (*message.Message, error) {
				b, err := ioutil.ReadAll(request.Body)
				if err != nil {
					return nil, errors.Wrap(err, "cannot read body")
				}

				return message.NewMessage(watermill.NewUUID(), b), nil
			},
		},
		logger,
	)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}

	r, err := message.NewRouter(
		message.RouterConfig{},
		logger,
	)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}

	r.AddMiddleware(
		middleware.Recoverer,
		middleware.CorrelationID,
	)
	r.AddPlugin(plugin.SignalsHandler)

	r.AddHandler(
		"http_to_ledger",
		"/webhooks", // this is the URL of our API
		httpSubscriber,
		"ledger.events", // this is the topic the message will be published to
		EventBus,
		func(msg *message.Message) ([]*message.Message, error) {
			webhook := SchedulerEvent{}

			if err := json.Unmarshal(msg.Payload, &webhook); err != nil {
				return nil, errors.Wrap(err, "cannot unmarshal message")
			}

			// Simply forward the message from HTTP Subscriber to Kafka Publisher
			return []*message.Message{msg}, nil
		},
	)
	// HTTP server needs to be started after router is ready.
	go func() {
		log.Println("waiting for router ready")
		<-r.Running()
		log.Println("starting http server")
		err = httpSubscriber.StartHTTPServer()
		if err != nil {
			zap.L().Error(err.Error())
		}
	}()

	err = r.Run(context.Background())
	if err != nil {
		zap.L().Error(err.Error())
	}
}
