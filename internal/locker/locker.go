package locker

import "sync"

var l sync.Mutex

func TryLock() bool {
	return l.TryLock()
}

func Unlock() {
	l.Unlock()
}
