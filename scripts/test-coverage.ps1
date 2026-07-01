param(
    [string]$Modules,
    [string]$CoverageDir,
    [string]$CoverageFile
)

if (Test-Path $CoverageDir) { Remove-Item -Recurse -Force $CoverageDir }
New-Item -ItemType Directory -Force -Path $CoverageDir | Out-Null

$ERR = 0
foreach ($mod in $Modules -split ' ') {
    Write-Host "Processing module: $mod"

    $pkgs = go list "./$mod/..." 2>$null |
        Where-Object { $_ -match '/(internal/(service|repository))' } |
        Where-Object { $_ -notmatch '/(mocks|testdata|pkg|api|proto|pb|cmd)' }

    if (-not $pkgs) {
        Write-Host "  No suitable packages in $mod, skipping"
        continue
    }

    $pkgStr = $pkgs -join ','
    $outFile = "$CoverageDir/$mod.out"

    go test "-coverpkg=$pkgStr" "-coverprofile=$outFile" -covermode=atomic $pkgs
    if ($LASTEXITCODE -ne 0) { $ERR = 1 }
}

if ($ERR -ne 0) {
    Write-Host "Errors during tests"
    exit $ERR
}

Write-Host "`nCoverage per module:"
foreach ($mod in $Modules -split ' ') {
    $outFile = "$CoverageDir/$mod.out"
    if (Test-Path $outFile) {
        Write-Host -NoNewline " * ${mod}: "
        go tool cover "-func=$outFile" | Select-Object -Last 1
    }
}

Write-Host "`nMerging coverage files..."
$totalFile = "$CoverageDir/$CoverageFile"
# ВАЖНО (Windows PowerShell 5.1): Set-Content/Out-File -Encoding utf8 добавляет BOM,
# из-за которого `go tool cover` падает с "bad mode line". Пишем UTF-8 без BOM через .NET.
$lines = [System.Collections.Generic.List[string]]::new()
$lines.Add("mode: atomic")
Get-ChildItem -Path $CoverageDir -Filter "*.out" |
    Where-Object { $_.Name -ne $CoverageFile } |
    ForEach-Object {
        Get-Content $_.FullName |
            Where-Object { $_ -notmatch "^mode:" } |
            ForEach-Object { $lines.Add($_) }
    }
[System.IO.File]::WriteAllLines($totalFile, $lines)

Write-Host "`nTotal coverage:"
go tool cover "-func=$totalFile" | Select-Object -Last 1
