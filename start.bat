@echo off
echo Starting Email Service Server...
echo.
echo Make sure you have:
echo 1. Copied .env.example to .env (optional)
echo 2. Updated SMTP settings in .env file (if using real email)
echo 3. Set NOTIFICATION_EMAIL in .env file (if using real email)
echo.

REM Check if compiled executable exists
if exist email-service.exe (
    echo Using compiled executable...
    .\email-service.exe
) else (
    echo Building and running from source...
    go mod tidy
    go run .
)
