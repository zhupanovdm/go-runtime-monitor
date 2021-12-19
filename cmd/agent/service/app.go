package service

import (
	"net/url"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/internal/app"
	"github.com/zhupanovdm/go-runtime-monitor/internal/encoder"
)

var PollInterval time.Duration
var ReportInterval time.Duration
var Url *url.URL

func StartAgent() {
	data := make(chan encoder.Encoder, 1024)

	app.Periodic(PollInterval, MetricsCollector(data))
	app.Periodic(ReportInterval, MetricsTransporter(data))

	app.Serve()
}
