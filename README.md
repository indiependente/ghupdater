[![Go Report Card](https://goreportcard.com/badge/github.com/indiependente/ghupdater)](https://goreportcard.com/report/github.com/indiependente/ghupdater)
[![Workflow Status](https://github.com/indiependente/ghupdater/workflows/lint-test/badge.svg)](https://github.com/indiependente/ghupdater/actions)

# ghupdater

A CLI tool to list, download and extract releases downloading the latest version from GitHub.

I currently use it in my homelab to keep my tools up to date. It can also restart a systemd service after the update.

## Installation

```go
go install github.com/indiependente/ghupdater@latest
```

## Usage

### Download and extract

```bash
ghupdater -owner cli -repo cli -os macOS -arch arm64
Get cli/cli latest release
Selected tag v2.64.0
gh_2.64.0_macOS_arm64.zip 100% |█████████████████| (13/13 MB, 38 MB/s)
Extracting asset gh_2.64.0_macOS_arm64.zip to /Users/francesco.f/go/src/github.com/indiependente/ghupdater/
Asset extracted
```

### List available assets

```bash
ghupdater -owner cli -repo cli -list
Get cli/cli latest release
Available assets with filters archive= os= arch=
gh_2.64.0_checksums.txt
gh_2.64.0_linux_386.deb
gh_2.64.0_linux_386.rpm
gh_2.64.0_linux_386.tar.gz
gh_2.64.0_linux_amd64.deb
gh_2.64.0_linux_amd64.rpm
gh_2.64.0_linux_amd64.tar.gz
gh_2.64.0_linux_arm64.deb
gh_2.64.0_linux_arm64.rpm
gh_2.64.0_linux_arm64.tar.gz
gh_2.64.0_linux_armv6.deb
gh_2.64.0_linux_armv6.rpm
gh_2.64.0_linux_armv6.tar.gz
gh_2.64.0_macOS_amd64.zip
gh_2.64.0_macOS_arm64.zip
gh_2.64.0_macOS_universal.pkg
gh_2.64.0_windows_386.msi
gh_2.64.0_windows_386.zip
gh_2.64.0_windows_amd64.msi
gh_2.64.0_windows_amd64.zip
gh_2.64.0_windows_arm64.zip
```

### List available assets with filters

```bash
ghupdater -owner cli -repo cli -arch arm64 -list
Get cli/cli latest release
Available assets with filters archive= os= arch=arm64
gh_2.64.0_linux_arm64.deb
gh_2.64.0_linux_arm64.rpm
gh_2.64.0_linux_arm64.tar.gz
gh_2.64.0_macOS_arm64.zip
gh_2.64.0_windows_arm64.zip
```

## Synopsis

```bash
ghupdater --help
Usage of ghupdater:
  -arch string
    	architecture
  -archive string
    	archive type
  -extract string
    	path to extract archive (default ".")
  -list
    	list available assets
  -os string
    	operating system
  -owner string
    	owner (mandatory)
  -repo string
    	repository (mandatory)
  -restart string
    	unit name to restart systemd service
```

## License

MIT
