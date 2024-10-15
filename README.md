# trace
Go语言实现的贯穿式链路追踪
---
# Example
```go
package main

import (
	"fmt"
	"time"

	"github.com/dongrv/trace"
)

func main() {
	ctx := trace.New("hello-world")
	defer ctx.Clear()

	ctx.Set("en", []interface{}{"hello", "world"})
	ctx.SetKV("zh", []interface{}{"你", "好"})

	time.Sleep(time.Second) // 模拟运行耗时

	fmt.Println(ctx.Stop().String())
}
```
