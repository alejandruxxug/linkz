<#
.SYNOPSIS
  Install Linkz on Windows (amd64) from a GitHub release of alejandruxxug/linkz.

.DESCRIPTION
  .\install.ps1                        latest release -> %LOCALAPPDATA%\Programs\Linkz
  .\install.ps1 -Version v0.1.0        a specific release
  .\install.ps1 -Dir C:\Tools\Linkz    another folder
  .\install.ps1 -FromFile linkz_0.1.0_windows_amd64.zip [-Checksums checksums.txt]

  Downloads come only from github.com/alejandruxxug/linkz releases. While that
  repo is private this works once a release is public, or with a token that can
  read the repo:  $env:GITHUB_TOKEN = "..."; .\install.ps1

  -FromFile installs an archive you already have (no network). The archive is
  checked when -Checksums is given.
#>
param(
  [string]$Version = "",
  [string]$Dir = (Join-Path $env:LOCALAPPDATA "Programs\Linkz"),
  [string]$FromFile = "",
  [string]$Checksums = ""
)
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$Repo = "alejandruxxug/linkz"
$Bin = "linkz.exe"

function Fail($msg) { Write-Error "install: $msg"; exit 1 }

if ($env:PROCESSOR_ARCHITECTURE -ne "AMD64") { Fail "no Linkz build for $env:PROCESSOR_ARCHITECTURE" }

$tmp = Join-Path ([System.IO.Path]::GetTempPath()) ("linkz-" + [guid]::NewGuid())
New-Item -ItemType Directory -Path $tmp | Out-Null

try {
  # Check ARCHIVE (named NAME) against the sha256 list in CHECKSUMS.
  function Test-Archive($archive, $name, $sums) {
    $line = Get-Content $sums | Where-Object { ($_ -split "\s+")[1] -eq $name } | Select-Object -First 1
    if (-not $line) { Fail "$name is not listed in the checksums" }
    $want = ($line -split "\s+")[0].ToLower()
    $got = (Get-FileHash -Algorithm SHA256 $archive).Hash.ToLower()
    if ($want -ne $got) { Fail "checksum mismatch for $name (wanted $want, got $got)" }
    Write-Host "checksum ok"
  }

  if ($FromFile) {
    if (-not (Test-Path $FromFile)) { Fail "no such file: $FromFile" }
    $name = Split-Path $FromFile -Leaf
    $archive = Join-Path $tmp $name
    Copy-Item $FromFile $archive
    if ($Checksums) { Test-Archive $archive $name $Checksums }
    else { Write-Warning "no -Checksums given, the archive is not verified" }
  } else {
    $headers = @{ Accept = "application/vnd.github+json" }
    if ($env:GITHUB_TOKEN) { $headers["Authorization"] = "Bearer $env:GITHUB_TOKEN" }
    $api = "https://api.github.com/repos/$Repo/releases"
    $rel = if ($Version) { "$api/tags/$Version" } else { "$api/latest" }
    try { $release = Invoke-RestMethod -Uri $rel -Headers $headers }
    catch { Fail "could not read the release ($rel). The repo is private until a release is public; set `$env:GITHUB_TOKEN to a token that can read it." }
    $tag = $release.tag_name
    $ver = $tag.TrimStart("v")
    $name = "linkz_${ver}_windows_amd64.zip"
    Write-Host "installing linkz $tag for windows/amd64"

    function Get-Asset($assetName, $dest) {
      $asset = $release.assets | Where-Object { $_.name -eq $assetName } | Select-Object -First 1
      if (-not $asset) { Fail "release $tag has no asset $assetName" }
      $h = @{ Accept = "application/octet-stream" }
      if ($env:GITHUB_TOKEN) { $h["Authorization"] = "Bearer $env:GITHUB_TOKEN" }
      Invoke-WebRequest -Uri $asset.url -Headers $h -OutFile $dest
    }
    $archive = Join-Path $tmp $name
    Get-Asset $name $archive
    Get-Asset "checksums.txt" (Join-Path $tmp "checksums.txt")
    Test-Archive $archive $name (Join-Path $tmp "checksums.txt")
  }

  $x = Join-Path $tmp "x"
  Expand-Archive -Path $archive -DestinationPath $x
  $exe = Join-Path $x $Bin
  if (-not (Test-Path $exe)) { Fail "the archive has no $Bin" }
  New-Item -ItemType Directory -Force -Path $Dir | Out-Null
  Copy-Item $exe (Join-Path $Dir $Bin) -Force
  Write-Host "installed $(Join-Path $Dir $Bin)"

  Write-Host ""
  Write-Host "Next steps"
  $onPath = ($env:PATH -split ";") -contains $Dir
  if (-not $onPath) {
    Write-Host "  - add it to your PATH (new terminals pick it up):"
    Write-Host "      [Environment]::SetEnvironmentVariable('PATH', [Environment]::GetEnvironmentVariable('PATH','User') + ';$Dir', 'User')"
  }
  Write-Host "  - try it:  linkz render examples\hello.lzdraw   (examples are in the archive)"
  Write-Host "  - open the editor:  linkz edit"
  Write-Host "  - if SmartScreen warns (the app is not signed yet): More info > Run anyway"
}
finally {
  Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
}
