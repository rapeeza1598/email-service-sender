@echo off
echo ============================================
echo Email Service Sender - Status Check
echo ============================================
echo.

if exist email-service.exe (
    echo ✅ Executable file found: email-service.exe
) else (
    echo ❌ Executable not found. Run: go build -o email-service.exe .
    pause
    exit /b 1
)

if exist main.go (
    echo ✅ Source file found: main.go
) else (
    echo ❌ Source file not found: main.go
)

if exist database.go (
    echo ✅ Database module found: database.go
) else (
    echo ❌ Database module not found: database.go
)

if exist email_service.go (
    echo ✅ Email service found: email_service.go
) else (
    echo ❌ Email service not found: email_service.go
)

if exist index.html (
    echo ✅ Web interface found: index.html
) else (
    echo ❌ Web interface not found: index.html
)

echo.
echo ============================================
echo Ready to run! Use one of these commands:
echo ============================================
echo.
echo 1. .\email-service.exe     (Run compiled version)
echo 2. go run .                (Run from source)
echo 3. start.bat               (Use batch file)
echo.
echo Web interface: http://localhost:3000
echo Health check:  http://localhost:3000/health
echo.
pause
