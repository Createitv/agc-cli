$ErrorActionPreference = "Stop"

$Repo = "Createitv/agc-cli"
$InstallDir = if ($env:AGC_INSTALL_DIR) { $env:AGC_INSTALL_DIR } else { Join-Path $HOME ".agc\bin" }
$Tag = if ($env:AGC_VERSION) { $env:AGC_VERSION } else { "latest" }

if ($Tag -eq "latest") {
  $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
  $Tag = $Release.tag_name
}

if (-not $Tag) {
  throw "Could not resolve latest agc-cli release."
}

$Version = $Tag.TrimStart("v")
$Archive = "agc-cli_${Version}_windows_amd64.zip"
$BaseUrl = "https://github.com/$Repo/releases/download/$Tag"
$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("agc-cli-" + [System.Guid]::NewGuid())
New-Item -ItemType Directory -Force -Path $TempDir | Out-Null

try {
  $ArchivePath = Join-Path $TempDir $Archive
  $ChecksumsPath = Join-Path $TempDir "checksums.txt"

  Invoke-WebRequest -Uri "$BaseUrl/$Archive" -OutFile $ArchivePath
  Invoke-WebRequest -Uri "$BaseUrl/checksums.txt" -OutFile $ChecksumsPath

  $Expected = (Get-Content $ChecksumsPath | Where-Object { $_ -match "\s+$([regex]::Escape($Archive))$" } | Select-Object -First 1).Split(" ")[0]
  $Actual = (Get-FileHash -Algorithm SHA256 $ArchivePath).Hash.ToLowerInvariant()
  if ($Expected.ToLowerInvariant() -ne $Actual) {
    throw "Checksum mismatch for $Archive."
  }

  Expand-Archive -Path $ArchivePath -DestinationPath $TempDir -Force
  New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
  Copy-Item -Force (Join-Path $TempDir "agc.exe") (Join-Path $InstallDir "agc.exe")

  Write-Host "agc installed to $(Join-Path $InstallDir "agc.exe")"
  Write-Host "Add $InstallDir to PATH if needed."
}
finally {
  Remove-Item -Recurse -Force $TempDir -ErrorAction SilentlyContinue
}
