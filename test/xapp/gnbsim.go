package main

import (
	"github.com/spf13/viper"
)

// A simple XApp sample program that use "xapp" skeleton

func main() {
	if viper.GetString("test.mode") == "generator" {
		generator()
	} else {
		forwarder()
	}
}
