package main

import (
	"flag"
	"github.com/thearchetypee/go-lc3-vm/vm"
	"log"
	"os"
)

func main() {
	imgFilePath := getProgramImageFilePath()
	flag.Parse()
	vm := vm.NewVM()
	vm.LoadProgramImage(imgFilePath)
	vm.Run()
}

func getProgramImageFilePath() string {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Printf("No program file specified")
		return ""
	}
	if info, err := os.Stat(args[0]); err != nil {
		log.Printf("Can't load %s: %s", args[0], err)
		return ""
	} else if info.IsDir() {
		log.Printf("%s is a directory, not a program file", args[0])
		return ""
	}
	return args[0]
}
