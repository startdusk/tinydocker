实验一，实现独立的命名空间:

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

增强实验
使用两个 shell，一个进行进程隔离，一个改变 hostname

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	fmt.Printf("Process => %v [%d]\n", os.Args, os.Getegid())
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic(fmt.Sprintf("[%s] not define", os.Args[1]))
	}
}

func run() {
	cmd := exec.Command(os.Args[0], append([]string{"child"}, os.Args[2])...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func child() {
	cmd := exec.Command(os.Args[2])
	syscall.Sethostname([]byte("container"))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

```

实验过程

```bash
go build
sudo ./tinydocker run /bin/bash

Process => [./tinydocker run /bin/bash] [0] # 第一次进入run函数，run创建了一个子进程创建bash，并在bash中再次调用 sudo ./tinydocker child /bin/bash，于是就进入下面的child函数
Process => [./tinydocker child /bin/bash] [0] # 进入到child，执行修改主机名，由于child没有做隔离，所以它修改的是run创建的进程的hostname
root@container:/home/benjamin-linux/project/github/tinydocker# hostname
container
```
