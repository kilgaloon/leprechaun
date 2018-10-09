package agent

import "sync"

// Vault interface provides lock/unlock
type Vault interface {
	GetMutex() *sync.Mutex
	Lock()
	Unlock()
}

// Lock wrapps around mutex lock and is alias {agent}.Agent.GetMutex().Lock
func (d Default) Lock() {
	d.GetMutex().Lock()
}

// Unlock wrapps around mutex unlock and is alias {agent}.Agent.GetMutex().Unlock
func (d Default) Unlock() {
	d.GetMutex().Unlock()
}
