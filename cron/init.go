package cron

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/VincentTC/simple-ecommerce/model"
	"github.com/VincentTC/simple-ecommerce/service"
	"github.com/robfig/cron/v3"
)

// Options for the API.
type Options struct {
	Sv           service.Sv
	CronSchedule CronSchedule
	Apps         model.AppsConfig
}

type CronSchedule struct {
	OrderReminderSchedule string
}

type Cron struct {
	cron    *cron.Cron
	options *Options
}

func NewCron(options *Options) *Cron {
	return &Cron{
		cron: cron.New(cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
		)),
		options: options,
	}
}

func (c *Cron) Init() {
	// register cron order reminder
	c.cron.AddFunc(c.options.CronSchedule.OrderReminderSchedule, orderReminder(c).Handle)
}

func (c *Cron) Run() {
	c.cron.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigChan
}

func (c *Cron) Stop() {
	c.cron.Stop()
}
