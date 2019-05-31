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

const THREAD_NAME_CONTEXT_NAME = "THREAD_NAME"

type FakeCmd struct {
	string
}

var counter uint64

func (*FakeCmd) InvokeWithTimeout(context context.Context) error {
	fmt.Printf("Executing %v... \n", context.Value(THREAD_NAME_CONTEXT_NAME))
	time.Sleep(50 * time.Millisecond)
	atomic.AddUint64(&counter, 1)
	wg.Done()
	fmt.Printf("Done %v... \n", context.Value(THREAD_NAME_CONTEXT_NAME))
	return nil
}

func (*FakeCmd) Fallback(context context.Context, message string, err error) {
	fmt.Printf("Fallback %v... message: %v error: %v \n", context.Value(THREAD_NAME_CONTEXT_NAME), message, err.Error())
}

var wg = sync.WaitGroup{}

func TestInvoke(t *testing.T) {
	pool := InitPool(100, 50)
	pool.start()
	for i := 0; i < 1000; i++ {
		go func() {
			duration := rand.Intn(100)
			time.Sleep(time.Duration(duration) * time.Millisecond)
			c := context.WithValue(context.TODO(), THREAD_NAME_CONTEXT_NAME, fmt.Sprint(i))
			f := &FakeCmd{}
			e := pool.Submit(c, f)
			if e == nil{
				wg.Add(1)
			}
		}()
	}
	wg.Wait()
	fmt.Println("counter: ", counter)
}
