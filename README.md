# myDocker
手写docker

### docker_run
使用 namespace 生成容器，编译后可使用 `. /myDocker run -ti /bin/ls` 可生成一个容器

### docker_cgroups
在 docker_run 基础上使用 cgroups 对容器进行资源限制

### docker_cgroups_addPipeline
在 docker_cgroups 基础上增加使用 pipeline 在父子进程之间传递消息

### 进入容器
目前没有提供 exec 命令进入容器，可以自己手动进入
```sh
ps -ef  # 找出容器的 pid
nsenter -t 容器PID  -m -u -i -n -p  # 敲该命令进入容器

```
