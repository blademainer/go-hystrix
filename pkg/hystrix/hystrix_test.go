package hystrix

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
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

func TestInvoke(t *testing.T) {
	pool := InitPool(100, 50)
	pool.start()
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
