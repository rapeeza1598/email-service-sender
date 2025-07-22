@echo off
echo Installing Go dependencies...
go mod download
go mod tidy

echo.
echo Building the application...
go build -o email-service.exe main.go email_service.go

echo.
echo Build completed! You can now run:
echo email-service.exe
pause
