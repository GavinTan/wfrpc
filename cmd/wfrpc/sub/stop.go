package sub

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gavintan/wfrpc/pkg/run"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Use: "help", Short: "帮助信息", Run: func(cmd *cobra.Command, args []string) { cmd.Parent().Help() }})
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止服务",
	Run: func(cmd *cobra.Command, args []string) {
		stop()
	},
}

func stop() {
	data, err := os.ReadFile(run.PidFile)
	if err != nil {
		log.Fatalln(err)
	}

	pid, _ := strconv.Atoi(strings.Trim(string(data), "\n"))
	p, _ := process.NewProcess(int32(pid))

	child, _ := p.Children()
	for _, c := range child {
		c.Kill()
	}

	p.Kill()
}
