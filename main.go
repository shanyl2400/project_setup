package main

import (
	"fmt"
	"os"
)

const (
	configFilePath = "./config.yaml"
)

func main() {
	loadConfig(configFilePath)
	InitRegex(config.Remove, config.Move)

	//flags
	params := os.Args[1:]

	if len(params) < 1 {
		Setup()
		return
	}
	switch params[0] {
	case "--setup":
		Setup()
	case "--release":
		Release()
	case "--clean":
		Clean()
	case "--restore":
		Restore(RestoreFlag(params[1]))
	default:
		fmt.Println("Unknown flag")
	}
}
