# 🚀 Email Service Sender - Quick Start

## วิธีการรัน

### 1. รันเซิร์ฟเวอร์

```bash
# วิธีที่ 1: รันจาก source code
go run .

# วิธีที่ 2: Build แล้วรัน
go build -o email-service.exe .
.\email-service.exe

# วิธีที่ 3: ใช้ batch file
start.bat
```

### 2. เข้าใช้งาน

-   เปิดเบราว์เซอร์ไปที่: http://localhost:3000
-   Health Check: http://localhost:3000/health
-   API Docs: http://localhost:3000/api/payments

### 3. ทดสอบ API

```bash
# รันไฟล์ทดสอบ
run_test.bat

# หรือทดสอบแยก (ต้องรันเซิร์ฟเวอร์ก่อน)
cd tests
go run test.go
```

## คุณสมบัติหลัก

✅ **ระบบ Logging**: บันทึก log ทุกการทำงานพร้อม Transaction ID และเวลา
✅ **ส่งเมลแจ้งเตือน**: รองรับทั้ง Mock และ Real SMTP  
✅ **ฐานข้อมูล JSON**: เก็บข้อมูลการชำระเงิน
✅ **Web Interface**: หน้าเว็บสำหรับกรอกข้อมูล
✅ **RESTful API**: API endpoints ครบถ้วน

## ไฟล์ที่สำคัญ

-   `main.go` - ไฟล์หลักของเซิร์ฟเวอร์
-   `database.go` - จัดการฐานข้อมูลและ log
-   `email_service.go` - จัดการการส่งเมล
-   `index.html` - หน้าเว็บสำหรับกรอกข้อมูล
-   `payment_logs.txt` - ไฟล์ log (สร้างอัตโนมัติ)
-   `payments.json` - ฐานข้อมูล (สร้างอัตโนมัติ)

## API Endpoints

-   `POST /api/payment-notification` - รับข้อมูลการชำระเงิน
-   `GET /api/payment-check/:id` - ตรวจสอบข้อมูล
-   `GET /api/logs/:id` - ดู logs ของ Transaction ID
-   `GET /api/payments` - ดูรายการทั้งหมด
-   `PUT /api/payment/:id/status` - อัปเดตสถานะ

## Log Example

```
[2024-01-15 14:30:15] Transaction ID: TEST123 - ได้รับข้อมูลการชำระเงิน
[2024-01-15 14:30:16] Transaction ID: TEST123 - รายละเอียด: นาย ทดสอบ ระบบ โอนเงิน 1500.00 บาท เมื่อ 2024-01-15 14:30
[2024-01-15 14:30:17] Transaction ID: TEST123 - บันทึกข้อมูลสำเร็จ
[2024-01-15 14:30:18] Transaction ID: TEST123 - เริ่มส่งเมลแจ้งเตือนการชำระเงิน
[2024-01-15 14:30:19] Transaction ID: TEST123 - ส่งเมลสำเร็จ - จำนวนเงิน: 1500.00 บาท
```

🎉 **พร้อมใช้งานแล้ว!**
