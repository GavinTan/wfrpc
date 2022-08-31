package run

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/gavintan/wfrpc/pkg/tools"
	"gopkg.in/ini.v1"
)

type configType map[string]string

var configList = []configType{}

func FrpcConfig() {
	_, err := os.Stat(FrpcConfigFile)
	var cfg *ini.File

	if err != nil {
		cfg = ini.Empty()
	} else {
		cfg, _ = ini.Load(FrpcConfigFile)
	}

	common, _ := cfg.GetSection("common")
	if common == nil {
		cfg.NewSection("common")
	}

	cfg.Section("common").Key("server_addr").SetValue(ServerAddr)
	cfg.Section("common").Key("server_port").SetValue(ServerPort)
	cfg.Section("common").Key("token").SetValue(ServerToken)

	n := tools.Abs(len(LocalPorts) - len(RemotePorts))
	if len(RemotePorts) != 0 && n != 0 {
		for i := 0; i < n; i++ {
			RemotePorts = append(RemotePorts, tools.RandomPort())
		}
	}

	for _, lp := range LocalPorts {
		sectionName := fmt.Sprintf("%v_%v", Host, lp)
		s, _ := cfg.GetSection(sectionName)
		if s == nil {
			s, _ = cfg.NewSection(sectionName)
		}
		s.Key("local_ip").SetValue(Host)
		s.Key("local_port").SetValue(lp)
		s.Key("type").SetValue(PortType)

		if len(RemotePorts) > 0 {
			s.Key("remote_port").SetValue(RemotePorts[tools.IndexOf(LocalPorts, lp)])
		} else {
			s.Key("remote_port").SetValue(tools.RandomPort())
		}
	}

	cfg.SaveTo(FrpcConfigFile)

	for _, s := range cfg.Sections() {
		config := configType{}
		if s.Name() != "DEFAULT" && s.Name() != "common" {
			for _, v := range s.Keys() {
				config[v.Name()] = v.Value()
			}
			configList = append(configList, config)
		}
	}
}

func ShowConfig() {
	FrpcConfig()
	f1 := strings.Repeat("-", 60)
	f2 := strings.Repeat(" ", 20)
	f3 := strings.Repeat("-", (len(f1)-len("config"))/2)

	fmt.Printf("%v+%vconfig%v+\n", f2, f3, f3)
	for _, c := range configList {
		msg := fmt.Sprintf("%v:%v --> %v:%v", c["local_ip"], c["local_port"], ServerAddr, c["remote_port"])
		n := math.Round(float64((60 - len(msg)) / 2))
		msgFill := strings.Repeat(" ", int(n))
		fmt.Printf("%v%v%v%v\n", f2, msgFill, msg, msgFill)
	}
	fmt.Printf("%v+%v+\n", f2, f1)
}

func ShowDir() {
	fmt.Printf("config: %v\n", FrpcConfigFile)
}
