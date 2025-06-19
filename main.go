package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/thewhitewizard/web3data-cli/cmd"
)

func main() {
	cmd.Execute()
	if err := fang.Execute(context.TODO(), cmd.RootCmd); err != nil {
		os.Exit(1)
	}
}
