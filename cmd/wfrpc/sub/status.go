package sub

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gavintan/wfrpc/pkg/run"

	psutil "github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Use: "help", Short: "帮助信息", Run: func(cmd *cobra.Command, args []string) { cmd.Parent().Help() }})
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "运行状态",
	Run: func(cmd *cobra.Command, args []string) {
		runnig := status()
		if !runnig {
			fmt.Println("wfrpc is stopped")
		}
	},
}

func status() bool {
	_, err := os.Stat(run.PidFile)

	if err == nil {
		data, err := os.ReadFile(run.PidFile)
		if err != nil {
			log.Fatalln(err)
		}

		pid, _ := strconv.Atoi(strings.Trim(string(data), "\n"))
		if ok, _ := psutil.PidExists(int32(pid)); ok {
			fmt.Printf("%s is running: %v\n", name, pid)
			return true
		} else {
			os.Remove(run.PidFile)
		}
	}

	return false
}
