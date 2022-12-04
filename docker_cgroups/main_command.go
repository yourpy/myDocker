package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"myDocker/cgroups/subsystems"
	"myDocker/container"
	"strings"
)

// 定义 runCommand Flags ，其作用类似于运行命令时使用一来指定参数
var runCommand = cli.Command{
	Name:  "run",
	Usage: `Create a container with namespace and cgroups limit ie: myDocker run -ti [image] [command]`,
	/*
		这里是 run 命令执行的真正函数。
		1. 判断参数是否包含 command
		2. 获取用户指定的 command
		3. 调用 Run function 去准备启动容器
	*/
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		// 这个是获取启动容器时的命令
		// 如果本次中的 stress --vm-bytes 200m --vm-keep -m 1，空格分隔后会存储在下面的cmdArray中
		var cmdArr []string
		for _, arg := range context.Args() {
			cmdArr = append(cmdArr, arg)
		}
		tty := context.Bool("ti")
		resConfig := &subsystems.ResourceConfig{
			MemoryLimit: context.String("mem"),
			CpuShare:    context.String("cpuShare"),
			CpuSet:      context.String("cpuSet"),
		}
		Run(tty, cmdArr, resConfig)
		return nil
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		// 增加内存等限制参数
		cli.StringFlag{
			Name:  "mem",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpu",
			Usage: "cpu limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},
}

// 定义 initCommand 的具体操作，此操作为内部方法，禁止外部调用
var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	/*
		1.获取传递过来的command参数
		2.执行容器初始化操作
	*/
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		cmd := context.Args().Get(0)
		args := strings.Split(context.Args().Get(1), " ")
		log.Infof("command: %s, args: %s", cmd, args)
		return container.RunContainerInitProcess(cmd, args)
	},
}
