package hystrix

import (
	"context"
	"errors"
	"fmt"
	pc "github.com/blademainer/commons/pkg/panic"
	"github.com/blademainer/go-hystrix/pkg/logger"
	"time"
)

type Pool struct {
	Size        int           // 执行池大小
	Timeout     time.Duration // 执行过期时间
	commandChan chan contextAndCommand
	doneChan    chan bool
}

type contextAndCommand struct {
	context context.Context
	ICommand
}

func InitPool(poolSize int, timeoutMillSeconds int) *Pool {
	if poolSize <= 0 {
		panic(fmt.Sprintf("Illegal pool size: %d \n", poolSize))
	}
	pool := &Pool{}
	pool.Size = poolSize
	pool.Timeout = time.Duration(timeoutMillSeconds) * time.Millisecond
	pool.commandChan = make(chan contextAndCommand, poolSize)
	pool.doneChan = make(chan bool, poolSize)
	for i := 0; i < poolSize; i++ {
		pool.doneChan <- true

	}
	pool.start();
	return pool
}

func (pool *Pool) start() {
	go pool.run()
}

func (pool *Pool) run() {
	for {
		select {
		case cmd := <-pool.commandChan:
			go pc.WithRecover(func() {
				pool.invoke(cmd)
			})
		}
	}
}

func (pool *Pool) invoke(cmd contextAndCommand) {
	timeout, cancelFunc := context.WithTimeout(cmd.context, pool.Timeout)
	defer func() {
		pool.doneChan <- true
	}()
	defer cancelFunc()
	cmd.InvokeWithTimeout(timeout)
}

func (pool *Pool) Submit(context context.Context, cmd ICommand) error {
	select {
	//case pool.commandChan <- cmd:
	//	if logger.Log.IsDebugEnabled() {
	//		logger.Log.Debugf("Submitted cmd: %v", cmd)
	//	}
	case <-pool.doneChan:
		pool.commandChan <- contextAndCommand{context: context, ICommand: cmd}
		if logger.Log.IsDebugEnabled() {
			logger.Log.Debugf("Submitted cmd: %v", cmd)
		}
		return nil
	default:
		e := errors.New("pool is full")
		cmd.Fallback(context, e.Error(), e)
		return e
	}
}
