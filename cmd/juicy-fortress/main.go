package main

import (
	"flag"
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

type args struct {
    installDir string
    verbose bool
}

func parseArgs() args {
    args := args{}
    flag.StringVar(&args.installDir, "installDir", "dwarffortress", "dir to install into")
    flag.BoolVar(&args.verbose, "verbose", false, "print more info")
    flag.Parse()
    return args
}

func main() {
    args := parseArgs()

	fmt.Println(infoStyle.Render("downloading Dwarf Fortress"))
	if err := stages.DownloadDwarfFortress(args.installDir, args.verbose); err != nil {
		exit("failed to download Dwarf Fortress: %v", err)
	}

	fmt.Println(infoStyle.Render("installing DFHack mod"))
	if err := stages.InstallDFHack(args.installDir, args.verbose); err != nil {
		exit("failed to install DFHack: %v", err)
	}

	fmt.Println(infoStyle.Render("installing TWBT plugin"))
	if err := stages.InstallTWBT(args.installDir, args.verbose); err != nil {
		exit("failed to install TWBT: %v", err)
	}

	fmt.Println(infoStyle.Render("installing Spacefox tileset"))
	if err := stages.InstallSpacefox(args.installDir, args.verbose); err != nil {
		exit("failed to install Spacefox: %v", err)
	}

    fmt.Println(infoStyle.Render("applying patches"))
    if err := stages.ApplyPatches(args.installDir, args.verbose); err != nil {
        exit("failed to apply patches: %v", err)
    }

	fmt.Println(successStyle.Render("all done!"))
	fmt.Println("enter the \"dwarffortress\" folder and run \"dfhack\" to start Dwarf Fortress.")
	fmt.Print("please note that on some platforms the first startup may take a while.")
	fmt.Println(" be patient while DFHack does its thing!")
	fmt.Println("all startups past the first one will be much faster.")
}
