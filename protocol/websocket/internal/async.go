package internal

import (
	"context"
	"fmt"
	uniformApi "github.com/zzy-rabbit/bp/tool/uniform/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"sync"
	"time"
)

type wait struct {
	timeout  time.Duration
	start    time.Time
	errors   chan xerror.IError
	response uniformApi.Frame
}

type async struct {
	mutex sync.RWMutex
	waits map[uint64]*wait
}

func NewSync(ctx context.Context) *async {
	s := &async{
		waits: make(map[uint64]*wait),
	}
	s.startMonitor(ctx)
	return s
}

func (s *async) startMonitor(ctx context.Context) {
	go func() {
		for {
			time.Sleep(time.Second * 1)
			s.mutex.Lock()
			for seq, wait := range s.waits {
				if time.Now().Sub(wait.start) > wait.timeout {
					wait.errors <- xerror.Extend(xerror.ErrTimeout, fmt.Sprintf("timeout: %s", wait.timeout))
					delete(s.waits, seq)
				}
			}
			s.mutex.Unlock()
		}
	}()
}

func (s *async) receive(ctx context.Context, frame uniformApi.Frame) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	wait, ok := s.waits[frame.Sequence]
	if ok {
		wait.response = frame
		wait.errors <- nil
		delete(s.waits, frame.Sequence)
	}
	return ok
}

func (s *async) wait(ctx context.Context, frame uniformApi.Frame, timeout time.Duration) *wait {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if timeout <= 0 {
		timeout = time.Second * 3
	}
	wait := &wait{
		timeout: timeout,
		start:   time.Now(),
		errors:  make(chan xerror.IError),
	}
	s.waits[frame.Sequence] = wait
	return wait
}
