param(
    [string]$Modules,
    [string]$GofumptPath,
    [string]$GciPath,
    [string]$Prefix
)
Write-Host "Форматируем через gofumpt..."
foreach ($m in $Modules.Split(' ')) {
    if (Test-Path $m) {
        Write-Host "  → $m"
        Get-ChildItem -Recurse -Filter '*.go' -Path $m |
            Where-Object { $_.FullName -notmatch '\\mocks\\' } |
            ForEach-Object { & $GofumptPath -extra -w $_.FullName }
    }
}
Write-Host "Сортируем импорты через gci..."
foreach ($m in $Modules.Split(' ')) {
    if (Test-Path $m) {
        Write-Host "  → $m"
        Get-ChildItem -Recurse -Filter '*.go' -Path $m |
            Where-Object { $_.FullName -notmatch '\\mocks\\' } |
            ForEach-Object { & $GciPath write -s standard -s default -s "prefix($Prefix)" $_.FullName }
    }
}
