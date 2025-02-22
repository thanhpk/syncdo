package syncdo

import (
	"fmt"
	"sync"
)

var ll = &sync.Mutex{}
var klm = make(map[string]*kmutex)
var kls = make([]*kmutex, 0)

func KLock(k string) func() {
	ll.Lock()
	kl, has := klm[k]
	if !has {
		kl = usekl()
		klm[k] = kl
	}

	kl.num++
	ll.Unlock()
	kl.Lock()
	return func() {
		ll.Lock()
		kl.num--
		if kl.num == 0 {
			kl.num = 0
			delete(klm, k)
		}
		ll.Unlock()
		kl.Unlock()
	}
}

func Status() string {
	ll.Lock()
	m := len(klm)
	n := len(kls)
	ll.Unlock()
	return fmt.Sprintf("locking %d/%d", m, n)
}

func usekl() *kmutex {
	for _, kl := range kls {
		if kl.num == 0 {
			return kl
		}
	}
	kl := &kmutex{Mutex: &sync.Mutex{}}
	kls = append(kls, kl)
	return kl
}

type kmutex struct {
	*sync.Mutex
	num int // number of waiters for the lock
}
