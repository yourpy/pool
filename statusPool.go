package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// 协程池状态
const (
	RUNNING uint64 = iota
	STOP
)

type Task struct {
	Params  []interface{}
	Handler func(prams ...interface{})
}

type GoroutinesPool struct {
	capacity uint64
	// 目前的协程数
	workerNum uint64
	// 状态
	status uint64

	//任务池
	taskChan chan *Task

	// 锁
	sync.Mutex
}

func NewGoroutinesPool(capacity uint64) *GoroutinesPool {
	return &GoroutinesPool{
		capacity:  capacity,
		workerNum: 0,
		status:    0,
		taskChan:  make(chan *Task),
		Mutex:     sync.Mutex{},
	}
}

// 获取协程池容量，固定参数，直接返回
func (p *GoroutinesPool) getPoolCapacity() uint64 {
	return p.capacity
}

// 返回协程池内运行的协程数量
func (p *GoroutinesPool) getPoolWorkerNum() uint64 {
	return atomic.LoadUint64(&p.workerNum)
}

// 增加运行的协程数量
func (p *GoroutinesPool) addNewWorker() {
	atomic.AddUint64(&p.workerNum, 1)
}

// 减少运行的协程数量
func (p *GoroutinesPool) decWorker() {
	atomic.AddUint64(&p.workerNum, ^uint64(0))
}

// 获取协程池运行的状态
func (p *GoroutinesPool) getPoolStatus() uint64 {
	return atomic.LoadUint64(&p.status)
}

// 设置协程池状态
func (p *GoroutinesPool) setPoolStatus(status uint64) bool {
	// 需要加锁
	p.Lock()
	defer p.Unlock()
	// 如果状态一致，无需修改
	if status == p.status {
		return false
	} else {
		p.status = status
		return true
	}
}

// NewTask 启动任务
// 启动一个新的任务需要先获取协程池的锁，获取到后判断协程池的状态，
// 若没有处于关闭状态，并且容量充足，则启动一个新的协程并向任务队列中写入新的任务。
func (p *GoroutinesPool) NewTask(task *Task) error {
	p.Lock()
	defer p.Unlock()

	if p.getPoolStatus() == RUNNING {
		if p.getPoolCapacity() > p.getPoolWorkerNum() {
			p.newTaskGoroutine()
		}
		p.taskChan <- task
		return nil
	} else {
		return errors.New("goroutines pool has already closed")
	}
}

// 启动 worker
func (p *GoroutinesPool) newTaskGoroutine() {
	p.addNewWorker()

	go func() {
		defer func() {
			p.decWorker()
			if r := recover(); r != nil {
				log.Printf("worker %s has panic\n", r)
			}
		}()

		for {
			select {
			case task, ok := <-p.taskChan:
				if !ok {
					return
				}
				task.Handler(task.Params...)
			}
		}
	}()
}

func (p *GoroutinesPool) ClosePool() {
	p.setPoolStatus(STOP)
	close(p.taskChan)
	for len(p.taskChan) > 0 {
		time.Sleep(time.Second * 60)
	}
}

// ---------------for test----------------
func add(nums ...interface{}) {
	sum := 0
	for i := 0; i < len(nums); i++ {
		sum += nums[i].(int)
	}
	fmt.Println(sum)
}

//func main() {
//	nums := []int{1, 2, 3, 4, 5}
//	p := NewGoroutinesPool(10)
//
//	ta := &Task{
//		Params:  []interface{}{nums},
//		Handler: add,
//	}
//	p.NewTask(ta)
//}
