param([string]$Modules)
Write-Host "Updating dependencies..."
foreach ($m in $Modules.Split(' ')) {
    if (Test-Path $m) {
        Write-Host "  -> $m"
        Push-Location $m
        go mod tidy
        if ($LASTEXITCODE -ne 0) { Pop-Location; exit 1 }
        Pop-Location
    }
}
