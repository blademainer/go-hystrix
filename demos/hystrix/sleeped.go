package main

import (
	"context"
	"fmt"
	"github.com/blademainer/go-hystrix/pkg/hystrix"
	"sync"
	"sync/atomic"
	"time"
)

type FakeCmd struct {
	string
}

var counter uint64

const THREAD_NAME_CONTEXT_NAME = "THREAD_NAME"

func (*FakeCmd) InvokeWithTimeout(context context.Context) error {
	fmt.Printf("Executing %v... \n", context.Value(THREAD_NAME_CONTEXT_NAME))
	time.Sleep(1000 * time.Millisecond)
	atomic.AddUint64(&counter, 1)
	wg.Done()
	fmt.Printf("Done %v... \n", context.Value(THREAD_NAME_CONTEXT_NAME))
	return nil
}

func (*FakeCmd) Fallback(context context.Context, message string, err error) {
	fmt.Printf("Fallback %v... message: %v error: %v \n", context.Value(THREAD_NAME_CONTEXT_NAME), message, err.Error())
}

var wg = sync.WaitGroup{}

func main() {
	pool := hystrix.InitPool(100, 50)
	for i := 0; i < 1000; i++ {
		time.Sleep(1 * time.Millisecond)
		go func() {
			wg.Add(1)
			c := context.WithValue(context.TODO(), THREAD_NAME_CONTEXT_NAME, fmt.Sprint(i))
			//duration := rand.Intn(100)
			f := &FakeCmd{}
			e := pool.Submit(c, f)
			if e != nil {
				fmt.Println("e: ", e.Error())
				wg.Done()
			}
		}()
	}
	wg.Wait()
	fmt.Println("counter: ", counter)
}
