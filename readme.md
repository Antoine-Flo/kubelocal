# 🚀 kubelocal

> **One-line installer for local Kubernetes development environment**

[![Release](https://img.shields.io/badge/release-v0.0.5-blue.svg)](https://github.com/Antoine-Flo/kubelocal/releases/latest)
[![Platform](https://img.shields.io/badge/platform-linux%20amd64-green.svg)](https://github.com/Antoine-Flo/kubelocal/releases)
[![License](https://img.shields.io/badge/license-MIT-orange.svg)](LICENSE)

## Quick Start

```bash
curl -fsSL https://kubelocal-website.pages.dev/install | sh
```

## What is kubelocal?

kubelocal is an interactive installer designed for **learning and experimentation** with local Kubernetes development. It simplifies the setup of essential tools needed to experiment with Kubernetes locally, making it easy to get started with hands-on learning.

> ⚠️ **Linux amd64 only** - This tool currently supports only Linux amd64 architecture

> ⚠️ **Learning Tool**: This project is designed for educational purposes and local experimentation. It should be used in VMs or WSL environments for learning Kubernetes concepts. **This is not intended for production environments.**
## Features

- 🐳 **Container Runtime Support**: Choose between Docker or Podman
- ☸️ **Kubernetes Distributions**: Install Kind or Minikube
- 🔧 **kubectl**: Installs the latest version of the Kubernetes command-line tool
- 🛠️ **CLI Tools**: Optional installation of Helm and Kustomize
- 🎯 **Interactive**: User-friendly prompts guide you through the installation
- ⚡ **Fast**: Optimized for quick local learning setup
- 🔗 **kubectl Alias**: Automatically creates a convenient 'k' alias for kubectl

## What Gets Installed

The installer will set up:

1. **kubectl** - Latest version of the Kubernetes command-line tool
2. **Container Runtime** - Docker or Podman (your choice)
3. **Kubernetes Distribution** - Kind or Minikube (your choice)
4. **CLI Tools** - Optional Helm and/or Kustomize (your choice)
5. **kubectl Alias** - Convenient 'k' alias for kubectl commands

> **Note**: This tool installs the necessary tools for learning Kubernetes locally. You'll need to create your first cluster manually using the provided instructions.

## Requirements

- **OS**: Debian-based Linux (tested on Debian)
- **Architecture**: Linux amd64
- **Permissions**: sudo access required for package installation

## Installation Options

### Recommended: Website Installer
```bash
# One-line installer from the official website
curl -fsSL https://kubelocal-website.pages.dev/install | sh
```

### Alternative: Manual Download
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

1. **Container Runtime Selection**: Choose Docker or Podman
2. **Kubernetes Distribution**: Choose Kind or Minikube
3. **CLI Tools Selection**: Choose optional tools like Helm and Kustomize
4. **Automatic Installation**: Sit back and let the installer do its work

## Example Output

```
 ██╗  ██╗██╗   ██╗██████╗ ███████╗██╗      ██████╗  ██████╗ █████╗ ██╗     
 ██║ ██╔╝██║   ██║██╔══██╗██╔════╝██║     ██╔═══██╗██╔════╝██╔══██╗██║     
 █████╔╝ ██║   ██║██████╔╝█████╗  ██║     ██║   ██║██║     ███████║██║     
 ██╔═██╗ ██║   ██║██╔══██╗██╔══╝  ██║     ██║   ██║██║     ██╔══██╗██║     
 ██║  ██╗╚██████╔╝██████╔╝███████╗███████╗╚██████╔╝╚██████╗██║  ██║███████╗
 ╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚══════╝╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝
                                                                           
          🚀 Quick Kubernetes Local Development Setup Tool 🚀

Welcome! Let's set up your local Kubernetes environment.

=== Configuration Summary ===
Container Runtime: docker
Kubernetes Local: kind
Deployment Tools: kubectl, helm

=== Installation in progress ===
✅ Docker installed successfully
✅ Kind installed successfully

🎉 Installation completed successfully!

📝 We created an alias 'k' for kubectl for you!

🚀 Getting Started:
   1. Create your first cluster:
      kind create cluster --name my-cluster

   2. Verify your cluster:
      k cluster-info --context kind-my-cluster
      k get nodes

   3. When you're done:
      kind delete cluster --name my-cluster

💡 Container Runtime (docker) Tips:
   • Check Docker status: docker --version
   • Test Docker access: docker ps
   • If permission denied, run: newgrp docker
   • Or logout/login to apply docker group changes

🔧 Additional Tools:
   • Helm: helm version
     Get started: helm create my-chart

✨ Happy Kubernetes development! ✨
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

## Website & Support

- 🌐 **Website**: [https://kubelocal-website.pages.dev/](https://kubelocal-website.pages.dev/) - Get more information and documentation
- 🐛 **Issues**: [GitHub Issues](https://github.com/Antoine-Flo/kubelocal/issues)
- 📖 **Documentation**: [GitHub Wiki](https://github.com/Antoine-Flo/kubelocal/wiki)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/Antoine-Flo/kubelocal/discussions)

---

**Made with ❤️ for the Kubernetes community**
