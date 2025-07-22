# 🚀 Email Service Sender - Golang Fiber

## คำอธิบาย

โปรเจกต์นี้เป็น **Golang Fiber service** สำหรับการรับข้อมูลการชำระเงินและส่งเมลแจ้งเตือนพร้อมระบบ logging ที่สมบูรณ์

### ✨ คุณสมบัติหลัก

-   📧 **ส่งเมลแจ้งเตือน**: รองรับทั้ง Mock และ Real SMTP
-   📝 **ระบบ Logging**: บันทึก log ทุกการทำงานพร้อม Transaction ID และเวลา
-   💾 **ฐานข้อมูล JSON**: เก็บข้อมูลการชำระเงินในไฟล์ JSON
-   🌐 **RESTful API**: API endpoints ครบถ้วนสำหรับจัดการข้อมูล
-   🎨 **Web Interface**: หน้าเว็บสำหรับกรอกข้อมูลการชำระเงิน

## 📁 โครงสร้างไฟล์

```
📦 email-service-sender/
├── 📄 main.go                 # ไฟล์หลักของเซิร์ฟเวอร์
├── 📄 database.go             # จัดการฐานข้อมูล JSON และ Log
├── 📄 email_service.go        # จัดการการส่งเมล (Real SMTP)
├── 📄 index.html              # หน้าเว็บสำหรับกรอกข้อมูล
├── 📄 go.mod                  # Go module dependencies
├── 📄 .env.example            # ตัวอย่างการตั้งค่า environment
├── 📄 build.bat               # สคริปต์สำหรับ build บน Windows
├── 📄 start.bat               # สคริปต์สำหรับเริ่มเซิร์ฟเวอร์
├── 📄 test.go                 # ไฟล์ทดสอบ API
└── 📄 README.md               # คู่มือนี้
```

## 🚀 การติดตั้งและเริ่มใช้งาน

### 1. ติดตั้ง Dependencies

```bash
go mod download
```

### 2. ตั้งค่า Environment Variables

```bash
# คัดลอกไฟล์ตัวอย่าง
copy .env.example .env

# แก้ไขไฟล์ .env ตามต้องการ
```

### 3. เริ่มเซิร์ฟเวอร์

**วิธีที่ 1: ใช้ Go command**

```bash
go run .
```

**วิธีที่ 2: ใช้ batch file (Windows)**

```bash
start.bat
```

เซิร์ฟเวอร์จะทำงานที่: `http://localhost:3000`

## 🔧 การตั้งค่า

### Environment Variables

| ตัวแปร               | ค่าเริ่มต้น         | อธิบาย                     |
| -------------------- | ------------------- | -------------------------- |
| `PORT`               | `3000`              | พอร์ตที่เซิร์ฟเวอร์จะทำงาน |
| `USE_REAL_EMAIL`     | `false`             | ใช้งาน Real Email Service  |
| `SMTP_HOST`          | `smtp.gmail.com`    | SMTP Server Host           |
| `SMTP_PORT`          | `587`               | SMTP Server Port           |
| `SMTP_USERNAME`      | -                   | อีเมลสำหรับส่ง             |
| `SMTP_PASSWORD`      | -                   | รหัสผ่านหรือ App Password  |
| `SMTP_FROM`          | -                   | อีเมลผู้ส่ง                |
| `NOTIFICATION_EMAIL` | `admin@example.com` | อีเมลที่จะรับการแจ้งเตือน  |

### การตั้งค่า Gmail SMTP

1. เปิดใช้งาน 2-Factor Authentication
2. สร้าง App Password สำหรับแอปพลิเคชัน
3. ใช้ App Password แทนรหัสผ่านปกติ

## 📡 API Endpoints

### 🏥 Health Check

```http
GET /health
```

ตรวจสอบสถานะเซิร์ฟเวอร์

### 💳 รับข้อมูลการชำระเงิน

```http
POST /api/payment-notification
Content-Type: application/json

{
  "account": "นาย สมชาย ใจดี",
  "transactionId": "DM13UWT9R4H7IF",
  "transferDate": "2024-01-15",
  "transferTime": "14:30",
  "recipientAccount": "123-0987699",
  "amount": 5000.00,
  "additional": "ชำระค่าสินค้า"
}
```

### 🔍 ตรวจสอบข้อมูลการชำระเงิน

```http
GET /api/payment-check/{transactionId}
```

### 📋 ดู Logs ของ Transaction

```http
GET /api/logs/{transactionId}
```

### 📊 ดูรายการการชำระเงินทั้งหมด

```http
GET /api/payments
GET /api/payments?status=รอตรวจสอบ
```

### ✏️ อัปเดตสถานะการชำระเงิน

```http
PUT /api/payment/{transactionId}/status
Content-Type: application/json

{
  "status": "ยืนยันแล้ว"
}
```

## 📝 ระบบ Logging

### การทำงานของ Log System

1. **Log File**: `payment_logs.txt`
2. **รูปแบบ Log**: `[YYYY-MM-DD HH:MM:SS] Transaction ID: {ID} - {Message}`
3. **การบันทึก**: บันทึกทั้งในไฟล์และ console

### ตัวอย่าง Log

```
[2024-01-15 14:30:15] Transaction ID: DM13UWT9R4H7IF - ได้รับข้อมูลการชำระเงิน
[2024-01-15 14:30:16] Transaction ID: DM13UWT9R4H7IF - รายละเอียด: นาย สมชาย ใจดี โอนเงิน 5000.00 บาท เมื่อ 2024-01-15 14:30
[2024-01-15 14:30:17] Transaction ID: DM13UWT9R4H7IF - บันทึกข้อมูลสำเร็จ
[2024-01-15 14:30:18] Transaction ID: DM13UWT9R4H7IF - เริ่มส่งเมลแจ้งเตือนการชำระเงิน
[2024-01-15 14:30:19] Transaction ID: DM13UWT9R4H7IF - ส่งเมลสำเร็จ - จำนวนเงิน: 5000.00 บาท
```

## 🗄️ ฐานข้อมูล

โปรเจกต์ใช้ JSON file เป็นฐานข้อมูล:

-   **ไฟล์**: `payments.json`
-   **โครงสร้าง**: เก็บข้อมูลการชำระเงินพร้อมสถานะและเวลา

### ตัวอย่างข้อมูลในฐานข้อมูล

```json
{
    "payments": {
        "DM13UWT9R4H7IF": {
            "account": "นาย สมชาย ใจดี",
            "transactionId": "DM13UWT9R4H7IF",
            "transferDate": "2024-01-15",
            "transferTime": "14:30",
            "recipientAccount": "123-0987699",
            "amount": 5000.0,
            "additional": "ชำระค่าสินค้า",
            "timestamp": "2024-01-15T14:30:15Z",
            "status": "รอตรวจสอบ",
            "submitted_at": "2024-01-15T14:30:15Z",
            "updated_at": "2024-01-15T14:30:15Z"
        }
    },
    "saved_at": "2024-01-15T14:30:15Z"
}
```

## 📧 การส่งเมล

### Mock Email Service (เริ่มต้น)

-   จำลองการส่งเมล
-   เขียน log เท่านั้น
-   ไม่ส่งเมลจริง

### Real Email Service

-   ส่งเมลจริงผ่าน SMTP
-   รองรับ TLS/SSL
-   เนื้อหาเมลแบบ HTML

## 🧪 การทดสอบ

### ทดสอบด้วย cURL

```bash
# ทดสอบ Health Check
curl http://localhost:3000/health

# ทดสอบส่งข้อมูลการชำระเงิน
curl -X POST http://localhost:3000/api/payment-notification \
  -H "Content-Type: application/json" \
  -d '{
    "account": "นาย ทดสอบ ระบบ",
    "transactionId": "TEST123456789",
    "transferDate": "2024-01-15",
    "transferTime": "14:30",
    "recipientAccount": "123-456-7890",
    "amount": 1500.00,
    "additional": "ทดสอบระบบ"
  }'

# ทดสอบดู logs
curl http://localhost:3000/api/logs/TEST123456789
```

### ทดสอบด้วยไฟล์ test.go

```bash
go run test.go
```

## 🌐 หน้าเว็บ

เข้าใช้งานผ่าน: `http://localhost:3000`

### คุณสมบัติของหน้าเว็บ:

-   📝 กรอกข้อมูลการชำระเงิน
-   📎 อัปโหลดหลักฐานการโอน
-   🔍 ตรวจสอบข้อมูลที่มีอยู่แล้ว
-   ✅ แสดงผลการส่งข้อมูล

## 🔧 การพัฒนาต่อ

### การเพิ่มฟีเจอร์ใหม่

1. **ฐานข้อมูลจริง**: MySQL, PostgreSQL
2. **Authentication**: JWT, OAuth
3. **File Upload**: รองรับไฟล์หลักฐาน
4. **Webhook**: แจ้งเตือนไปยังระบบอื่น
5. **Dashboard**: หน้าจัดการข้อมูล

### การติดตั้งเพิ่มเติม

```bash
# สำหรับ MySQL
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql

# สำหรับ JWT
go get -u github.com/golang-jwt/jwt/v5

# สำหรับ File Upload
go get -u github.com/gofiber/fiber/v2/middleware/filesystem
```

## 🐛 การแก้ไขปัญหา

### ปัญหาที่พบบ่อย

1. **ไม่สามารถส่งเมลได้**

    - ตรวจสอบการตั้งค่า SMTP
    - ตรวจสอบ App Password
    - ตรวจสอบการเชื่อมต่ออินเทอร์เน็ต

2. **Port ถูกใช้งานแล้ว**

    - เปลี่ยน PORT ในไฟล์ .env
    - หรือปิดโปรแกรมที่ใช้ port นั้น

3. **ไม่สามารถเขียนไฟล์ได้**
    - ตรวจสอบสิทธิ์การเขียนไฟล์
    - รันในโหมด Administrator (Windows)

## 📄 License

โปรเจกต์นี้เป็น Open Source สามารถนำไปใช้งานและพัฒนาต่อได้อย่างอิสระ

## 👨‍💻 ผู้พัฒนา

พัฒนาโดย GitHub Copilot สำหรับการจัดการระบบแจ้งชำระเงินและส่งเมลแจ้งเตือน

---

🎉 **ขอให้มีความสุขกับการใช้งาน Email Service Sender!** 🎉
