Write-Host "[>>]: Building for production..." -ForegroundColor Cyan
Write-Host "[NPM]: Building `"web > panel`"..." -ForegroundColor Cyan

Set-Location -Path (Join-Path $PWD "web\panel")
& "pnpm" "run" "build"
Write-Host ""

if ($LASTEXITCODE -eq 0) { # Check if the pnpm build was successful
  Write-Host "[NPM]: Web panel built successfully" -ForegroundColor Green
} else {
  Write-Host "[NPM]: build failed." -ForegroundColor Red
}

Set-Location -Path (Join-Path $PWD "..\..")

$platform = Read-Host "Build for Windows or Unix? (Enter 'w' or 'u')" # Ask user for the target platform
$appName = "d2m"

function Build-Windows {
  Write-Host "Building for Windows..." -ForegroundColor Cyan

  $env:GOOS = "windows"
  $env:GOARCH = "amd64"
  go build -o "$appName.exe" ./app/cli
  
  if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful: $appName.exe" -ForegroundColor Green
  } else {
    Write-Host "Build failed" -ForegroundColor Red
  }
}

function Build-Unix {
  Write-Host "Building for Unix..." -ForegroundColor Cyan
  
  $env:GOOS = "linux"
  $env:GOARCH = "amd64"
  go build -o "$appName" ./app/cli

  if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful: $appName" -ForegroundColor Green
  } else {
    Write-Host "Build failed" -ForegroundColor Red
  }
}

switch ($platform.ToLower()) {
  "u" { Build-Unix }
  "w" { Build-Windows }
  default { Write-Host "Invalid input. Please enter 'w' or 'u'" -ForegroundColor Red }
}
