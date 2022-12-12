package subsystems

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct {
}

// Set 设置 cgroupPath 对应的 cgroup 内存资源限制
func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	// 获取自定义cgroup的路径，没有则创建，如：/sys/fs/cgroup/memory/mydocker-cgroup
	subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true)
	if err != nil {
		return err
	}

	log.Infof("%s cgroup path: %s", s.Name(), subsysCgroupPath)

	if res.MemoryLimit != "" {
		// 设置这个 cgroup 的内存限制，写入 cgroup 对应目录的 memory.limit_in_bytes 文件中
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
			return fmt.Errorf("set cgroup memory fail %v", err)
		}
	}
	return nil
}

// Remove 删除 cgroupPath 对应的 cgroup
func (s *MemorySubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}

// Apply 将一个进程加入到 cgroupPath 对应的 cgroup 中
func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}

func (s *MemorySubSystem) Name() string {
	return "memory"
}
