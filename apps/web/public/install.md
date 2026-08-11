# Install agc CLI

Last updated: 2026-08-11

agc CLI is a Go command-line tool for AppGallery Connect automation. It helps developers discover AppGallery Connect API families, bind project profiles, dry-run requests, invoke supported endpoints, export reports, upload release assets, and run a local REST command center.

## Recommended package manager install

These commands install the versioned GitHub Release build instead of compiling from source.

### macOS with Homebrew

```bash
brew tap createitv/tap && brew install agc-cli && agc version
```

The command adds the Createitv tap, installs the released formula, and verifies that `agc` is on PATH.

### Windows with Scoop

Run in Windows PowerShell:

```powershell
if (!(Get-Command scoop -ErrorAction SilentlyContinue)) { Set-ExecutionPolicy RemoteSigned -Scope CurrentUser -Force; irm get.scoop.sh | iex }; scoop bucket add createitv https://github.com/Createitv/scoop-bucket; scoop install agc-cli; agc version
```

Manual Scoop form:

```powershell
scoop bucket add createitv https://github.com/Createitv/scoop-bucket
scoop install agc-cli
agc version
```

### Winget

Winget support has been submitted to Microsoft for review. After approval, the command will be:

```powershell
winget install --id Createitv.AgcCli -e
```

## Go install fallback

If a package manager is not available, Go can still compile agc CLI from the public module.

macOS / Linux:

```bash
mkdir -p "$HOME/.local/bin" && GOBIN="$HOME/.local/bin" go install github.com/Createitv/agc-cli/cmd/agc@latest && export PATH="$HOME/.local/bin:$PATH" && agc version
```

Windows PowerShell:

```powershell
$p="$env:LOCALAPPDATA\Programs\agc\bin"; New-Item -ItemType Directory -Force $p; $env:GOBIN=$p; go install github.com/Createitv/agc-cli/cmd/agc@latest; $env:Path="$p;$env:Path"; agc version
```

Go install requires Go 1.22 or later and may show `dev` in `agc version` because it does not use the GoReleaser version ldflags.

## Package sources

- Homebrew tap: https://github.com/Createitv/homebrew-tap
- Scoop bucket: https://github.com/Createitv/scoop-bucket
- Winget PR: https://github.com/microsoft/winget-pkgs/pull/415361
- GitHub Release: https://github.com/Createitv/agc-cli/releases/tag/v0.1.0

## First commands after install

```bash
agc auth login --service-account-file ~/.agc/service-account.json --name production
agc init --app-id <app-id> --project-id <project-id> --package-name com.example.app --default-profile production
agc capabilities --output table
agc publishing endpoints --output table
```

## Useful links

- Website: https://agccli.app/
- Go package: https://pkg.go.dev/github.com/Createitv/agc-cli/cmd/agc
- Source repository: https://github.com/Createitv/agc-cli
