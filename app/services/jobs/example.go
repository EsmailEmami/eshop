package jobs

import (
	"github.com/robfig/cron/v3"
)

func ExampleJob() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 * * * * *" /*every minute*/, func() {

	})
	c.Start()
}
