实验一，测试 UTS 隔离:

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic(fmt.Sprintf("[%s] not define", os.Args[1]))
	}
}

func run() {
	cmd := exec.Command(os.Args[2])
	cmd.SysProcAttr = &syscall.SysProcAttr{
        // 这里是使用UTS隔离
        // 注意: 2022/08/24 Mac M1 尚不支持此特性
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

```

编译运行上面代码

```bash
$ go build -o tinydocker .
```

执行命令 shell

```bash
$ sudo ./tinydocker run sh
```

修改主机名

```bash
# hostname
会展示你的主机名
# hostname -b ben 修改主机名为ben
# hostname
ben 会在当前的fork的子进程中修改hostname
```

在电脑上打开另一个 shell，查看主机名，发现没有被改变，说明隔离起作用了

```bash
$ hostname
```
