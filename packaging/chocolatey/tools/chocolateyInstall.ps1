$ErrorActionPreference = 'Stop'

$packageName = 'agc-cli'
$toolsDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$installDir = Join-Path $toolsDir 'agc-cli'
$zipPath = Join-Path $toolsDir 'agc-cli.zip'

$packageArgs = @{
  packageName   = $packageName
  url64bit      = 'https://github.com/Createitv/agc-cli/releases/download/v0.0.0/agc-cli_0.0.0_windows_amd64.zip'
  checksum64    = '0000000000000000000000000000000000000000000000000000000000000000'
  checksumType64= 'sha256'
  fileFullPath  = $zipPath
}

Get-ChocolateyWebFile @packageArgs
Get-ChocolateyUnzip -FileFullPath $zipPath -Destination $installDir
Install-BinFile -Name agc -Path (Join-Path $installDir 'agc.exe')
