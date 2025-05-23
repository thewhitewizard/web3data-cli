package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const (
	AESKeySize256 = 256

	secretDirPerm   = 0700
	keyFilePerm     = 0600
	bitsToByteRatio = 8
)

var EncryptionCmd = &cobra.Command{
	Use:   "encryption",
	Short: "🔐 Manage data encryption and decryption",
}

func init() {
	EncryptionCmd.AddCommand(aesEncryptCommand())
	EncryptionCmd.AddCommand(aesDecryptCommand())
	EncryptionCmd.AddCommand(aesGenerateKeyCommand())
}

func aesEncryptCommand() *cobra.Command {
	var inputPath, outputPath, keyPath string

	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "🔒 Encrypt a file using AES (key is base64-encoded in a file)",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🔐 Starting AES encryption...")

			keyB64, err := os.ReadFile(keyPath)
			if err != nil {
				return fmt.Errorf("❌ Failed to read key file: %w", err)
			}
			key, err := base64.StdEncoding.DecodeString(string(keyB64))
			if err != nil {
				return fmt.Errorf("❌ Failed to decode base64 key: %w", err)
			}
			if len(key) != (AESKeySize256 / bitsToByteRatio) {
				return fmt.Errorf("❌ Invalid AES key length: %d bytes", len(key))
			}

			plaintext, err := os.ReadFile(inputPath)
			if err != nil {
				return fmt.Errorf("❌ Failed to read input file: %w", err)
			}

			plaintext = pkcs7Pad(plaintext, aes.BlockSize)

			block, err := aes.NewCipher(key)
			if err != nil {
				return fmt.Errorf("❌ Failed to create cipher: %w", err)
			}

			iv := make([]byte, aes.BlockSize)
			if _, errGenIV := rand.Read(iv); errGenIV != nil {
				return fmt.Errorf("❌ Failed to generate IV: %w", errGenIV)
			}

			ciphertext := make([]byte, len(plaintext))
			mode := cipher.NewCBCEncrypter(block, iv)
			mode.CryptBlocks(ciphertext, plaintext)

			outFile, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("❌ Failed to create output file: %w", err)
			}
			defer outFile.Close()

			if _, err := outFile.Write(iv); err != nil {
				return fmt.Errorf("❌ Failed to write IV: %w", err)
			}
			if _, err := outFile.Write(ciphertext); err != nil {
				return fmt.Errorf("❌ Failed to write ciphertext: %w", err)
			}

			fmt.Println("✅ File encrypted! 📁 Saved to:", outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputPath, "in", "i", "", "Plaintext input file path")
	cmd.Flags().StringVarP(&outputPath, "out", "o", "", "Encrypted output file path")
	cmd.Flags().StringVarP(&keyPath, "key", "k", "", "Base64-encoded AES key file")

	if err := cmd.MarkFlagRequired("in"); err != nil {
		log.Fatalf("❌ Failed to mark 'in' flag as required: %v", err)
	}

	if err := cmd.MarkFlagRequired("out"); err != nil {
		log.Fatalf("❌ Failed to mark 'out' flag as required: %v", err)
	}

	if err := cmd.MarkFlagRequired("key"); err != nil {
		log.Fatalf("❌ Failed to mark 'key' flag as required: %v", err)
	}

	return cmd
}

func aesDecryptCommand() *cobra.Command {
	var inputPath, outputPath, keyPath string

	cmd := &cobra.Command{
		Use:   "decrypt",
		Short: "🔓 Decrypt a file using AES (key is base64-encoded in a file)",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🔓 Starting AES decryption...")

			keyB64, err := os.ReadFile(keyPath)
			if err != nil {
				return fmt.Errorf("❌ Failed to read key file: %w", err)
			}
			key, err := base64.StdEncoding.DecodeString(string(keyB64))
			if err != nil {
				return fmt.Errorf("❌ Failed to decode base64 key: %w", err)
			}
			if len(key) != (AESKeySize256 / bitsToByteRatio) {
				return fmt.Errorf("❌ Invalid AES key length: %d bytes", len(key))
			}

			inFile, err := os.Open(inputPath)
			if err != nil {
				return fmt.Errorf("❌ Failed to open input file: %w", err)
			}
			defer inFile.Close()

			iv := make([]byte, aes.BlockSize)
			if _, errReadIV := io.ReadFull(inFile, iv); errReadIV != nil {
				return fmt.Errorf("❌ Failed to read IV: %w", errReadIV)
			}

			ciphertext, err := io.ReadAll(inFile)
			if err != nil {
				return fmt.Errorf("❌ Failed to read ciphertext: %w", err)
			}
			if len(ciphertext)%aes.BlockSize != 0 {
				return errors.New("❌ Ciphertext is not a multiple of the block size")
			}

			block, err := aes.NewCipher(key)
			if err != nil {
				return fmt.Errorf("❌ Failed to create cipher: %w", err)
			}

			mode := cipher.NewCBCDecrypter(block, iv)
			plaintext := make([]byte, len(ciphertext))
			mode.CryptBlocks(plaintext, ciphertext)

			plaintext, err = pkcs7Unpad(plaintext)
			if err != nil {
				return fmt.Errorf("❌ Failed to unpad plaintext: %w", err)
			}

			if err := os.WriteFile(outputPath, plaintext, keyFilePerm); err != nil {
				return fmt.Errorf("❌ Failed to write output file: %w", err)
			}

			fmt.Println("✅ File decrypted! 📁 Saved to:", outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputPath, "in", "i", "", "Encrypted input file path")
	cmd.Flags().StringVarP(&outputPath, "out", "o", "", "Decrypted output file path")
	cmd.Flags().StringVarP(&keyPath, "key", "k", "", "Base64-encoded AES key file")

	if err := cmd.MarkFlagRequired("in"); err != nil {
		log.Fatalf("❌ Failed to mark 'in' flag as required: %v", err)
	}

	if err := cmd.MarkFlagRequired("out"); err != nil {
		log.Fatalf("❌ Failed to mark 'out' flag as required: %v", err)
	}

	if err := cmd.MarkFlagRequired("key"); err != nil {
		log.Fatalf("❌ Failed to mark 'key' flag as required: %v", err)
	}

	return cmd
}

func aesGenerateKeyCommand() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "genkey",
		Short: "🔑 Generate a random AES key and store it in base64 in .secrets/",
		RunE: func(cmd *cobra.Command, args []string) error {
			keyLen := AESKeySize256 / bitsToByteRatio
			key := make([]byte, keyLen)
			if _, err := rand.Read(key); err != nil {
				return fmt.Errorf("❌ Failed to generate key: %w", err)
			}

			b64Key := base64.StdEncoding.EncodeToString(key)

			if output == "" {
				if err := os.MkdirAll(".secrets", secretDirPerm); err != nil {
					return fmt.Errorf("❌ Failed to create secrets directory: %w", err)
				}
				output = fmt.Sprintf(".secrets/aes-key-%d.b64", AESKeySize256)
			}

			if err := os.WriteFile(output, []byte(b64Key), keyFilePerm); err != nil {
				return fmt.Errorf("❌ Failed to write key file: %w", err)
			}

			fmt.Printf("✅ AES-%d key generated and saved to %s\n", AESKeySize256, output)
			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "out", "o", "", "Output file path (default: .secrets/aes-key-<size>.b64)")

	return cmd
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(data, padText...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid padding size")
	}
	paddingLen := int(data[len(data)-1])

	if paddingLen == 0 || paddingLen > len(data) {
		return nil, errors.New("invalid padding")
	}

	return data[:len(data)-paddingLen], nil
}
