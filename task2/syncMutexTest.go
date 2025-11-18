package main

import (
	"fmt"
	"sync"
)

/*✅锁机制
题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ： sync.Mutex 的使用、并发数据安全。*/

// Counter 结构体包含一个共享计数器和互斥锁
type Counter struct {
	Value int
	mutex sync.Mutex
}

// Increment 使用互斥锁安全地递增计数器
func (c *Counter) Increment() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Value++
}

// Value 使用互斥锁安全地获取计数器值
func (c *Counter) GetValue() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Value
}

func main() {

	//var mutexInstance sync.Mutex
	// 创建计数器实例
	counter := &Counter{
		Value: 0,
		//mutex: mutexInstance,  写不写都一样，上面的Counter结构体已经借用外面包里面sync.Mutex
	}

	// 使用 WaitGroup 等待所有协程完成
	var wg sync.WaitGroup

	// 启动10个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// 每个协程对计数器进行1000次递增操作
			for j := 0; j < 1000; j++ {
				counter.Increment()
			}
			fmt.Printf("协程 %d 完成工作\n", id)
		}(i)
	}

	// 等待所有协程完成
	wg.Wait()

	// 输出最终的计数器值
	fmt.Printf("最终计数器：%d\n", counter.GetValue())
	fmt.Printf("期望值: %d\n", 10*1000)
}
