package ipfs

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	nodeURL string
)

var IPFSCmd = &cobra.Command{
	Use:   "ipfs",
	Short: "📤 Interact with IPFS",
}

func init() {
	IPFSCmd.PersistentFlags().StringVarP(&nodeURL, "node", "n", "https://ipfs-gateway.v8-bellecour.iex.ec", "IPFS node API URL")
	IPFSCmd.AddCommand(ipfsUploadCommand())
	IPFSCmd.AddCommand(ipfsDownloadCommand())
}

func ipfsUploadCommand() *cobra.Command {
	var filePath string

	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a local file to IPFS",
		RunE: func(cmd *cobra.Command, args []string) error {
			if filePath == "" {
				return fmt.Errorf("❌ you must provide a file path with --file")
			}

			fmt.Printf("📁 File: %s\n🔗 Node: %s\n", filePath, nodeURL)

			file, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("cannot open file: %w", err)
			}
			defer file.Close()

			pr, pw := io.Pipe()
			multipartWriter := multipart.NewWriter(pw)

			errChan := make(chan error, 1)

			go func() {
				defer pw.Close()
				defer multipartWriter.Close()

				part, errCreateFormFile := multipartWriter.CreateFormFile("file", filepath.Base(file.Name()))
				if errCreateFormFile != nil {
					errChan <- fmt.Errorf("failed to create form file: %w", errCreateFormFile)
					return
				}

				if _, errCopyFile := io.Copy(part, file); errCopyFile != nil {
					errChan <- fmt.Errorf("failed to copy file data: %w", errCopyFile)
					return
				}

				errChan <- nil
			}()

			req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v0/add", nodeURL), pr)
			if err != nil {
				return fmt.Errorf("failed to create request: %w", err)
			}
			req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("failed to send request: %w", err)
			}
			defer resp.Body.Close()

			if goroutineErr := <-errChan; goroutineErr != nil {
				return goroutineErr
			}

			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("IPFS upload failed: %s", string(body))
			}

			fmt.Println("✅ IPFS Response:", string(body))
			return nil
		},
	}

	uploadCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the local file to upload")

	return uploadCmd
}

func ipfsDownloadCommand() *cobra.Command {
	var cid string
	var output string

	downloadCmd := &cobra.Command{
		Use:   "download",
		Short: "Download a file from IPFS using its CID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cid == "" {
				return fmt.Errorf("❌ you must provide a CID with --cid")
			}
			if output == "" {
				output = cid
			}

			url := fmt.Sprintf("%s/ipfs/%s", nodeURL, cid)

			fmt.Printf("🔗 Downloading from: %s\n📁 Saving to: %s\n", url, output)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				return fmt.Errorf("failed to create HTTP request: %w", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("failed to download from IPFS: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("IPFS download failed: %s", string(body))
			}

			outFile, err := os.Create(output)
			if err != nil {
				return fmt.Errorf("cannot create output file: %w", err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, resp.Body); err != nil {
				return fmt.Errorf("error saving file: %w", err)
			}

			fmt.Println("✅ Download complete")
			return nil
		},
	}

	downloadCmd.Flags().StringVarP(&cid, "cid", "c", "", "CID of the file to download (required)")
	downloadCmd.Flags().StringVarP(&output, "out", "o", "", "Output file path (optional)")

	return downloadCmd
}
