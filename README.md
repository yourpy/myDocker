# myDocker
手写docker

### docker_run
使用 namespace 生成容器
```
# 编译
go build .
# 启动一个容器
. /myDocker run -ti /bin/sh
```

### docker_cgroups
在 docker_run 基础上使用 cgroups 对容器进行资源限制
```
# 编译
go build .
# 启动一个容器
. /myDocker run -ti -mem 100m stress --vm-bytes 200m --vm-keep -m 1
# 通过 top 命令可以看到内存占用被限制在 100m
```

### docker_cgroups_addPipeline
在 docker_cgroups 基础上增加使用 pipeline 在父子进程之间传递消息

### docker_cgroups_addPivotroot
在 docker_cgroups_addPipeline 基础上增加实现新创建的容器和父进程目录不同的功能，目前docker 中还是使用的系统原有proc，不怎么纯净，所以使用 busybox 来更换 docker 的系统挂载点


### 进入容器
目前没有提供 exec 命令进入容器，可以自己手动进入
```sh
ps -ef  # 找出容器的 pid
nsenter -t 容器PID  -m -u -i -n -p  # 敲该命令进入容器

```
