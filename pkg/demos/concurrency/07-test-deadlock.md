我将创建一个简单的死锁示例。死锁是指两个或多个协程互相等待对方持有的资源，导致所有相关协程都无法继续执行的情况。我将通过两个协程互相等待对方持有的锁来演示这个问题。

让我先搜索一下当前代码库中是否已经有相关实现：

Ran tool
我将创建一个新文件来演示死锁的情况。我会创建一个简单的场景，其中两个协程互相等待对方持有的锁，从而导致死锁。

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// 演示死锁的简单示例
func testDeadlock() {
	fmt.Println("=== 测试死锁 ===")

	// 创建两个互斥锁
	lock1 := &sync.Mutex{}
	lock2 := &sync.Mutex{}

	// 创建等待组
	var wg sync.WaitGroup
	wg.Add(2)

	// 协程1：先获取lock1，然后尝试获取lock2
	go func() {
		defer wg.Done()

		fmt.Println("协程1: 尝试获取lock1")
		lock1.Lock()
		fmt.Println("协程1: 已获取lock1")

		// 短暂等待，确保死锁发生
		time.Sleep(time.Millisecond * 100)

		fmt.Println("协程1: 尝试获取lock2") // 这里会发生死锁
		lock2.Lock()
		fmt.Println("协程1: 已获取lock2")

		// 释放锁
		lock2.Unlock()
		lock1.Unlock()
	}()

	// 协程2：先获取lock2，然后尝试获取lock1
	go func() {
		defer wg.Done()

		fmt.Println("协程2: 尝试获取lock2")
		lock2.Lock()
		fmt.Println("协程2: 已获取lock2")

		// 短暂等待，确保死锁发生
		time.Sleep(time.Millisecond * 100)

		fmt.Println("协程2: 尝试获取lock1") // 这里会发生死锁
		lock1.Lock()
		fmt.Println("协程2: 已获取lock1")

		// 释放锁
		lock1.Unlock()
		lock2.Unlock()
	}()

	// 等待协程完成（注意：由于死锁，这里会一直等待）
	fmt.Println("等待协程完成...")
	wg.Wait()
	fmt.Println("程序结束") // 这一行不会被执行到，因为发生了死锁
}
```
现在我需要在 main.go 中添加对这个新函数的调用：

```go:features/concurrency/main.go
package main

// 并发编程

// 1. 协程
// 2. 通道
// 3. 互斥锁
// 4. 读写锁
// 5. 条件变量
// 6. 信号量
// 7. 死锁

func main() {
	// 1.测试协程
	// testCreateGoroutine()

	// 2.测试通道
	// testCreateChannel()

	// 3.测试互斥锁
	// testMutex()

	// 4.测试读写锁
	// testReadWriteLock()

	// 5.测试条件变量
	// testCond()

	// 7.测试死锁
	testDeadlock()
}
```
我已经创建了一个简单的死锁示例。这个示例展示了死锁最经典的情况之一：两个协程互相等待对方持有的资源。让我解释一下代码的关键部分：

1. **死锁场景设置**：
   - 创建了两个互斥锁：`lock1` 和 `lock2`
   - 创建了两个协程，每个协程都试图以不同的顺序获取这两个锁

2. **死锁形成过程**：
   ```go
   // 协程1的加锁顺序：lock1 -> lock2
   go func() {
       lock1.Lock()
       time.Sleep(time.Millisecond * 100)
       lock2.Lock() // 死锁点
   }()

   // 协程2的加锁顺序：lock2 -> lock1
   go func() {
       lock2.Lock()
       time.Sleep(time.Millisecond * 100)
       lock1.Lock() // 死锁点
   }()
   ```

3. **死锁发生的原因**：
   - 协程1先获取了lock1，然后等待lock2
   - 协程2先获取了lock2，然后等待lock1
   - 两个协程都在等待对方释放锁，但都不会释放自己持有的锁
   - 结果是两个协程永远相互等待，形成死锁

4. **运行结果**：
   当你运行这个程序时，你会看到类似这样的输出：
   ```
   === 测试死锁 ===
   等待协程完成...
   协程1: 尝试获取lock1
   协程1: 已获取lock1
   协程2: 尝试获取lock2
   协程2: 已获取lock2
   协程1: 尝试获取lock2
   协程2: 尝试获取lock1
   ```
   之后程序会一直挂起，因为发生了死锁。

5. **如何避免这种死锁**：
   - 保持固定的锁获取顺序：总是按照相同的顺序获取锁
   - 使用超时机制：在获取锁时设置超时时间
   - 使用 `tryLock` 机制：如果无法立即获取锁就放弃
   - 使用更高级的同步原语：如通道或条件变量

要运行这个示例，你可以直接执行：
```bash
cd features/concurrency && go run .
```

这个示例虽然简单，但很好地展示了死锁的本质：多个协程因为资源竞争而相互等待，最终导致程序无法继续执行。在实际开发中，死锁可能会更复杂，可能涉及多个资源和多个协程，但基本原理是相同的。
