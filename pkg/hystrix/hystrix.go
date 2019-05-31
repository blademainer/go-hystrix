package hystrix

import "context"

type ICommand interface {
	// 实际执行的函数
	InvokeWithTimeout(context context.Context) error

	// 当Invoke服务调用失败时回调该函数
	Fallback(context context.Context, message string, err error)
}
