package main

import (
	"testing"
	"time"
)

func BenchmarkNewGoroutinesPoolPool_Run(b *testing.B) {

	//设置最大线程数
	num := 100 * 20 * 100

	// 注册工作池，传入任务
	// 参数1 初始化worker(工人)并发个数 20万个
	p := NewGoroutinesPool(uint64(num))

	// dataNum := 100 * 100 * 100    //模拟百万请求
	dataNum := 100 * 100 * 100

	for i := 1; i <= dataNum; i++ {
		t := &Task{
			Params: nil,
			Handler: func(...interface{}) {
				//fmt.Println("创建一个Task:", time.Now().Format("2006-01-02 15:04:05"))
				time.Sleep(50 * time.Millisecond)
			},
		}
		p.NewTask(t) //往线程池 的通道中 写参数   每个参数相当于一个请求  来了100万个请求
	}

	p.ClosePool()
}
