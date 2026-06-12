param(
    [string]$BinDir,
    [string]$GofumptPath,
    [string]$GciPath,
    [string]$GofumptVersion,
    [string]$GciVersion
)
$env:GOBIN = $BinDir
if (-not (Test-Path $GofumptPath)) {
    Write-Host "Installing gofumpt $GofumptVersion..."
    go install "mvdan.cc/gofumpt@$GofumptVersion"
    if ($LASTEXITCODE -ne 0) { exit 1 }
}
if (-not (Test-Path $GciPath)) {
    Write-Host "Installing gci $GciVersion..."
    go install "github.com/daixiang0/gci@$GciVersion"
    if ($LASTEXITCODE -ne 0) { exit 1 }
}
Write-Host "Formatters ready"
