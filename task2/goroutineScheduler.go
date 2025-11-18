package main

import (
	"fmt"
	"sync"
	"time"
)

/*
✅Goroutine
题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
考察点 ：协程原理、并发任务调度。
*/

// 任务类型：一个无参数无返回值的函数
type Task func()

// 任务执行结果
type TaskResult struct {
	Taskname string
	Duration time.Duration
	Error    error
}

// 任务调度器
type TaskScheduler struct {
	wg        sync.WaitGroup // 每个 TaskScheduler 实例都有自己的 WaitGroup,不同的调度器实例可以独立管理自己的任务组,多个调度器实例可以同时运行，互不干扰
	result    chan TaskResult
	taskCount int
}

func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		result:    make(chan TaskResult),
		wg:        sync.WaitGroup{},
		taskCount: 0,
	}
}

// 添加并执行任务
func (s *TaskScheduler) AddTask(name string, task Task) {
	s.wg.Add(1)
	s.taskCount++

	go func(taskName string, t Task) {
		defer s.wg.Done()
		start := time.Now()

		// 执行任务（可以添加错误处理）
		defer func() {
			if r := recover(); r != nil {
				s.result <- TaskResult{
					Taskname: taskName,
					Duration: time.Since(start),
					Error:    fmt.Errorf("任务执行失败：%v", r),
				}
			}
		}()
		// 执行实际任务
		t()

		s.result <- TaskResult{
			Taskname: taskName,
			Duration: time.Since(start),
			Error:    nil,
		}
	}(name, task)

}

// 等待所有任务完成并返回结果
func (s *TaskScheduler) Wait() []TaskResult {
	go func() {
		s.wg.Wait()
		close(s.result)
	}()

	result := make([]TaskResult, 0, s.taskCount)
	for res := range s.result {
		result = append(result, res)
	}
	return result
}

func main() {
	// 创建调度器
	scheduler := NewTaskScheduler()

	// 定义一些示例任务
	tasks := map[string]Task{
		"快速任务": func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("快速任务执行完成")
		},
		"中等任务": func() {
			time.Sleep(300 * time.Millisecond)
			fmt.Println("中等任务执行完成")
		},
		"慢速任务": func() {
			time.Sleep(500 * time.Millisecond)
			fmt.Println("慢速任务执行完成")
		},
		"计算任务": func() {
			sum := 0
			for i := 0; i < 1000000; i++ {
				sum += i
			}
			fmt.Println("计算任务执行完成：%d\n", sum)
		},
	}

	// 添加所有任务到调度器
	startTime := time.Now()

	for name, task := range tasks {
		scheduler.AddTask(name, task)
	}

	// 等待所有任务完成并获取结果
	results := scheduler.Wait()
	totalTime := time.Since(startTime)

	for _, result := range results {
		if result.Error != nil {
			fmt.Printf("任务%s：失败-(耗时: %v)\n", result.Taskname, result.Error, result.Duration)
		} else {
			fmt.Printf("任务 %s: 成功 (耗时: %v)\n",
				result.Taskname, result.Duration)
		}
	}
	fmt.Printf("\n总执行时间: %v\n", totalTime)
	fmt.Printf("并发执行节省的时间: 显著!\n")
}
