# BAGR Backend Test Runner
# This script runs all tests in the proper order

Write-Host "=== BAGR Backend Test Suite ===" -ForegroundColor Green
Write-Host ""

# Import helper functions
. "$PSScriptRoot\scripts\test_db_helpers.ps1"

# Check prerequisites
Write-Host "Checking prerequisites..." -ForegroundColor Yellow

if (-not (Test-DatabaseConnection)) {
    Write-Host "❌ Database not available. Please start Docker containers first." -ForegroundColor Red
    Write-Host "Run: docker-compose up -d" -ForegroundColor Yellow
    exit 1
}

if (-not (Wait-ForServer)) {
    Write-Host "❌ Server not available. Please start the server first." -ForegroundColor Red
    Write-Host "Run: go run cmd/main.go -config config.yaml" -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ Prerequisites met!" -ForegroundColor Green
Write-Host ""

# Clear any existing test data
Write-Host "Clearing existing test data..." -ForegroundColor Yellow
Clear-TestData

Write-Host ""

# Run integration tests
Write-Host "=== Running Integration Tests ===" -ForegroundColor Cyan

$integrationTests = @(
    "test_auth_flow.ps1",
    "test_email_verification.ps1"
)

foreach ($test in $integrationTests) {
    $testPath = Join-Path $PSScriptRoot "integration\$test"
    if (Test-Path $testPath) {
        Write-Host "Running $test..." -ForegroundColor Yellow
        try {
            & $testPath
            Write-Host "✅ $test completed successfully!" -ForegroundColor Green
        } catch {
            Write-Host "❌ $test failed: $($_.Exception.Message)" -ForegroundColor Red
        }
        Write-Host ""
    } else {
        Write-Host "⚠️  Test file not found: $testPath" -ForegroundColor Yellow
    }
}

# Run unit tests (if Go tests exist)
Write-Host "=== Running Unit Tests ===" -ForegroundColor Cyan

if (Get-Command go -ErrorAction SilentlyContinue) {
    Write-Host "Running Go unit tests..." -ForegroundColor Yellow
    try {
        $unitTestResult = go test ./... -v
        Write-Host "✅ Unit tests completed!" -ForegroundColor Green
    } catch {
        Write-Host "❌ Unit tests failed: $($_.Exception.Message)" -ForegroundColor Red
    }
} else {
    Write-Host "⚠️  Go not found, skipping unit tests" -ForegroundColor Yellow
}

Write-Host ""

# Final cleanup
Write-Host "=== Final Cleanup ===" -ForegroundColor Cyan
Clear-TestData

# Run frontend test
Write-Host "=== Frontend Test Available ===" -ForegroundColor Cyan
Write-Host "To test the complete flow with a visual interface, run:" -ForegroundColor Yellow
Write-Host "  .\tests\frontend\run_frontend_test.ps1" -ForegroundColor White
Write-Host ""

Write-Host "=== All Tests Complete ===" -ForegroundColor Green
Write-Host "Check the output above for any failures." -ForegroundColor Yellow
