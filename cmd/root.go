package cmd

import (
	"github.com/thewhitewizard/web3data-cli/cmd/arweave"
	"github.com/thewhitewizard/web3data-cli/cmd/encryption"
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
	rootCmd.AddCommand(encryption.EncryptionCmd)
	rootCmd.AddCommand(ipfs.IPFSCmd)
	rootCmd.AddCommand(arweave.ArweaveCmd)
	rootCmd.AddCommand(version.VersionCmd)
}
