package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/thewhitewizard/web3data-cli/cmd"
)

func main() {
	if err := fang.Execute(context.TODO(), cmd.RootCmd); err != nil {
		os.Exit(1)
	}
}
