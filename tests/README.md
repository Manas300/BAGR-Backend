# BAGR Backend Tests

This directory contains all test files for the BAGR auction system backend.

## Directory Structure

```
tests/
├── README.md                    # This file
├── run_all_tests.ps1           # Main test runner
├── frontend/                   # Frontend tests
│   ├── index.html              # Visual test interface
│   └── run_frontend_test.ps1   # Frontend test runner
├── integration/                # Integration tests
│   ├── test_email_verification.ps1
│   ├── test_auth_flow.ps1
│   └── test_api_endpoints.ps1
├── unit/                       # Unit tests (Go)
│   ├── auth/
│   ├── models/
│   └── services/
├── e2e/                        # End-to-end tests
│   └── test_complete_flow.ps1
└── scripts/                    # Test helper scripts
    ├── test_db_helpers.ps1
    └── cleanup_test_data.ps1
```

## Test Types

### Integration Tests
- Test complete API flows
- Test database interactions
- Test email verification
- Test authentication flows

### Unit Tests
- Test individual functions
- Test service methods
- Test model validations
- Test utility functions

### End-to-End Tests
- Test complete user journeys
- Test full application workflows
- Test cross-service interactions

## Running Tests

### Frontend Tests
```powershell
# Run visual frontend test
.\tests\frontend\run_frontend_test.ps1

# Or open the HTML file directly
.\tests\frontend\index.html
```

### Integration Tests
```powershell
# Run email verification test
.\tests\integration\test_email_verification.ps1

# Run all integration tests
Get-ChildItem .\tests\integration\*.ps1 | ForEach-Object { & $_.FullName }
```

### Unit Tests
```bash
# Run all unit tests
go test ./...

# Run specific package tests
go test ./internal/auth/...

# Run with coverage
go test -cover ./...
```

### End-to-End Tests
```powershell
# Run complete flow test
.\tests\e2e\test_complete_flow.ps1
```

## Test Data

- Test users are created with `test_` prefix
- Test data is cleaned up after tests
- Use separate test database for integration tests
- Mock external services (email, payment, etc.)

## Best Practices

1. **Clean Code**: Keep tests organized and readable
2. **Isolation**: Each test should be independent
3. **Cleanup**: Always clean up test data
4. **Documentation**: Document what each test does
5. **Naming**: Use descriptive test names
6. **Assertions**: Make clear assertions about expected behavior
