package cron

import (
	crontab "github.com/robfig/cron"
	"time"
)

var (
	tLocal, _ = time.LoadLocation("Asia/Shanghai")
	cn        = crontab.New(crontab.WithSeconds(), crontab.WithLocation(tLocal))
)

func Add(spec string, cmd func()) {
	cn.AddFunc(spec, cmd)
	cn.Start()
}
