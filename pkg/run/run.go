package run

import (
	"os"
	"os/signal"
	"os/user"
	"path"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/frp/pkg/util/log"
)

var (
	u, _           = user.Current()
	FrpcBasePath   = path.Join(u.HomeDir, "/.wfrpc")
	SysMachine     = runtime.GOOS
	SysArch        = runtime.GOARCH
	logFile        = path.Join(os.TempDir(), "wfrpc.log")
	PidFile        = path.Join(os.TempDir(), "wfrpc.pid")
	FrpcConfigFile = func() string {
		p := path.Join(FrpcBasePath, "wfrpc.ini")
		if SysMachine == "windows" {
			return strings.Replace(p, "/", "\\", -1)
		}
		return p
	}()

	NoConfig                                            bool
	LocalPorts, RemotePorts                             []string
	Host, PortType, ServerAddr, ServerPort, ServerToken string
)

func RunFrpc() {
	os.Mkdir(FrpcBasePath, os.ModePerm)
	if !NoConfig {
		ShowConfig()
	}

	err := runClient(FrpcConfigFile)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

}

func DaemonRunFrpc() {
	cmd, err := RunDaemon(logFile)
	if err != nil {
		log.Error("后台运行wfrpc进程失败: %s", err)
		os.Exit(1)
	}

	if cmd != nil {
		return
	}

	RunFrpc()

}

func runClient(cfgFilePath string) error {
	cfg, pxyCfgs, visitorCfgs, err := config.ParseClientConfig(cfgFilePath)
	if err != nil {
		return err
	}
	return startService(cfg, pxyCfgs, visitorCfgs, cfgFilePath)
}

func startService(
	cfg config.ClientCommonConf,
	pxyCfgs map[string]config.ProxyConf,
	visitorCfgs map[string]config.VisitorConf,
	cfgFile string,
) (err error) {

	log.InitLog(cfg.LogWay, cfg.LogFile, cfg.LogLevel,
		cfg.LogMaxDays, cfg.DisableLogColor)

	if cfgFile != "" {
		log.Trace("start frpc service for config file [%s]", cfgFile)
		defer log.Trace("frpc service for config file [%s] stopped", cfgFile)
	}
	svr, errRet := client.NewService(cfg, pxyCfgs, visitorCfgs, cfgFile)
	if errRet != nil {
		err = errRet
		return
	}

	kcpDoneCh := make(chan struct{})
	// Capture the exit signal if we use kcp.
	if cfg.Protocol == "kcp" {
		go handleSignal(svr, kcpDoneCh)
	}

	err = svr.Run()
	if err == nil && cfg.Protocol == "kcp" {
		<-kcpDoneCh
	}
	return
}

func handleSignal(svr *client.Service, doneCh chan struct{}) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.GracefulClose(500 * time.Millisecond)
	close(doneCh)
}
