# Email Service Sender

Golang Fiber service สำหรับรับข้อมูลการชำระเงินและส่งเมลแจ้งเตือน พร้อมระบบ logging

## คุณสมบัติ

-   รับข้อมูลการชำระเงินจากฟอร์ม HTML
-   ส่งเมลแจ้งเตือนการชำระเงิน
-   บันทึก log ทุกการทำงานพร้อม Transaction ID และเวลา
-   API สำหรับตรวจสอบข้อมูลการชำระเงิน
-   Health check endpoint

## การติดตั้ง

1. ติดตั้ง Go dependencies:

```bash
go mod download
```

2. รันเซิร์ฟเวอร์:

```bash
go run main.go
```

เซิร์ฟเวอร์จะทำงานที่ `http://localhost:3000`

## API Endpoints

### POST /api/payment-notification

รับข้อมูลการชำระเงินและส่งเมลแจ้งเตือน

**ตัวอย่างข้อมูล:**

```json
{
    "account": "นาย สมชาย ใจดี",
    "transactionId": "DM13UWT9R4H7IF",
    "transferDate": "2024-01-15",
    "transferTime": "14:30",
    "recipientAccount": "123-0987699",
    "amount": 5000.0,
    "additional": "ชำระค่าสินค้า"
}
```

### GET /api/payment-check/:transactionId

ตรวจสอบข้อมูลการชำระเงินที่มีอยู่

### GET /api/logs/:transactionId

ดู logs ของ transaction ที่ระบุ

### GET /health

Health check endpoint

## ไฟล์ Log

ระบบจะสร้างไฟล์ `payment_logs.txt` เพื่อบันทึก log ทุกการทำงาน โดยจะระบุ:

-   วันที่และเวลา
-   Transaction ID
-   รายละเอียดการทำงาน

## การพัฒนาต่อ

หากต้องการเพิ่มการส่งเมลจริง สามารถใช้ library เช่น:

-   `gopkg.in/mail.v2` สำหรับ SMTP
-   `github.com/sendgrid/sendgrid-go` สำหรับ SendGrid API
-   `github.com/mailgun/mailgun-go` สำหรับ Mailgun API
