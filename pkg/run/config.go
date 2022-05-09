package run

import (
	"fmt"
	"math"
	"os"
	"strings"
	"wfrpc/pkg/tools"

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
		section, _ := cfg.NewSection("common")
		section.NewKey("server_addr", "nb33.3322.org")
		section.NewKey("server_port", "7000")
		section.NewKey("token", "")
	} else {
		cfg.Section("common").Key("server_addr").SetValue(ServerAddr)
		cfg.Section("common").Key("server_port").SetValue(ServerPort)
		cfg.Section("common").Key("token").SetValue(ServerToken)
	}

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
	fmt.Printf("\t\t+%vconfig%v+\n", strings.Repeat("-", 27), strings.Repeat("-", 27))
	for _, c := range configList {
		msg := fmt.Sprintf("%v:%v --> %v:%v", c["local_ip"], c["local_port"], ServerAddr, c["remote_port"])
		n := math.Round(float64((60 - len(msg)) / 2))
		o := strings.Repeat(" ", int(n))
		fmt.Printf("\t\t%v%v%v\n", o, msg, o)
	}
	fmt.Printf("\t\t+%v+\n", strings.Repeat("-", 60))
}

func ShowDir() {
	fmt.Printf("bin: %v\nconfig: %v\n", FrpcBin, FrpcConfigFile)
}
