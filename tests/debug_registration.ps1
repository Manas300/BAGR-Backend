# Debug Registration Test
Write-Host "Testing Registration API..." -ForegroundColor Green

$testData = @{
    email = "test@example.com"
    username = "testuser123"
    first_name = "Test"
    last_name = "User"
    password = "TestPass123"
    confirm_password = "TestPass123"
    role = "fan"
} | ConvertTo-Json

Write-Host "Sending registration request..." -ForegroundColor Yellow
Write-Host "Data: $testData" -ForegroundColor Gray

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/register" -Method POST -Body $testData -ContentType "application/json"
    
    Write-Host "Status Code: $($response.StatusCode)" -ForegroundColor Green
    Write-Host "Response:" -ForegroundColor Green
    Write-Host $response.Content -ForegroundColor White
}
catch {
    Write-Host "Error occurred:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response Body:" -ForegroundColor Red
        Write-Host $responseBody -ForegroundColor White
    }
}

