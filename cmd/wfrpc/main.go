package main

import (
	_ "github.com/fatedier/frp/assets/frpc"
	"github.com/gavintan/wfrpc/cmd/wfrpc/sub"
	"github.com/spf13/cobra"
)

func init() {
	// log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.LstdFlags)
	//允许windows下双击运行
	cobra.MousetrapHelpText = ""
}

func main() {
	sub.Execute()
}
