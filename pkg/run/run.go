package run

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strings"
)

var (
	u, _         = user.Current()
	FrpcBasePath = path.Join(u.HomeDir, "/.wfrpc")
	SysMachine   = runtime.GOOS
	SysArch      = runtime.GOARCH
	logFile      = path.Join(os.TempDir(), "wfrpc.log")
	PidFile      = path.Join(os.TempDir(), "wfrpc.pid")
	FrpcBin      = func() string {
		p := path.Join(FrpcBasePath, "frpc")
		if SysMachine == "windows" {
			return strings.Replace(p, "/", "\\", -1)
		}
		return p
	}()
	FrpcConfigFile = func() string {
		p := path.Join(FrpcBasePath, "frpc.ini")
		if SysMachine == "windows" {
			return strings.Replace(p, "/", "\\", -1)
		}
		return p
	}()

	LocalPorts, RemotePorts                             []string
	Host, PortType, ServerAddr, ServerPort, ServerToken string
)

func RunFrpc() {
	os.Mkdir(FrpcBasePath, os.ModePerm)

	DownloadFrpc()
	ShowConfig()

	cmd := exec.Command(FrpcBin, "-c", FrpcConfigFile)

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}

func DaemonRunFrpc() {
	cmd, err := RunDaemon(logFile)
	if err != nil {
		log.Fatal("运行frpc进程失败:", err)
	}

	if cmd != nil {
		return
	}

	RunFrpc()

}
