package pkg

import "github.com/zhupanovdm/go-runtime-monitor/pkg/task"

type (
	// Service is used to mark type as service
	Service interface {
		Name() string
	}

	// BackgroundService is used to mark type as background service
	BackgroundService interface {
		Service
		BackgroundTask() task.Task
	}
)
