@echo off
echo Running Email Service Tests...
echo.

REM Start the server in background
echo Starting server...
start /B ..\email-service.exe

REM Wait for server to start
echo Waiting for server to start...
timeout /t 5 /nobreak >nul

REM Run tests
echo Running tests...
cd tests
go run test.go

REM Kill the server process
echo.
echo Stopping server...
taskkill /F /IM email-service.exe >nul 2>&1

echo.
echo Tests completed!
pause
