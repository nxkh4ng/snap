package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

const version = "0.1.0"

var (
	withAll            *bool
	withDesc           *bool
	withFooter         *bool
	withBreakingChange *bool
	withTicket         *bool
	showVersion        *bool
	showHelp           *bool
)

func initFlags() {
	withAll = flag.BoolP("all", "a", false, "add all extras")
	withDesc = flag.BoolP("desc", "d", false, "add description")
	withFooter = flag.BoolP("footer", "f", false, "add footer")
	withBreakingChange = flag.BoolP("breaking", "b", false, "add breaking change")
	withTicket = flag.BoolP("ticket", "t", false, "add ticket")
	showVersion = flag.BoolP("version", "v", false, "show version")
	showHelp = flag.BoolP("help", "h", false, "show help")

	flag.CommandLine.SortFlags = false
	flag.Usage = printUsage

	flag.Parse()
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "snap — conventional commit CLI\n\n")
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  snap [flags]\n\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	fmt.Fprintf(os.Stderr, "  -a, --all        add all extras\n")
	fmt.Fprintf(os.Stderr, "  -d, --desc       add description\n")
	fmt.Fprintf(os.Stderr, "  -f, --footer     add footer\n")
	fmt.Fprintf(os.Stderr, "  -b, --breaking   add breaking change\n")
	fmt.Fprintf(os.Stderr, "  -t, --ticket     add ticket\n")
	fmt.Fprintf(os.Stderr, "  -v, --version    show version\n")
	fmt.Fprintf(os.Stderr, "  -h, --help       show help\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  snap\n")
	fmt.Fprintf(os.Stderr, "  snap -t\n")
	fmt.Fprintf(os.Stderr, "  snap -dt\n")
	fmt.Fprintf(os.Stderr, "  snap -a\n")
}

func handleFlags() bool {
	if *withAll {
		*withDesc = true
		*withFooter = true
		*withBreakingChange = true
		*withTicket = true
	}

	if *showVersion {
		fmt.Printf("snap version %s\n", version)
		return true
	}

	if *showHelp {
		printUsage()
		return true
	}

	return false
}
