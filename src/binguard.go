package main

import (
	"binlock/lockops"
	"binlock/progress"
	"binlock/title"
	"flag"
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	title.PrintBinGuardTitle()
	inputFile := flag.String("i", "", "Input binary path")
	outputFile := flag.String("o", "", "Output protected binary path")
	password := flag.String("pass", "", "Password to protect the binary")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" || *password == "" {
		fmt.Println("Please provide all required parameters")
		flag.PrintDefaults()
		return
	}

	ciphertext, err := lockops.BinLockWithSalt(*inputFile, *password)
	if err != nil {
		fmt.Println("Error in encrypting file:", err)
		return
	}

	err = lockops.CreateBinLockerFile(*outputFile, ciphertext)
	if err != nil {
		fmt.Println("Error in creating locker file:", err)
		return
	}
	m := progress.NewModel()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Progressbar error:", err)
		return
	}

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true).PaddingLeft(2).PaddingTop(0).Render("Binguard applied!"))
	outpath, err := filepath.Abs(*outputFile)
	if err != nil {
		fmt.Println("Error in creating getting absolute filepath:", err)
		return
	}
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true).PaddingLeft(2).PaddingTop(0).Render("Protected file path:", outpath))
}
