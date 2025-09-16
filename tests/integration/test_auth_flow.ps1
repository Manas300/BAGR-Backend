# Test Complete Authentication Flow
# This script tests all authentication endpoints and flows

Write-Host "=== BAGR Authentication Flow Test ===" -ForegroundColor Green
Write-Host ""

# Test 1: Get available roles
Write-Host "1. Testing Get Roles..." -ForegroundColor Yellow
try {
    $rolesResponse = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/roles" -Method GET
    Write-Host "✅ Roles endpoint working!" -ForegroundColor Green
    Write-Host "Status: $($rolesResponse.StatusCode)"
    Write-Host "Roles: $($rolesResponse.Content)"
} catch {
    Write-Host "❌ Roles endpoint failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# Test 2: Register different user types
Write-Host "2. Testing User Registration (Different Roles)..." -ForegroundColor Yellow

$testUsers = @(
    @{
        email = "producer@example.com"
        username = "producer123"
        first_name = "Producer"
        last_name = "User"
        role = "producer"
    },
    @{
        email = "artist@example.com"
        username = "artist123"
        first_name = "Artist"
        last_name = "User"
        role = "artist"
    },
    @{
        email = "fan@example.com"
        username = "fan123"
        first_name = "Fan"
        last_name = "User"
        role = "fan"
    }
)

foreach ($user in $testUsers) {
    Write-Host "  Registering $($user.role)..." -ForegroundColor Cyan
    $registerBody = @{
        email = $user.email
        username = $user.username
        first_name = $user.first_name
        last_name = $user.last_name
        password = "MySecurePass123"
        confirm_password = "MySecurePass123"
        role = $user.role
    } | ConvertTo-Json

    try {
        $registerResponse = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/register" -Method POST -Body $registerBody -ContentType "application/json"
        Write-Host "    ✅ $($user.role) registered successfully!" -ForegroundColor Green
    } catch {
        Write-Host "    ❌ $($user.role) registration failed: $($_.Exception.Message)" -ForegroundColor Red
    }
}

Write-Host ""

# Test 3: Test password validation
Write-Host "3. Testing Password Validation..." -ForegroundColor Yellow

$weakPasswords = @(
    "password",      # No uppercase
    "PASSWORD",      # No lowercase
    "Password",      # No numbers
    "Password123"    # Common password
)

foreach ($weakPass in $weakPasswords) {
    Write-Host "  Testing weak password: $weakPass" -ForegroundColor Cyan
    $registerBody = @{
        email = "test@example.com"
        username = "testuser"
        first_name = "Test"
        last_name = "User"
        password = $weakPass
        confirm_password = $weakPass
        role = "fan"
    } | ConvertTo-Json

    try {
        $registerResponse = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/register" -Method POST -Body $registerBody -ContentType "application/json"
        Write-Host "    ❌ Weak password should have been rejected!" -ForegroundColor Red
    } catch {
        Write-Host "    ✅ Weak password correctly rejected" -ForegroundColor Green
    }
}

Write-Host ""

# Test 4: Test duplicate email/username
Write-Host "4. Testing Duplicate Registration..." -ForegroundColor Yellow

$duplicateBody = @{
    email = "producer@example.com"  # Already registered
    username = "newuser"
    first_name = "New"
    last_name = "User"
    password = "MySecurePass123"
    confirm_password = "MySecurePass123"
    role = "fan"
} | ConvertTo-Json

try {
    $duplicateResponse = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/register" -Method POST -Body $duplicateBody -ContentType "application/json"
    Write-Host "❌ Duplicate email should have been rejected!" -ForegroundColor Red
} catch {
    Write-Host "✅ Duplicate email correctly rejected" -ForegroundColor Green
}

Write-Host ""

# Test 5: Test forgot password
Write-Host "5. Testing Forgot Password..." -ForegroundColor Yellow

$forgotBody = @{
    email = "producer@example.com"
} | ConvertTo-Json

try {
    $forgotResponse = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/forgot-password" -Method POST -Body $forgotBody -ContentType "application/json"
    Write-Host "✅ Forgot password request successful!" -ForegroundColor Green
    Write-Host "Status: $($forgotResponse.StatusCode)"
} catch {
    Write-Host "❌ Forgot password failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

Write-Host "=== Authentication Flow Test Complete ===" -ForegroundColor Green

