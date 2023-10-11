package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lizongying/go-crawler/pkg"
	context2 "github.com/lizongying/go-crawler/pkg/context"
	request2 "github.com/lizongying/go-crawler/pkg/request"
	"github.com/segmentio/kafka-go"
	"golang.org/x/time/rate"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func (s *Scheduler) Request(ctx pkg.Context, request pkg.Request) (response pkg.Response, err error) {
	defer func() {
		s.Spider().StateRequest().Set()
	}()

	if request == nil {
		err = errors.New("nil request")
		return
	}

	s.logger.Debugf("request: %+v", request)

	response, err = s.Download(ctx, request)
	if err != nil {
		if errors.Is(err, pkg.ErrIgnoreRequest) {
			s.logger.Info(err)
			err = nil
			return
		}

		s.HandleError(ctx, response, err, request.ErrBack())
		return
	}

	s.logger.Debugf("request %+v", request)
	return
}

func (s *Scheduler) handleRequest(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	slot := "*"
	value, _ := s.RequestSlotLoad(slot)
	requestSlot := value.(*rate.Limiter)

	for {
		req, err := s.kafkaReader.FetchMessage(ctx)
		if err != nil {
			s.logger.Warn(err)
			continue
		}
		if len(req.Value) == 0 {
			err = errors.New("req is empty")
			s.logger.Warn(err)
			continue
		}

		s.logger.Debugf("request: %s", req)
		var requestJsonWithContext request2.JsonWithContext
		requestJsonWithContext.ContextJson = new(context2.Json)
		requestJsonWithContext.RequestJson = new(request2.Json)
		err = json.Unmarshal(req.Value, &requestJsonWithContext)
		if err != nil {
			s.logger.Warn(err)
			continue
		}

		request, err := requestJsonWithContext.ToRequest()
		c := requestJsonWithContext.ToContext()
		s.logger.Debugf("request: %+v", request)
		if err != nil {
			s.logger.Warn(err)
			continue
		}
		slot = request.Slot()
		if slot == "" {
			slot = "*"
		}
		slotValue, ok := s.RequestSlotLoad(slot)
		if !ok {
			concurrency := uint8(1)
			if request.Concurrency() != nil {
				concurrency = *request.Concurrency()
			}
			if concurrency < 1 {
				concurrency = 1
			}
			requestSlot = rate.NewLimiter(rate.Every(request.Interval()/time.Duration(concurrency)), int(concurrency))
			s.RequestSlotStore(slot, requestSlot)
		}

		requestSlot = slotValue.(*rate.Limiter)

		err = requestSlot.Wait(ctx)
		if err != nil {
			s.logger.Error(err)
		}
		go func(c pkg.Context, request pkg.Request) {
			response, e := s.Request(c, request)
			if e != nil {
				s.Spider().StateRequest().Out()
				return
			}

			go func(ctx pkg.Context, response pkg.Response) {
				defer func() {
					if r := recover(); r != nil {
						buf := make([]byte, 1<<16)
						runtime.Stack(buf, true)
						err = errors.New(string(buf))
						s.logger.Error(err)
						s.HandleError(ctx, response, err, request.ErrBack())
					}
				}()

				s.Spider().StateMethod().In()
				if err = s.Spider().CallBack(request.CallBack())(ctx, response); err != nil {
					s.logger.Error(err)
					s.HandleError(ctx, response, err, request.ErrBack())
				}
				s.Spider().StateMethod().Out()
				s.Spider().StateRequest().Out()
			}(c, response)
		}(c, request)
	}

	return
}

func (s *Scheduler) YieldRequest(c pkg.Context, request pkg.Request) (err error) {
	defer func() {
		s.Spider().StateRequest().Set()
	}()

	meta := c.Meta()

	// add referrer to request
	if meta.Referrer != nil {
		request.SetReferrer(meta.Referrer.String())
	}

	// add cookies to request
	if len(meta.Cookies) > 0 {
		for _, cookie := range meta.Cookies {
			request.AddCookie(cookie)
		}
	}

	s.Spider().StateRequest().In()
	bs, err := (&request2.WithContext{
		Context: c,
		Request: request,
	}).MarshalWithContext()
	s.logger.Info("request with context:", string(bs))
	if err != nil {
		s.logger.Error(err)
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err = s.kafkaWriter.WriteMessages(ctx, kafka.Message{
		Value: bs,
	}); err != nil {
		s.logger.Error(err)
		return
	}

	return
}

func (s *Scheduler) YieldExtra(extra any) (err error) {
	defer func() {
		s.Spider().StateRequest().In()
		s.Spider().StateRequest().Set()
	}()

	extraValue := reflect.ValueOf(extra)
	if extraValue.Kind() != reflect.Ptr || extraValue.IsNil() {
		err = errors.New("extra must be a non-null pointer")
		return
	}

	name := extraValue.Elem().Type().Name()

	bs, err := json.Marshal(extra)
	if err != nil {
		s.logger.Error(err)
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	kafkaWriter := &kafka.Writer{
		Addr:                   kafka.TCP(strings.Split(s.config.KafkaUri(), ",")...),
		AllowAutoTopicCreation: true,
		Topic:                  fmt.Sprintf("%s-%s-extra-%s", s.config.GetBotName(), s.Spider().Name(), name),
	}
	defer func() {
		err = kafkaWriter.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}()
	if err = kafkaWriter.WriteMessages(ctx, kafka.Message{
		Value: bs,
	}); err != nil {
		s.logger.Error(err)
		return
	}

	return
}

func (s *Scheduler) GetExtra(extra any) (err error) {
	defer func() {
		s.Spider().StateRequest().Out()
	}()

	extraValue := reflect.ValueOf(extra)
	if extraValue.Kind() != reflect.Ptr || extraValue.IsNil() {
		err = errors.New("extra must be a non-null pointer")
		return
	}

	name := extraValue.Elem().Type().Name()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.CloseReasonQueueTimeout())*time.Second)
	defer cancel()

	resultChan := make(chan struct{})
	go func() {
		var msg kafka.Message
		msg, err = kafka.NewReader(kafka.ReaderConfig{
			Brokers:  s.kafkaReader.Config().Brokers,
			MaxBytes: 10e6, // 10MB
			Topic:    fmt.Sprintf("%s-%s-extra-%s", s.config.GetBotName(), s.Spider().Name(), name),
			GroupID:  s.config.GetBotName(),
		}).FetchMessage(ctx)
		if err != nil {
			s.logger.Error(err)
			return
		}

		if len(msg.Value) == 0 {
			err = errors.New("msg error")
			s.logger.Error(err)
			return
		}

		err = json.Unmarshal(msg.Value, extra)
		if err != nil {
			s.logger.Error(err)
			return
		}

		resultChan <- struct{}{}
	}()

	select {
	case <-resultChan:
		return
	case <-ctx.Done():
		close(resultChan)
		err = pkg.ErrQueueTimeout
		return
	}
}
