package sub

import (
	"fmt"
	"os"
	"wfrpc/pkg/run"
	"wfrpc/pkg/tools"

	"github.com/spf13/cobra"
)

const version = "v1.0.0"
const name = "wfrpc"

var (
	clear, daemon, showDir, showConfig, showVersion, help bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "帮助信息")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "查看版本")
	rootCmd.PersistentFlags().BoolVarP(&showDir, "showdir", "", false, "查看配置目录")
	rootCmd.PersistentFlags().BoolVarP(&showConfig, "showconfig", "", false, "查看配置")
	rootCmd.PersistentFlags().BoolVarP(&clear, "clear", "", false, "清理配置")
	rootCmd.PersistentFlags().BoolVarP(&daemon, "daemon", "d", false, "后台运行")
	rootCmd.PersistentFlags().StringArrayVarP(&run.LocalPorts, "lport", "l", []string{}, "转发端口, 支持指定多个 example: --lport 22 --lport 80")
	rootCmd.PersistentFlags().StringArrayVarP(&run.RemotePorts, "rport", "r", []string{}, "转发后访问的端口, 不指定将随机生成, 支持指定多个与lport对应 example: --rport 22 --rport 80")
	rootCmd.PersistentFlags().StringVarP(&run.PortType, "type", "", "tcp", "指定端口类型 tcp or udp")
	rootCmd.PersistentFlags().StringVarP(&run.Host, "host", "", tools.GetHostIP(), "转发目标IP地址, 默认本机ip")
	rootCmd.PersistentFlags().StringVarP(&run.ServerAddr, "server_addr", "", "nb33.3322.org", "frp服务器地址, 默认nb33.3322.org")
	rootCmd.PersistentFlags().StringVarP(&run.ServerPort, "server_port", "", "7000", "frp服务器端口, 默认7000")
	rootCmd.PersistentFlags().StringVarP(&run.ServerToken, "server_token", "", "", "frp服务器认证token")
}

var rootCmd = &cobra.Command{
	Use:   "wfrpc",
	Short: "一个便捷的frp客户端 (https://github.com/gavintan/wfrpc)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if help {
			fmt.Println(cmd.Help())
			return nil
		}
		if showVersion {
			fmt.Println(version)
			return nil
		}
		if showDir {
			run.ShowDir()
			return nil
		}
		if showConfig {
			run.ShowConfig()
			return nil
		}
		if clear {
			os.Remove(run.FrpcConfigFile)
			fmt.Println("clear config success")
			return nil
		}

		start(daemon)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func start(daemon bool) {

	runnig := status()

	if !runnig {

		if daemon {
			run.DaemonRunFrpc()

		} else {
			run.RunFrpc()
		}

	}
}
