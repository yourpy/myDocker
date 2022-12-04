package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

/*
NewParentProcess
这里是父进程，也就是当前进程执行的内容
1. 这里的/proc/self/exe 调用中，/proc/self/ 指的是当前运行进程自己的环境， exec 其实就是自己
调用了自己，使用这种方式对创建出来的进程进行初始化
2. 后面的 args 是参数，其中 init 是传递给本进程的第 1 个参数，其实就是会去调用 initCommand
去初始化进程的一些环境和资源
3. 下面的 clone 参数就是去 fork 出来一个新进程，并且使用了 namespace 隔离新创建的进程和外部环境
4. 如果用户指定了－ti 参数，就需要把当前进程的输入输出导入到标准输入输出上
*/
func NewParentProcess(tty bool, cmdArr []string) *exec.Cmd {
	commands := strings.Join(cmdArr, " ")
	log.Infof("command all is %s ", commands)
	// 传入，cmdArray[0]是启动的程序，比如本篇中的是stress，command是完整的命令：stress --vm-bytes 200m --vm-keep -m 1
	cmd := exec.Command("/proc/self/exe", "init", cmdArr[0], commands)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}
