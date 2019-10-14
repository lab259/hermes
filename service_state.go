package hermes

import "sync"

type serviceState struct {
	mutex sync.RWMutex
	_running bool
}

func (service *serviceState) isRunning() bool {
	service.mutex.RLock()
	defer service.RUnlock()
	return service._running
}

func (service *serviceState) setRunning(v bool) {
	service.mutex.Lock()
	service._running = v
	service.mutex.Unlock()
}
