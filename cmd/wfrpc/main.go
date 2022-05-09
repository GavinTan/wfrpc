package main

import (
	"fmt"
	"log"
	"wfrpc/cmd/wfrpc/sub"

	"github.com/spf13/cobra"
)

type sliceFlag []string

var (
	saddr, sport, ptype                string
	lport, rport                       sliceFlag
	clear, daemon, showdir, showconfig bool
)

func (f *sliceFlag) String() string {
	return fmt.Sprintf("%v", []string(*f))
}

func (f *sliceFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.LstdFlags)
	//允许windows下双击运行
	cobra.MousetrapHelpText = ""
}

func main() {

	sub.Execute()
}
