package pkg

import (
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
)

type Service interface {
	Name() string
}

type BackgroundService interface {
	Service
	BackgroundTask() task.Task
}
