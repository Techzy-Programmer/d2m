Write-Host "[NOTE]: Running in DEBUG mode not suitable for production" -ForegroundColor DarkYellow
Write-Host "[WARN]: Web assets are not built automagically!!" -ForegroundColor Yellow

# ENV for Go build (Windows)
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$appName = "d2m"

$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) "go-build-$appName-debug"

if (-Not (Test-Path -Path $tempDir)) {
  New-Item -ItemType Directory -Path $tempDir
}

$exePath = Join-Path $tempDir "$appName.exe"

Write-Host "[GO]: Building the Go application..." -ForegroundColor Cyan
# Build the Go application into the temp dir with debug flags
go build -gcflags="all=-N -l" -ldflags "-X 'main.Release=test'" -o $exePath ./app/cli # Useful for debugging with "Delve"

if ($LASTEXITCODE -eq 0) { # Check if the go build was successful
  Write-Host "[GO]: Binary built successfully targeting `"$env:GOOS-$env:GOARCH`"" -ForegroundColor Green
  Start-Sleep -Seconds 2
  Clear-Host
  Write-Host "[LAUNCH]: `"$appName.exe`" from `"$tempDir`"`n" -ForegroundColor DarkGray
  & $exePath
} else {
  Write-Host "[GO]: build failed." -ForegroundColor Red
}

# Optional Cleanup: remove the temp dir
# Remove-Item -Recurse -Force $tempDir
