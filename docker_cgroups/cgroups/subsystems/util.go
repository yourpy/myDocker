package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

/*
91 86 0:47 / /sys/fs/cgroup/cpu rw,nosuid,nodev,noexec,relatime - cgroup cgroup rw,cpu
92 86 0:48 / /sys/fs/cgroup/cpuacct rw,nosuid,nodev,noexec,relatime - cgroup cgroup rw,cpuacct
95 86 0:49 / /sys/fs/cgroup/blkio rw,nosuid,nodev,noexec,relatime - cgroup cgroup rw,blkio
97 86 0:21 / /sys/fs/cgroup/memory rw,nosuid,nodev,noexec,relatime - cgroup cgroup rw,memory
*/

// FindCgroupMountpoint 通过 /proc/self/mountinfo 找出挂载了某个 subsystem 的 hierarchy cgroup 根节点所在的目录
func FindCgroupMountpoint(subsystem string) string {
	file, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

// GetCgroupPath 找到 cgroup 在文件中的绝对路径
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	// 找到cgroup的根目录，如：/sys/fs/cgroup/memory
	cgroupRoot := FindCgroupMountpoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}
