package cmd

import (
	"github.com/thewhitewizard/web3data-cli/cmd/arweave"
	"github.com/thewhitewizard/web3data-cli/cmd/encryption"
	"github.com/thewhitewizard/web3data-cli/cmd/ipfs"
	"github.com/thewhitewizard/web3data-cli/cmd/version"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:  "web3datacli",
	Long: `A CLI tool to manage web3data`,
}

func init() {
	RootCmd.AddCommand(encryption.EncryptionCmd)
	RootCmd.AddCommand(ipfs.IPFSCmd)
	RootCmd.AddCommand(arweave.ArweaveCmd)
	RootCmd.AddCommand(version.VersionCmd)
}
