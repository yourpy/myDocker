package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

/*
RunContainerInitProcess
这里的 init 函数是在容器内部执行的，也就是说 代码执行到这里后，容器所在进程其实就已经创建出来了，
这是本容器执行的第 1 个进程。
使用 mount 先去挂载 proc 文件系统，以便后面通过 ps 等系统命令去查看当前进程资源情况。
*/
func RunContainerInitProcess(command string, args []string) (err error) {
	logrus.Infof("RunContainerInitProcess command %s, args %s", command, args)

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

	// 使用lookPath的方式去查找命令进行执行
	path, err := exec.LookPath(command)
	if err != nil {
		logrus.Errorf("can't find exec path: %s %v", command, err)
		return err
	}
	logrus.Infof("find path: %s", path)

	// 完成初始化动作并将用户程序运行起来
	// 将当前进程的PID置为1
	// 调用了Kernel的 int execve(const char *filename, char *const argv[], char *const envp[])
	// 作用是执行 filename 对应程序。覆盖当前进程的镜像、数据和堆栈等信息，包括 PID 这些都会被将要运行的进程覆盖掉。
	// 将用户指定的进程（filename）运行起来，把最初的 in it 进程给替换掉
	if err = syscall.Exec(path, args, os.Environ()); err != nil {
		logrus.Errorf("syscall exec err: %v", err.Error())
		return
	}

	return nil
}
