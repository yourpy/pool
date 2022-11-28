package main

import (
	"fmt"
	"time"
)

/**
pool的实现：将需要执行的函数转换成一个个 task，发到协程池，接着协程池内随机一个 worker 执行
协程池里的 worker 用 Goroutine 承载
*/

// Task 每一个任务Task都可以抽象成一个函数
type Task struct {
	f func() error
}

// NewTask 创建 task
func NewTask(f func() error) *Task {
	return &Task{f: f}
}

// Execute 执行 task 任务
func (t *Task) Execute() {
	t.f()
}

type Pool struct {
	// 最大 worker 数量
	workerNum int
	// 从外部接收 task
	EntryChannel chan *Task
	// 协程池内的就绪任务队列
	JobsChannel chan *Task
}

func NewPool(cap int) *Pool {
	return &Pool{
		workerNum:    cap,
		EntryChannel: make(chan *Task),
		JobsChannel:  make(chan *Task),
	}
}

// 协程池创建 worker 并开始工作
func (p *Pool) worker(workerID int) {
	// 监听就绪任务队列，不断地从 chan 获取 task 执行
	for task := range p.JobsChannel {
		task.Execute()
		//fmt.Println("worker ID ", workerID, " 执行完毕任务")
	}
}

func (p *Pool) Run() {
	// 根据 worker 数量初始化 worker，每个 worker 用 goroutine 承载
	for i := 0; i < p.workerNum; i++ {
		//fmt.Println("开启 Worker:", i)
		go p.worker(i)
	}

	// 将 EntryChannel 的 task 写入 JobsChannel
	for task := range p.EntryChannel {
		p.JobsChannel <- task
	}

	// 关闭 JobsChannel
	close(p.JobsChannel)
	fmt.Println("执行完毕需要关闭JobsChannel")
}

func main() {

	//创建一个协程池,最大开启3个协程worker
	p := NewPool(3)

	//开一个协程 不断的向 Pool 输送打印一条时间的task任务
	go func() {
		for i := 1; i <= 10; i++ {
			t := NewTask(func() error {
				fmt.Println("创建一个Task:", time.Now().Format("2006-01-02 15:04:05"))
				return nil
			})
			p.EntryChannel <- t
		}
		close(p.EntryChannel)
	}()

	// 启动协程池p
	p.Run()
}
