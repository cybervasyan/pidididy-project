param(
    [string]$Modules
)

$ERR = 0
foreach ($mod in $Modules -split ' ') {
    if (Test-Path $mod) {
        Write-Host "Testing $mod"
        go test -v "./$mod/..."
        if ($LASTEXITCODE -ne 0) { $ERR = 1 }
    }
}
exit $ERR
