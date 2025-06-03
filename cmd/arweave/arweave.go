package arweave

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	apiBaseURL string
)

var ArweaveCmd = &cobra.Command{
	Use:   "arweave",
	Short: "🕸️ Interact with Arweave",
}

func init() {

	ArweaveCmd.AddCommand(arweaveUploadCommand())
	ArweaveCmd.AddCommand(arweaveDownloadCommand())
}

func arweaveUploadCommand() *cobra.Command {
	var filePath string

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a file to Arweave through a relay API",
		RunE: func(cmd *cobra.Command, args []string) error {
			if filePath == "" {
				return fmt.Errorf("❌ you must provide a file path with --file")
			}

			file, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("cannot open file: %w", err)
			}
			defer file.Close()

			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)
			part, err := writer.CreateFormFile("file", filepath.Base(filePath))
			if err != nil {
				return fmt.Errorf("cannot create form file: %w", err)
			}

			if _, err := io.Copy(part, file); err != nil {
				return fmt.Errorf("failed to copy file data: %w", err)
			}

			if err := writer.Close(); err != nil {
				return fmt.Errorf("failed to close multipart writer: %w", err)
			}

			resp, err := http.Post(apiBaseURL+"/upload", writer.FormDataContentType(), &buf)
			if err != nil {
				return fmt.Errorf("upload request failed: %w", err)
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("arweave upload failed: %s", string(body))
			}

			var res struct {
				ArweaveID string `json:"arweaveId"`
				URL       string `json:"url"`
			}
			if err := json.Unmarshal(body, &res); err != nil {
				return fmt.Errorf("invalid response: %w", err)
			}

			fmt.Println("✅ Upload successful")
			fmt.Printf("🆔 Arweave ID: %s\n🔗 URL: %s\n", res.ArweaveID, res.URL)
			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the local file to upload")
	cmd.Flags().StringVarP(&apiBaseURL, "api", "a", "http://localhost:3000", "Base URL of the Arweave relay API (without /upload)")

	return cmd
}

func arweaveDownloadCommand() *cobra.Command {
	var id string
	var output string

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download a file from Arweave using its transaction ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == "" {
				return fmt.Errorf("❌ you must provide an Arweave transaction ID with --id")
			}
			if output == "" {
				output = id
			}

			url := fmt.Sprintf("https://arweave.net/%s", id)
			fmt.Printf("🔗 Downloading from: %s\n📁 Saving to: %s\n", url, output)

			resp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("failed to download from Arweave: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("arweave download failed: %s", string(body))
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

	cmd.Flags().StringVarP(&id, "id", "i", "", "Arweave transaction ID to download (required)")
	cmd.Flags().StringVarP(&output, "out", "o", "", "Output file path (optional)")

	return cmd
}
