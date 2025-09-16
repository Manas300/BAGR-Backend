# Test Email Verification Flow
# This script tests the complete registration -> email verification -> login flow

# Import helper functions
. "$PSScriptRoot\..\scripts\test_db_helpers.ps1"

Write-Host "=== BAGR Authentication Test ===" -ForegroundColor Green
Write-Host ""

# Check prerequisites
if (-not (Test-DatabaseConnection)) {
    Write-Host "❌ Database not available. Please start Docker containers first." -ForegroundColor Red
    exit 1
}

if (-not (Wait-ForServer)) {
    Write-Host "❌ Server not available. Please start the server first." -ForegroundColor Red
    exit 1
}

# Clear any existing test data
Clear-TestData

# Test 1: Register a new user
Write-Host "1. Testing User Registration..." -ForegroundColor Yellow
$registerBody = @{
    email = "testuser@example.com"
    username = "testuser123"
    first_name = "Test"
    last_name = "User"
    password = "MySecurePass123"
    confirm_password = "MySecurePass123"
    role = "producer"
} | ConvertTo-Json

try {
    $registerResponse = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/register" -Method POST -Body $registerBody -ContentType "application/json"
    Write-Host "✅ Registration successful!" -ForegroundColor Green
    Write-Host "Status: $($registerResponse.StatusCode)"
    Write-Host "Response: $($registerResponse.Content)"
} catch {
    Write-Host "❌ Registration failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Write-Host ""

# Test 2: Try to login (should fail - email not verified)
Write-Host "2. Testing Login (should fail - email not verified)..." -ForegroundColor Yellow
$loginBody = @{
    email = "testuser@example.com"
    password = "MySecurePass123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -Body $loginBody -ContentType "application/json"
    Write-Host "❌ Login should have failed!" -ForegroundColor Red
} catch {
    Write-Host "✅ Login correctly blocked (email not verified)" -ForegroundColor Green
    Write-Host "Error: $($_.Exception.Message)"
}

Write-Host ""

# Test 3: Get verification token from database
Write-Host "3. Getting verification token from database..." -ForegroundColor Yellow
$verificationToken = Get-VerificationToken -Email "testuser@example.com"

if ($verificationToken -and $verificationToken -ne "") {
    Write-Host "✅ Found verification token: $verificationToken" -ForegroundColor Green
} else {
    Write-Host "❌ No verification token found!" -ForegroundColor Red
    exit 1
}

Write-Host ""

# Test 4: Verify email
Write-Host "4. Testing Email Verification..." -ForegroundColor Yellow
$verifyUrl = "http://localhost:8080/api/v1/auth/verify?token=$verificationToken"

try {
    $verifyResponse = Invoke-WebRequest -Uri $verifyUrl -Method GET
    Write-Host "✅ Email verification successful!" -ForegroundColor Green
    Write-Host "Status: $($verifyResponse.StatusCode)"
} catch {
    Write-Host "❌ Email verification failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# Test 5: Login after verification (should succeed)
Write-Host "5. Testing Login after verification..." -ForegroundColor Yellow

try {
    $loginResponse = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -Body $loginBody -ContentType "application/json"
    Write-Host "✅ Login successful after verification!" -ForegroundColor Green
    Write-Host "Status: $($loginResponse.StatusCode)"
    Write-Host "Response: $($loginResponse.Content)"
    
    # Extract tokens from response
    $responseData = $loginResponse.Content | ConvertFrom-Json
    $accessToken = $responseData.data.access_token
    Write-Host "Access Token: $($accessToken.Substring(0, 50))..." -ForegroundColor Cyan
} catch {
    Write-Host "❌ Login failed after verification: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# Test 6: Test protected endpoint
Write-Host "6. Testing Protected Endpoint..." -ForegroundColor Yellow
if ($accessToken) {
    try {
        $headers = @{
            "Authorization" = "Bearer $accessToken"
        }
        $profileResponse = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/profile" -Method GET -Headers $headers
        Write-Host "✅ Protected endpoint access successful!" -ForegroundColor Green
        Write-Host "Status: $($profileResponse.StatusCode)"
        Write-Host "Profile: $($profileResponse.Content)"
    } catch {
        Write-Host "❌ Protected endpoint access failed: $($_.Exception.Message)" -ForegroundColor Red
    }
} else {
    Write-Host "⚠️  Skipping protected endpoint test (no access token)" -ForegroundColor Yellow
}

Write-Host ""

# Test 7: Show final statistics
Write-Host "7. Final Database Statistics..." -ForegroundColor Yellow
Show-UserStats

Write-Host ""

# Cleanup
Write-Host "8. Cleaning up test data..." -ForegroundColor Yellow
Clear-TestData

Write-Host ""
Write-Host "=== Test Complete ===" -ForegroundColor Green