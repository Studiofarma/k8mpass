package config

import "flag"

var (
	Config string
	Plugin string
)

func LoadFlags() {
	flag.StringVar(&Config, "kubeconfig", "", "specify kubernetes config file to use")
	flag.StringVar(&Plugin, "plugin", "", "path to plugin file")
	flag.Parse()
}
