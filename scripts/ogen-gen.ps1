param(
    [string]$OpenApiFiles,
    [string]$YqPath,
    [string]$OgenPath
)
Get-ChildItem -Recurse -Path $OpenApiFiles -Include '*.yaml','*.yml' | ForEach-Object {
    $f = $_.FullName
    if (Select-String -Path $f -Pattern 'x-ogen:' -Quiet) {
        Write-Host "Generating from: $f"
        $target  = & $YqPath e '."x-ogen".target'  $f
        $package = & $YqPath e '."x-ogen".package' $f
        Write-Host "  target=$target  package=$package"
        & $OgenPath --target $target --package $package --clean $f
        if ($LASTEXITCODE -ne 0) { exit 1 }
    }
}
