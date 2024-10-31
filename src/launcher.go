package main

import (
	"binlock/lockops"
	"binlock/title"
	"fmt"
	"log"
	"os"
)

func main() {
	title.PrintLauncherTitle()
	/*inputFile := flag.String("ipf", "", "Input protected binary path")
	password := flag.String("pass", "", "Password to decrypt the binary")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	binaryArgs := flag.Args()

	if *inputFile == "" || *password == "" {
		fmt.Println("Please provide all required parameters")
		flag.PrintDefaults()
		return
	}

	plaintext, err := lockops.UnlockWithSalt(*inputFile, *password, *debug)
	if err != nil {
		fmt.Println("Error in decrypting file:", err)
		return
	}
	lockops.RunProtectedBinary(plaintext, *inputFile, binaryArgs, *debug)*/

	if len(os.Args) < 3 {
		log.Fatal("Usage: ./launcher <protected-binary> <password> [binary args...]")
	}

	protectedBinaryPath := os.Args[1]
	password := os.Args[2]
	var binaryArgs []string

	if len(os.Args) > 3 {
		binaryArgs = os.Args[3:]
	}

	// Check if the protected binary file exists
	if _, err := os.Stat(protectedBinaryPath); os.IsNotExist(err) {
		log.Fatalf("Protected binary does not exist: %s", protectedBinaryPath)
	}

	plaintext, err := lockops.UnlockWithSalt(protectedBinaryPath, password, false)
	if err != nil {
		fmt.Println("Error in decrypting file:", err)
		return
	}
	lockops.RunProtectedBinary(plaintext, protectedBinaryPath, binaryArgs, false)

}
