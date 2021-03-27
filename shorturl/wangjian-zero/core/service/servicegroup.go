package service

import (
	"log"

	"shorturl/wangjian-zero/core/proc"
	"shorturl/wangjian-zero/core/syncx"
	"shorturl/wangjian-zero/core/threading"
)

type (
	// Starter is the interface wraps the Start method.
	Starter interface {
		Start()
	}

	// Stopper is the interface wraps the Stop method.
	Stopper interface {
		Stop()
	}

	// Service is the interface that groups Start and Stop methods.
	Service interface {
		Starter
		Stopper
	}

	// A ServiceGroup is a group of services.
	ServiceGroup struct {
		services []Service
		stopOnce func()
	}
)

// NewServiceGroup returns a ServiceGroup.
func NewServiceGroup() *ServiceGroup {
	sg := new(ServiceGroup)
	sg.stopOnce = syncx.Once(sg.doStop)
	return sg
}

// Add adds service into sg.
func (sg *ServiceGroup) Add(service Service) {
	sg.services = append(sg.services, service)
}

// Start starts the ServiceGroup.
// There should not be any logic code after calling this method, because this method is a blocking one.
// Also, quitting this method will close the logx output.
func (sg *ServiceGroup) Start() {
	proc.AddShutdownListener(func() {
		log.Println("Shutting down...")
		sg.stopOnce()
	})

	sg.doStart()
}

// Stop stops the ServiceGroup.
func (sg *ServiceGroup) Stop() {
	sg.stopOnce()
}

func (sg *ServiceGroup) doStart() {
	routineGroup := threading.NewRoutineGroup()

	for i := range sg.services {
		service := sg.services[i]
		routineGroup.RunSafe(func() {
			service.Start()
		})
	}

	routineGroup.Wait()
}

func (sg *ServiceGroup) doStop() {
	for _, service := range sg.services {
		service.Stop()
	}
}

// WithStart wraps a start func as a Service.
func WithStart(start func()) Service {
	return startOnlyService{
		start: start,
	}
}

// WithStarter wraps a Starter as a Service.
func WithStarter(start Starter) Service {
	return starterOnlyService{
		Starter: start,
	}
}

type (
	stopper struct {
	}

	startOnlyService struct {
		start func()
		stopper
	}

	starterOnlyService struct {
		Starter
		stopper
	}
)

func (s stopper) Stop() {
}

func (s startOnlyService) Start() {
	s.start()
}