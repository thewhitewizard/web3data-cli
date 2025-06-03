# 🌐 Web3 Data CLI

A simple, secure, and extensible command-line tool written in Go for managing Web3-related data.

**Web3 Data CLI** (`web3datacli`) helps you handle sensitive data through strong AES encryption and enables seamless file sharing via IPFS.

---

## ✨ Features

- 🔐 AES encryption & decryption (256-bit)
- 🗂 Local key generation and management
- 📤 File uploads and download to IPFS
- 🧩 Built with extensibility in mind — easy to plug in more Web3 features
- 💻 Native binaries for Linux and macOS

---

## 📦 Installation

Download the latest binary from the [Releases](https://github.com/your-org/web3datacli/releases) page.

### Linux

```bash
wget https://github.com/your-org/web3datacli/releases/latest/download/web3datacli-linux -O web3datacli
chmod +x web3datacli
sudo mv web3datacli /usr/local/bin/web3datacli
```

### macOS

```bash
wget https://github.com/your-org/web3datacli/releases/latest/download/web3datacli-darwin -O web3datacli
chmod +x web3datacli
sudo mv web3datacli /usr/local/bin/web3datacli
```

---

## 🚀 Usage

```bash
web3datacli [command]
```

### Available Commands

| Command      | Description                             |
| ------------ | --------------------------------------- |
| `arweave`    | 🕸️ Interact with Arweave                 |
| `encryption` | 🔐 Manage data encryption and decryption |
| `ipfs`       | 📤 Interact with IPFS                    |
| `version`    | Show the CLI version                    |
| `completion` | Generate autocompletion for your shell  |
| `help`       | Help about any command                  |

---

## 🧠 Autocompletion

Enable autocompletion for your shell (e.g., bash, zsh):

```bash
web3datacli completion zsh >  ~/.zsh_completions/_web3datacli 
source ~/.zsh_completions/_web3datacli 
```

---

## 🧾 License

MIT License © 2025 F.CORDIER

---

## 🤝 Contributions

Feel free to open issues or submit PRs — contributions are welcome!
