@echo off
echo Testing Email Service...
echo.

REM Wait a moment for server to start
timeout /t 3 /nobreak >nul

echo Testing Health Check...
curl -s http://localhost:3000/health
echo.
echo.

echo Testing Payment Notification...
curl -X POST http://localhost:3000/api/payment-notification ^
  -H "Content-Type: application/json" ^
  -d "{\"account\":\"นาย ทดสอบ ระบบ\",\"transactionId\":\"TEST123\",\"transferDate\":\"2024-01-15\",\"transferTime\":\"14:30\",\"recipientAccount\":\"123-456-7890\",\"amount\":1500.00,\"additional\":\"ทดสอบระบบ\"}"
echo.
echo.

echo Testing Logs...
curl -s http://localhost:3000/api/logs/TEST123
echo.
echo.

echo Testing Complete!
pause
