package pkg

import (
	"context"
	"errors"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type SchedulerType string

const (
	SchedulerUnknown SchedulerType = ""
	SchedulerMemory  SchedulerType = "memory"
	SchedulerRedis   SchedulerType = "redis"
	SchedulerKafka   SchedulerType = "kafka"
)

type Scheduler interface {
	GetDownloader() Downloader
	SetDownloader(Downloader)
	GetExporter() Exporter
	SetExporter(Exporter)
	SetItemDelay(time.Duration)
	SetItemConcurrency(int)
	SetRequestRate(string, time.Duration, int)
	YieldItem(Context, Item) error
	MustYieldItem(Context, Item)
	Request(Context, Request) (Response, error)
	YieldRequest(Context, Request) error
	MustYieldRequest(Context, Request)
	YieldExtra(any) error
	MustYieldExtra(any)
	GetExtra(any) error
	MustGetExtra(any)
	StartScheduler(context.Context) error
	StopScheduler(context.Context) error
	Spider() Spider
	SetSpider(spider Spider)
	Interval() time.Duration
	SetInterval(time.Duration)

	SetScheduler(Scheduler) Scheduler
	SetLogger(Logger) Scheduler

	Downloader
	Exporter
}

type UnimplementedScheduler struct {
	scheduler          Scheduler
	logger             Logger
	itemConcurrency    int
	itemConcurrencyNew int
	itemDelay          time.Duration
	requestSlots       sync.Map

	Downloader
	Exporter
	middlewares []Middleware

	spider Spider
}

func (s *UnimplementedScheduler) Spider() Spider {
	return s.spider
}
func (s *UnimplementedScheduler) SetSpider(spider Spider) {
	s.spider = spider
}
func (s *UnimplementedScheduler) GetDownloader() Downloader {
	return s.Downloader
}
func (s *UnimplementedScheduler) SetDownloader(downloader Downloader) {
	s.Downloader = downloader
}
func (s *UnimplementedScheduler) GetExporter() Exporter {
	return s.Exporter
}
func (s *UnimplementedScheduler) SetExporter(exporter Exporter) {
	s.Exporter = exporter
}
func (s *UnimplementedScheduler) GetMiddlewares() []Middleware {
	return s.middlewares
}
func (s *UnimplementedScheduler) SetMiddlewares(middlewares []Middleware) {
	s.middlewares = middlewares
}
func (s *UnimplementedScheduler) GetItemDelay() time.Duration {
	return s.itemDelay
}
func (s *UnimplementedScheduler) SetItemDelay(itemDelay time.Duration) {
	s.itemDelay = itemDelay
}
func (s *UnimplementedScheduler) ItemConcurrencyNew() int {
	return s.itemConcurrencyNew
}
func (s *UnimplementedScheduler) SetItemConcurrencyNew(itemConcurrency int) {
	s.itemConcurrencyNew = itemConcurrency
}
func (s *UnimplementedScheduler) ItemConcurrency() int {
	return s.itemConcurrency
}
func (s *UnimplementedScheduler) SetItemConcurrencyRaw(itemConcurrency int) {
	s.itemConcurrency = itemConcurrency
}
func (s *UnimplementedScheduler) SetItemConcurrency(itemConcurrency int) {
	if s.itemConcurrency == itemConcurrency {
		return
	}

	if itemConcurrency < 1 {
		itemConcurrency = 1
	}

	s.itemConcurrencyNew = itemConcurrency
}
func (s *UnimplementedScheduler) RequestSlotLoad(slot string) (value any, ok bool) {
	return s.requestSlots.Load(slot)
}
func (s *UnimplementedScheduler) RequestSlotStore(slot string, value any) {
	s.requestSlots.Store(slot, value)
}
func (s *UnimplementedScheduler) SetRequestRate(slot string, interval time.Duration, concurrency int) {
	if slot == "" {
		slot = "*"
	}

	if concurrency < 1 {
		concurrency = 1
	}

	slotValue, ok := s.requestSlots.Load(slot)
	if !ok {
		requestSlot := rate.NewLimiter(rate.Every(interval/time.Duration(concurrency)), concurrency)
		s.requestSlots.Store(slot, requestSlot)
		return
	}

	limiter := slotValue.(*rate.Limiter)
	limiter.SetBurst(concurrency)
	limiter.SetLimit(rate.Every(interval / time.Duration(concurrency)))

	return
}

func (s *UnimplementedScheduler) SetScheduler(scheduler Scheduler) Scheduler {
	s.scheduler = scheduler
	return s.scheduler
}
func (s *UnimplementedScheduler) SetLogger(logger Logger) Scheduler {
	s.logger = logger
	return s.scheduler
}
func (s *UnimplementedScheduler) MustYieldItem(c Context, item Item) {
	if err := s.scheduler.YieldItem(c, item); err != nil {
		s.logger.Error(err)
	}
}
func (s *UnimplementedScheduler) MustYieldRequest(c Context, request Request) {
	if err := s.scheduler.YieldRequest(c, request); err != nil {
		s.logger.Error(err)
	}
}
func (s *UnimplementedScheduler) MustYieldExtra(extra any) {
	if err := s.scheduler.YieldExtra(extra); err != nil {
		s.logger.Error(err)
	}
}
func (s *UnimplementedScheduler) MustGetExtra(extra any) {
	if err := s.scheduler.GetExtra(extra); err != nil {
		s.logger.Error(err)
		if errors.Is(err, ErrQueueTimeout) {
			panic(ErrQueueTimeout)
		}
	}
}
func (s *UnimplementedScheduler) HandleError(ctx Context, response Response, err error, errBackName string) {
	spider := s.Spider()

	processed := false
	for _, v := range s.middlewares {
		next := v.ProcessError(ctx, response, err)
		if !next {
			break
		}
		processed = true
	}
	if processed {
		s.logger.Debug("error processed")
	}
	spider.ErrBack(errBackName)(ctx, response, err)
	spider.IncRequestError()
}
