# Test Database Helper Functions
# This script provides helper functions for test database operations

function Get-VerificationToken {
    param(
        [string]$Email
    )
    
    $query = "SELECT ev.token FROM email_verifications ev JOIN users u ON ev.user_id = u.id WHERE u.email = '$Email' ORDER BY ev.created_at DESC LIMIT 1;"
    $result = docker exec bagr-postgres psql -U bagr_user -d bagr_db -t -c $query
    return $result.Trim()
}

function Get-ResetToken {
    param(
        [string]$Email
    )
    
    $query = "SELECT pr.token FROM password_resets pr JOIN users u ON pr.user_id = u.id WHERE u.email = '$Email' ORDER BY pr.created_at DESC LIMIT 1;"
    $result = docker exec bagr-postgres psql -U bagr_user -d bagr_db -t -c $query
    return $result.Trim()
}

function Get-UserByEmail {
    param(
        [string]$Email
    )
    
    $query = "SELECT id, email, username, first_name, last_name, role, status, email_verified FROM users WHERE email = '$Email';"
    $result = docker exec bagr-postgres psql -U bagr_user -d bagr_db -c $query
    return $result
}

function Clear-TestData {
    Write-Host "Clearing test data..." -ForegroundColor Yellow
    
    # Clear all test data
    $queries = @(
        "DELETE FROM password_resets;",
        "DELETE FROM email_verifications;",
        "DELETE FROM bids;",
        "DELETE FROM auctions;",
        "DELETE FROM tracks;",
        "DELETE FROM users WHERE email LIKE '%@example.com';"
    )
    
    foreach ($query in $queries) {
        docker exec bagr-postgres psql -U bagr_user -d bagr_db -c $query | Out-Null
    }
    
    Write-Host "✅ Test data cleared!" -ForegroundColor Green
}

function Show-UserStats {
    Write-Host "=== Database User Statistics ===" -ForegroundColor Cyan
    
    $statsQuery = @"
SELECT 
    role,
    COUNT(*) as count,
    COUNT(CASE WHEN email_verified = true THEN 1 END) as verified_count
FROM users 
GROUP BY role
ORDER BY role;
"@
    
    $result = docker exec bagr-postgres psql -U bagr_user -d bagr_db -c $statsQuery
    Write-Host $result
}

function Test-DatabaseConnection {
    Write-Host "Testing database connection..." -ForegroundColor Yellow
    
    try {
        $result = docker exec bagr-postgres psql -U bagr_user -d bagr_db -c "SELECT 1;"
        Write-Host "✅ Database connection successful!" -ForegroundColor Green
        return $true
    } catch {
        Write-Host "❌ Database connection failed: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

function Wait-ForServer {
    param(
        [string]$Url = "http://localhost:8080/health",
        [int]$MaxAttempts = 30,
        [int]$DelaySeconds = 1
    )
    
    Write-Host "Waiting for server to be ready..." -ForegroundColor Yellow
    
    for ($i = 1; $i -le $MaxAttempts; $i++) {
        try {
            $response = Invoke-WebRequest -Uri $Url -Method GET -TimeoutSec 5
            if ($response.StatusCode -eq 200) {
                Write-Host "✅ Server is ready!" -ForegroundColor Green
                return $true
            }
        } catch {
            Write-Host "Attempt $i/$MaxAttempts - Server not ready yet..." -ForegroundColor Yellow
            Start-Sleep -Seconds $DelaySeconds
        }
    }
    
    Write-Host "❌ Server failed to start within $($MaxAttempts * $DelaySeconds) seconds" -ForegroundColor Red
    return $false
}

# Export functions for use in other scripts
Export-ModuleMember -Function Get-VerificationToken, Get-ResetToken, Get-UserByEmail, Clear-TestData, Show-UserStats, Test-DatabaseConnection, Wait-ForServer

