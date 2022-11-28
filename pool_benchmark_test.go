package main

import (
	"testing"
	"time"
)

func BenchmarkPool_Run(b *testing.B) {

	//设置最大线程数
	num := 20* 100 * 100

	// 注册工作池，传入任务
	// 参数1 初始化worker(工人)并发个数 20万个
	p := NewPool(num)

	// dataNum := 100 * 100 * 100    //模拟百万请求
	dataNum := 100 * 100 * 100
	go func() { //这是一个独立的协程 保证可以接受到每个用户的请求
		for i := 1; i <= dataNum; i++ {
			t := NewTask(func() error {
				//fmt.Println("创建一个Task:", time.Now().Format("2006-01-02 15:04:05"))
				time.Sleep(50 * time.Millisecond)
				return nil
			})
			p.EntryChannel <- t //往线程池 的通道中 写参数   每个参数相当于一个请求  来了100万个请求
		}
		close(p.EntryChannel)
	}()

	p.Run() //有任务就去做，没有就阻塞，任务做不过来也阻塞
}

func Benchmark_Run(b *testing.B) {

	// dataNum := 100 * 100 * 100    //模拟百万请求
	dataNum := 100 * 100 * 100

	// 控制流程
	ch := make(chan int)
	count := 0

	go func() {
		for i := 1; i <= dataNum; i++ {
			t := NewTask(func() error {
				//fmt.Println("创建一个Task:", time.Now().Format("2006-01-02 15:04:05"))
				time.Sleep(50 * time.Millisecond)
				return nil
			})
			go t.Execute()
			ch <- count
			count++
		}
		close(ch)
	}()

	exitFlag := false
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				exitFlag = true
				break
			} else {
				//fmt.Println(num)
			}
		}

		if exitFlag {
			//fmt.Println("跳出循环")
			break
		}
	}

	//fmt.Println("end--------------")
}
