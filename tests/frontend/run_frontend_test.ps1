# Frontend Test Runner
# This script opens the HTML frontend for testing the authentication flow

Write-Host "=== BAGR Frontend Test ===" -ForegroundColor Green
Write-Host ""

# Check if server is running
Write-Host "Checking if server is running..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method GET -TimeoutSec 5
    if ($response.StatusCode -eq 200) {
        Write-Host "✅ Server is running!" -ForegroundColor Green
    } else {
        Write-Host "❌ Server returned status: $($response.StatusCode)" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "❌ Server is not running. Please start the server first:" -ForegroundColor Red
    Write-Host "   go run cmd/main.go -config config.yaml" -ForegroundColor Yellow
    exit 1
}

Write-Host ""

# Get the full path to the HTML file
$htmlPath = Join-Path $PSScriptRoot "index.html"
$fullPath = Resolve-Path $htmlPath

Write-Host "Opening frontend test page..." -ForegroundColor Yellow
Write-Host "Path: $fullPath" -ForegroundColor Cyan
Write-Host ""

# Open the HTML file in default browser
try {
    Start-Process $fullPath
    Write-Host "✅ Frontend test page opened in browser!" -ForegroundColor Green
} catch {
    Write-Host "❌ Failed to open browser. Please manually open:" -ForegroundColor Red
    Write-Host "   $fullPath" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=== Test Instructions ===" -ForegroundColor Cyan
Write-Host "1. Choose your role (Producer, Artist, Fan, or Moderator)" -ForegroundColor White
Write-Host "2. Fill out the registration form" -ForegroundColor White
Write-Host "3. Check the server console for the verification token" -ForegroundColor White
Write-Host "4. Paste the token to verify your email" -ForegroundColor White
Write-Host "5. Login with your credentials" -ForegroundColor White
Write-Host "6. Test the protected endpoint" -ForegroundColor White
Write-Host ""
Write-Host "Watch the server console for verification tokens!" -ForegroundColor Yellow
Write-Host "Press Ctrl+C to stop this script" -ForegroundColor Gray

# Keep the script running so user can see the instructions
try {
    while ($true) {
        Start-Sleep -Seconds 1
    }
} catch {
    Write-Host "`nFrontend test completed!" -ForegroundColor Green
}



