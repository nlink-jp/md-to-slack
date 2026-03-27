package main

import (
	"fmt"
	"io"
	"os"

	"github.com/nlink-jp/md-to-slack/internal/converter"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-V":
			fmt.Printf("md-to-slack %s\n", version)
			return
		case "--help", "-h":
			fmt.Fprintln(os.Stderr, "Usage: md-to-slack [--version]")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "Read Markdown from stdin, write Slack Block Kit JSON to stdout.")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "Flags:")
			fmt.Fprintln(os.Stderr, "  --version, -V   Print version and exit")
			fmt.Fprintln(os.Stderr, "  --help,    -h   Print this help and exit")
			return
		default:
			fmt.Fprintf(os.Stderr, "md-to-slack: unknown flag %q\n", os.Args[1])
			os.Exit(1)
		}
	}

	src, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "md-to-slack: read stdin: %v\n", err)
		os.Exit(1)
	}

	out, err := converter.ConvertToJSON(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "md-to-slack: convert: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stdout.Write(out); err != nil {
		fmt.Fprintf(os.Stderr, "md-to-slack: write stdout: %v\n", err)
		os.Exit(1)
	}
}
