package cgroups

import (
	"github.com/sirupsen/logrus"
	"myDocker/cgroups/subsystems"
)

type CgroupManager struct {
	// cgroup 在 hierarchy 中的路径，相当于创建的 cgroup 目录相对于root cgroup 目录的路径
	Path string
	// 资源配置
	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// Apply 将进程 PID 加入到每个 cgroup
func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

// Set 设置各个 subsystem 挂载中的 cgroup 资源限制
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

// Destroy 释放各个 subsystem 挂载中的 cgroup
func (c *CgroupManager) Destroy() error {
	for _, subSysIn := range subsystems.SubsystemsIns {
		if err := subSysIn.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}

	return nil
}
