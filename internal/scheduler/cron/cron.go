package cron

import "time"

type Cron struct {
	interval time.Duration
	handler  func()
}

func NewCron(
	interval time.Duration,
	handler func(),
) *Cron {
	return &Cron{
		interval: interval,
		handler:  handler,
	}
}

func (c *Cron) Start() {
	ticker := time.NewTicker(c.interval)

	go func() {

		for {
			select {
			case <-ticker.C:
				c.handler()
			}
		}
	}()
}
