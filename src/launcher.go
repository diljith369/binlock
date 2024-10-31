package main

import (
	"binlock/lockops"
	"binlock/title"
	"flag"
	"fmt"
)

func main() {
	title.PrintLauncherTitle()
	inputFile := flag.String("ipf", "", "Input protected binary path")
	password := flag.String("pass", "", "Password to decrypt the binary")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *inputFile == "" || *password == "" {
		fmt.Println("Please provide all required parameters")
		flag.PrintDefaults()
		return
	}

	err := lockops.UnlockWithSalt(*inputFile, *password, *debug)
	if err != nil {
		fmt.Println("Error in decrypting file:", err)
		return
	}
}
