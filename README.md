# 🚀 gcp - Go 协程池 (Goroutine Pool)

`gcp` 是一个轻量级、高性能的 Go 语言协程池实现。它可以有效地控制并发数量，复用 Goroutine 资源，并提供任务提交和池管理功能。

## ✨ 特性

- **🏗️ 固定容量**: 限制系统最大 Goroutine 数量，防止资源耗尽。
- **♻️ 资源复用**: 任务执行完成后，Goroutine 进入空闲状态等待新任务，减少创建和销毁开销。
- **📤 任务提交**: 提供简单直观的 `Submit` 接口提交异步任务。
- **🛑 平滑关闭**: 支持 `ShutDown` 操作，确保所有已提交的任务执行完毕后再关闭。
- **🛡️ Panic 恢复**: 内置 Panic 处理机制，防止单个任务崩溃导致整个进程退出。
- **📋 错误处理**: 明确的错误定义，如池满或池已关闭。

## 📦 安装

```bash
go get github.com/yourusername/gcp
```

## ⚡ 快速开始

### 基础用法

```go
package main

import (
	"fmt"
	"github.com/yourusername/gcp"
	"time"
)

func main() {
	// 1. 创建一个容量为 5 的协程池
	pool := gcp.New(5)
	
	// 确保在程序退出前关闭协程池
	defer pool.ShutDown()

	// 2. 提交任务
	for i := 1; i <= 10; i++ {
		innerI := i
		err := pool.Submit(func() {
			fmt.Printf("正在处理任务 %d\n", innerI)
			time.Sleep(time.Second)
			fmt.Printf("任务 %d 处理完成\n", innerI)
		})

		if err != nil {
			fmt.Printf("提交任务 %d 失败: %v\n", innerI, err)
		}
	}

	// 等待一会查看输出
	time.Sleep(time.Second * 3)
}
```

## 📖 API 说明

### `New(capacity int32) *Pool`
创建一个新的协程池。
- `capacity`: 协程池的最大并行度。

### `(p *Pool) Submit(task Task) error`
向池中提交一个任务（`func()`）。
- 如果池中有空闲的 Worker，则立即执行。
- 如果池已满且无法扩展，则返回 `PoolFullError`。
- 如果池已关闭，则返回 `PoolClosedError`。

### `(p *Pool) ShutDown()`
关闭协程池。该方法会阻塞，直到所有正在运行的任务处理完毕。

## ⚠️ 错误类型

| 错误变量 | 说明 |
| :--- | :--- |
| `InvalidCapacityError` | 容量设置无效（必须大于 0） |
| `PoolFullError` | 协程池已满，无法接收新任务 |
| `PoolClosedError` | 协程池已关闭，禁止提交任务 |

## 📄 许可证

MIT License
