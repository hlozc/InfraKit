package middleware

/*
context 的使用，主要分为两类
1. 一个是用来存储上下文，比如一个请求的处理过程，我们希望一个结构体可以存储这里面存储的所有数据

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	process(ctx)

	ctx = context.WithValue(ctx, "traceId", "qcrao-2019")
	process(ctx)
}

func process(ctx context.Context) {
	traceId, ok := ctx.Value("traceId").(string)
	if ok {
		fmt.Printf("process over. trace_id=%s\n", traceId)
	} else {
		fmt.Printf("process over. no trace_id\n")
	}
}


2. 一个协程，我们希望在合适的时机就及时销毁这个协程，那么 context 内部有一个 Done()，会返回一个只读管道，合适的时机会发生读事件，比如 timeout

func Perform(ctx context.Context) {
    for {
        calculatePos()
        sendResult()

        select {
        case <-ctx.Done():
            // 被取消，直接返回
            return
        case <-time.After(time.Second):
            // block 1 秒钟
        }
    }
}

而主流程可能是这样的：


ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
go Perform(ctx)

// ……
// app 端返回页面，调用cancel 函数
cancel()

3. 接口和使用说明

ctx, cancel := context.WithTimeout(parent, timeout)
会创建一个子上下文，并设置超时时间
如果超时时间到了，那么子上下文就会自动触发取消信号，这里会通过 Done() 这个接口所反馈
cancel() 可以主动取消

*/
