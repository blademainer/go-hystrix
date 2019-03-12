package main

import (
	"context"
	"fmt"
	"github.com/blademainer/go-hystrix/pkg/hystrix"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type FakeCmd struct {
	string
}

var counter uint64

func (*FakeCmd) InvokeWithTimeout(context context.Context) error {
	time.Sleep(50 * time.Millisecond)
	atomic.AddUint64(&counter, 1)
	wg.Done()
	return nil
}

func (*FakeCmd) Fallback(message string, err error) {
	fmt.Println(message)
}

var wg = sync.WaitGroup{}

func main() {
	pool := hystrix.InitPool(100, 50)
	for i := 0; i < 1000; i++ {
		go func() {
			duration := rand.Intn(100)
			time.Sleep(time.Duration(duration) * time.Millisecond)
			f := &FakeCmd{}
			e := pool.Submit(f)
			if e == nil{
				wg.Add(1)
			}
		}()
	}
	wg.Wait()
	fmt.Println("counter: ", counter)
}