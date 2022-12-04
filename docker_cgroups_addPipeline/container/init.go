package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

/*
RunContainerInitProcess
这里的 init 函数是在容器内部执行的，也就是说 代码执行到这里后，容器所在进程其实就已经创建出来了，
这是本容器执行的第 1 个进程。
使用 mount 先去挂载 proc 文件系统，以便后面通过 ps 等系统命令去查看当前进程资源情况。
*/
func RunContainerInitProcess() (err error) {
	// private 方式挂载，不影响宿主机的挂载
	if err = syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		logrus.Errorf("private 方式挂载 failed: %v", err)
		return err
	}

	// MS_NOEXEC 本文件系统不允许执行其他程序
	// MS_NOSUID 不允许 set-user-ID 和 set-group-ID
	// MS_NODEV  默认参数
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		logrus.Errorf("proc 挂载 failed: %v", err)
		return err
	}

	cmdArr := readUserCommand()
	if cmdArr == nil || len(cmdArr) == 0 {
		return fmt.Errorf("run container get user command error, cmdArr is nil")
	}

	// 使用lookPath的方式去查找命令进行执行
	path, err := exec.LookPath(cmdArr[0])
	if err != nil {
		logrus.Errorf("can't find exec path: %s %v", cmdArr[0], err)
		return err
	}
	logrus.Infof("find path: %s", path)

	// 完成初始化动作并将用户程序运行起来
	// 将当前进程的PID置为1
	// 调用了Kernel的 int execve(const char *filename, char *const argv[], char *const envp[])
	// 作用是执行 filename 对应程序。覆盖当前进程的镜像、数据和堆栈等信息，包括 PID 这些都会被将要运行的进程覆盖掉。
	// 将用户指定的进程（filename）运行起来，把最初的 init 进程给替换掉
	if err = syscall.Exec(path, cmdArr, os.Environ()); err != nil {
		logrus.Errorf("syscall exec err: %v", err.Error())
		return
	}

	return nil
}

func readUserCommand() []string {
	// 进程默认三个管道，从 fork 那边传过来的就是第 4 个
	// uintpr(3) 指 index 为 3 的文件描述符，就是传递进来的管道的一端
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("read init argv pipe err: %v", err)
		return nil
	}
	return strings.Split(string(msg), " ")
}
