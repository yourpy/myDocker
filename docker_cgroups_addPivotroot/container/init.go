package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
	// 挂载
	if err := setUpMount(); err != nil {
		logrus.Errorf("setUpMount fail: %v", err)
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

func setUpMount() error {
	// 首先设置根目录为私有模式，防止影响pivot_root
	// private 方式挂载，不影响宿主机的挂载
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		logrus.Errorf("private 方式挂载 failed: %v", err)
		return err
	}

	// 获取当前路径
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("get current location err: %v", err)
		return err
	}
	logrus.Infof("current location: %s", pwd)

	err = privotRoot(pwd)
	if err != nil {
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

	syscall.Mount("tmpfs", "/dev", "tempfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
	return nil
}

func privotRoot(root string) error {
	// 为了使当前root的老root和新root不在同一个文件系统下，需要把root重新mount一次
	// bind mount 是把相同的内容换了一个挂载点的挂载方法
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}

	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	// 判断当前目录是否已有该文件夹
	if _, err := os.Stat(pivotDir); err == nil {
		// 存在则删除
		if err := os.Remove(pivotDir); err != nil {
			return err
		}
	}
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return fmt.Errorf("mkdir of pivot_root err: %v", err)
	}

	// pivot_root 到新的 rootfs，老的 old_root 现在挂载到 rootfs/.pivot_root 上
	// 挂载点目前依然可以在mount命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root err: %v", err)
	}

	// 修改当前工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir root err: %v", err)
	}

	// 取消临时文件.pivot_root的挂载并删除它
	// 注意当前已经在根目录下，所以临时文件的目录也改变了
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir err: %v", err)
	}

	return os.Remove(pivotDir)
}
