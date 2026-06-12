param(
    [string]$Modules,
    [string]$LintPath
)
$err = 0
Write-Host "Линтим все модули..."
foreach ($m in $Modules.Split(' ')) {
    if (Test-Path $m) {
        Write-Host "  → $m"
        & $LintPath run "$m/..." --config=.golangci.yml
        if ($LASTEXITCODE -ne 0) { $err = 1 }
    }
}
exit $err
