package run

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

func RunDaemon(logFile string) (*exec.Cmd, error) {
	envName := "WDAEMON" //环境变量名称
	envValue := "WFRPC"  //环境变量值

	val := os.Getenv(envName) //读取环境变量的值,若未设置则为空字符串
	if val == envValue {      //监测到特殊标识, 判断为子进程,不再执行后续代码
		err := ioutil.WriteFile(PidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	/*以下是父进程执行的代码*/

	//因为要设置更多的属性, 这里不使用`exec.Command`方法, 直接初始化`exec.Cmd`结构体
	cmd := &exec.Cmd{
		Path: os.Args[0],
		Args: os.Args,      //注意,此处是包含程序名的
		Env:  os.Environ(), //父进程中的所有环境变量
	}

	//为子进程设置特殊的环境变量标识
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", envName, envValue))

	//若有日志文件, 则把子进程的输出导入到日志文件
	if logFile != "" {
		stdout, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
		cmd.Stderr = stdout
		cmd.Stdout = stdout
	}

	//异步启动子进程
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
