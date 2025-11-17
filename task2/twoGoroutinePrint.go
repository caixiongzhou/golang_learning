package main

/*
✅Goroutine
题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
考察点 ： go 关键字的使用、协程的并发执行。
*/

// 1.go关键字
/*func main() {

	go func() {
		for i := 1; i <= 10; i += 2 {
			fmt.Printf("协程A,奇数:%d\n", i)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 启动偶数协程
	go func() {
		for i := 2; i <= 10; i += 2 {
			fmt.Printf("协程B,偶数:%d\n", i)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 等待协程执行完成
	time.Sleep(2 * time.Second)
	fmt.Println("打印完成！")
}*/

// 2. 使用多个协程（推荐使用 WaitGroup）
/*func main() {
	var wg sync.WaitGroup

	// 启动3个协程

	wg.Add(1) // 计数器+1
	go func() {
		defer wg.Done() // 计数器-1
		printOddNumbers(fmt.Sprintf("%s协程：", "奇数"))
	}()

	wg.Add(1) // 计数器+1
	go func() {
		defer wg.Done() // 计数器-1
		printEvenNumbers(fmt.Sprintf("%s协程：", "偶数"))
	}()

	fmt.Println("等待所有协程完成...")
	wg.Wait() // 等待所有协程完成
	fmt.Println("所有协程已完成！")
}

func printOddNumbers(name string) {
	for i := 1; i <= 10; i += 2 {
		fmt.Printf("%s: %d\n", name, i)
		time.Sleep(100 * time.Millisecond)
	}

}

func printEvenNumbers(name string) {
	for i := 2; i <= 10; i += 2 {
		fmt.Printf("%s: %d\n", name, i)
		time.Sleep(100 * time.Millisecond)
	}
}*/

// 3.使用通道协调
/*主函数
├── 创建 WaitGroup
├── 创建通道 ch
├── 启动奇数协程 → 发送奇数到通道
├── 启动偶数协程 → 发送偶数到通道
├── 启动监控协程 → 等待并关闭通道
└── 主协程循环 → 从通道接收并打印
*/
/*func main() {
	var wg sync.WaitGroup
	ch := make(chan int)

	wg.Add(2)

	// 奇数协程
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i += 2 {
			ch <- i
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 偶数协程
	go func() {
		defer wg.Done()
		for i := 2; i <= 10; i += 2 {
			ch <- i
			time.Sleep(100 * time.Millisecond)
		}
	}()
	// 监控协程（关键）
	go func() {
		wg.Wait()
		close(ch)
	}()

	//接收并打印
	for num := range ch {
		if num%2 == 0 {
			fmt.Printf("偶数协程：%d\n", num)
		} else {
			fmt.Printf("奇数协程：%d\n", num)
		}
	}
}
*/

// 4.使用带缓冲的通道，更优雅的解决方案，避免监控协程启动早于生产者
/*
func main() {
	var wg sync.WaitGroup
	ch := make(chan int, 10) // 带缓冲的通道，避免阻塞

	wg.Add(2)

	// 生产者
	go produceOdd(&wg, ch)
	go produceEven(&wg, ch)

	//监控
	go func() {
		wg.Wait()
		close(ch)
	}()
	//消费者
	comsumer(ch)

}

// 偶数生产者
func produceEven(wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()
	for i := 2; i <= 10; i += 2 {
		ch <- i
		time.Sleep(100 * time.Millisecond)
	}
}

// 奇数生产者
func produceOdd(wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()
	for i := 1; i <= 10; i += 2 {
		ch <- i
		time.Sleep(100 * time.Millisecond)
	}
}
func comsumer(ch <-chan int) {
	for num := range ch {
		if num%2 == 0 {
			fmt.Printf("偶数协程：%d\n", num)
		} else {
			fmt.Printf("奇数协程：%d\n", num)
		}
	}
}
*/
