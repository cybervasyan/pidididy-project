param(
    [string]$BinDir,
    [string]$ToolPath,
    [string]$Package,
    [string]$Version
)
if (Test-Path $ToolPath) {
    exit 0
}
Write-Host "Installing $Package @ $Version..."
$env:GOBIN = $BinDir
$pkg = $Package + '@' + $Version
go install $pkg
if ($LASTEXITCODE -ne 0) { exit 1 }