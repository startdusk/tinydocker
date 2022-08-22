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
		Run()
	case "init":
		Init()
	default:
		panic(fmt.Sprintf("[%s] not define", os.Args[1]))
	}
}

func Run() {
	cmd := exec.Command(os.Args[0], "init", os.Args[2])
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// syscall.CLONE_NEWUTS hostname 隔离
		// syscall.CLONE_NEWPID 进程隔离
		// syscall.CLONE_NEWNS 文件系统隔离
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

// func child() {
// 	cmd := exec.Command(os.Args[2])
// 	syscall.Sethostname([]byte("container"))
// 	// MS_NOEXEC: 在本文件系统中不允许运行其他程序
// 	// MS_NOSUID: 在本系统中运行程序的时候，不允许set-user-id和set-group-id
// 	// MS_NODEV:  从linux2.4以来，所有mount的系统都会默认设定的参数
// 	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
// 	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	if err := cmd.Run(); err != nil {
// 		panic(err)
// 	}
// 	syscall.Unmount("/proc", 0)
// }

func Init() {
	syscall.Sethostname([]byte("container"))
	syscall.Chroot("rootfs")
	syscall.Chdir("/")
	syscall.Mount("proc", "/proc", "proc", 0, "")
	syscall.Exec(os.Args[2], os.Args[2:], os.Environ())
	syscall.Unmount("/proc", 0)
}
