package hermes

import "sync"

type serviceState struct {
	sync.RWMutex
	_running bool
}

func (service *serviceState) isRunning() bool {
	service.RLock()
	defer service.RUnlock()
	return service._running
}

func (service *serviceState) setRunning(v bool) {
	service.Lock()
	service._running = v
	service.Unlock()
}
