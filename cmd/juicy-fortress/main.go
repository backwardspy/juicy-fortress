package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"

	"github.com/backwardspy/juicy-fortress/internal/stages"
)

var infoStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
var successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

func exit(format string, args ...any) {
	fmt.Fprintf(os.Stderr, errorStyle.Render(fmt.Sprintf(format, args...)))
	println()
	os.Exit(1)
}

func main() {
	installDir := "dwarffortress"

	fmt.Println(infoStyle.Render("downloading Dwarf Fortress"))
	if err := stages.DownloadDwarfFortress(installDir); err != nil {
		exit("failed to download Dwarf Fortress: %v", err)
	}

	fmt.Println(infoStyle.Render("installing DFHack mod"))
	if err := stages.InstallDFHack(installDir); err != nil {
		exit("failed to install DFHack: %v", err)
	}

	fmt.Println(infoStyle.Render("installing TWBT plugin"))
	if err := stages.InstallTWBT(installDir); err != nil {
		exit("failed to install TWBT: %v", err)
	}

	fmt.Println(infoStyle.Render("installing Spacefox tileset"))
	if err := stages.InstallSpacefox(installDir); err != nil {
		exit("failed to install Spacefox: %v", err)
	}

	fmt.Println(successStyle.Render("all done!"))
	fmt.Println("enter the \"dwarffortress\" folder and run \"dfhack\" to start Dwarf Fortress.")
	fmt.Print("please note that on some platforms the first startup may take a while.")
	fmt.Println(" be patient while DFHack does its thing!")
	fmt.Println("all startups past the first one will be much faster.")
}
