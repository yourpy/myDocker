package main

import (
	log "github.com/sirupsen/logrus"
	"myDocker/cgroups"
	"myDocker/cgroups/subsystems"
	"myDocker/container"
	"os"
)

/*
Run
Start 调用command，它首先会 clone 出来一个 name space 隔离的进程，然后在子进程中，调用/proc/self/exe，也就是调用自己，发送 init 参数，调用 init 方法，
去初始化容器的一些资源。
*/
func Run(tty bool, comArr []string, res *subsystems.ResourceConfig) {
	// 进程已经启动完成
	parent := container.NewParentProcess(tty, comArr)
	if err := parent.Start(); err != nil {
		log.Error(err)
		return
	}

	// 创建 cgroup manager，并通过调用 set 和 apply 设置资源限制并使限制在容器上生效
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")

	defer cgroupManager.Destroy()

	// 将当前的进程的PID加入资源限制列表，如果内存的就是将PID写入tasks文件中
	if err := cgroupManager.Apply(parent.Process.Pid); err != nil {
		log.Errorf("cgroup apply err: %v", err)
		return
	}
	// 设置资源限制
	if err := cgroupManager.Set(res); err != nil {
		log.Errorf("cgoup set err: %v", err)
		return
	}

	log.Infof("parent process run")

	parent.Wait()
	os.Exit(-1)
}
