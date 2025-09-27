# 🚀 kubelocal

> **One-line installer for local Kubernetes development environment**

[![Release](https://img.shields.io/badge/release-v0.0.5-blue.svg)](https://github.com/Antoine-Flo/kubelocal/releases/latest)
[![Platform](https://img.shields.io/badge/platform-linux%20amd64-green.svg)](https://github.com/Antoine-Flo/kubelocal/releases)
[![License](https://img.shields.io/badge/license-MIT-orange.svg)](LICENSE)

## Quick Start

```bash
curl -fsSL https://github.com/Antoine-Flo/kubelocal/releases/download/v0.0.5/kubelocal-v0.0.5-linux-amd64.tar.gz | tar -xz && ./kubelocal
```

## What is kubelocal?

kubelocal is an interactive installer that sets up a complete local Kubernetes development environment on your machine. It handles all the complexity of installing and configuring the necessary tools for local Kubernetes development.
## Features

- 🐳 **Container Runtime Support**: Choose between Docker or Podman
- ☸️ **Kubernetes Distributions**: Install Kind or Minikube
- 🔧 **Complete Setup**: Automatically installs kubectl and configures your cluster
- 💾 **Local Storage**: Sets up persistent storage for your workloads
- 🎯 **Interactive**: User-friendly prompts guide you through the installation
- ⚡ **Fast**: Optimized for quick local development setup

## What Gets Installed

The installer will set up:

1. **Prerequisites** - Essential system packages
2. **Container Runtime** - Docker or Podman (your choice)
3. **kubectl** - Latest version of the Kubernetes command-line tool
4. **Kubernetes Distribution** - Kind or Minikube (your choice)
5. **Cluster Configuration** - Ready-to-use local cluster
6. **Local Storage** - Persistent volume support

## Requirements

- **OS**: Debian-based Linux (tested on Debian)
- **Architecture**: Linux amd64
- **Permissions**: sudo access required for package installation

## Installation Options

### Option 1: Direct Download & Run
```bash
# Download and extract
curl -fsSL https://github.com/Antoine-Flo/kubelocal/releases/download/v0.0.5/kubelocal-v0.0.5-linux-amd64.tar.gz | tar -xz

# Run the installer
./kubelocal
```

### Option 2: Manual Download
```bash
# Download the release
wget https://github.com/Antoine-Flo/kubelocal/releases/download/v0.0.5/kubelocal-v0.0.5-linux-amd64.tar.gz

# Extract
tar -xzf kubelocal-v0.0.5-linux-amd64.tar.gz

# Run
./kubelocal
```

## Usage

Once installed, the script will guide you through:

1. **Container Runtime Selection**: Choose Docker (d) or Podman (p)
2. **Kubernetes Distribution**: Choose Kind (k) or Minikube (m)
3. **Automatic Installation**: Sit back and let the installer do its work

## Example Output

```
🚀 Kubernetes Local Environment Installer

Install Docker (d) or Podman (p) : d
Install Kind (k) or Minikube (m) : k

🔧 Starting installation process...

[1/7] Updating package lists... ✓ Done
[2/7] Installing prerequisites... ✓ Done
[3/7] Installing Docker... ✓ Done
[4/7] Installing kubectl (latest version)... ✓ Done
[5/7] Installing Kind... ✓ Done
[6/7] Configuring cluster... ✓ Done
[7/7] Setting up local storage... ✓ Done

🎉 Installation completed successfully!
Your local Kubernetes environment is ready with Docker and Kind
```

## Version

Check the installed version:
```bash
./kubelocal --version
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/Antoine-Flo/kubelocal/issues)
- 📖 **Documentation**: [GitHub Wiki](https://github.com/Antoine-Flo/kubelocal/wiki)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/Antoine-Flo/kubelocal/discussions)

---

**Made with ❤️ for the Kubernetes community**
