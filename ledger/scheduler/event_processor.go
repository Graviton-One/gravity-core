package scheduler

import (
	"encoding/json"
	"sync"

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
