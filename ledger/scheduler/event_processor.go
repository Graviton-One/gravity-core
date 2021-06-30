package scheduler

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/mitchellh/mapstructure"
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
	err = GlobalScheduler.signOraclesByNebula(Event.RoundId, nebulaId, Event.ChainType, Event.Sender)
	time.Sleep(time.Second * 5)

	success := false
	attempts := 4
	for {
		if success || attempts == 0 {
			break
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
