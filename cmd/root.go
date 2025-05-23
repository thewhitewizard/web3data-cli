package cmd

import (
	"github.com/thewhitewizard/web3data-cli/cmd/ipfs"
	"github.com/thewhitewizard/web3data-cli/cmd/version"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "web3datacli",
	Long: `A CLI tool to manage web3data`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(ipfs.IPFSCmd)
	rootCmd.AddCommand(version.VersionCmd)
}
