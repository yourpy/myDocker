package main

import (
	log "github.com/sirupsen/logrus"
	"myDocker/container"
	"os"
)

/*
Run
Start 调用command，它首先会 clone 出来一个 name space 隔离的进程，然后在子进程中，调用/proc/self/exe，也就是调用自己，发送 init 参数，调用 init 方法，
去初始化容器的一些资源。
*/
func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}
