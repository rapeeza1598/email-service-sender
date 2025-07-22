package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	EmailTimeFormat = "2006-01-02 15:04:05"
)

func getEmailEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type RealEmailService struct {
	Config  SMTPConfig
	LogFile *os.File
}

func NewRealEmailService() (*RealEmailService, error) {
	config := SMTPConfig{
		Host:     getEmailEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:     getEmailEnv("SMTP_PORT", "587"),
		Username: getEmailEnv("SMTP_USERNAME", "your-email@gmail.com"),
		Password: getEmailEnv("SMTP_PASSWORD", "your-app-password"),
		From:     getEmailEnv("SMTP_FROM", "your-email@gmail.com"),
	}

	logFile, err := os.OpenFile("payment_logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถสร้างไฟล์ log ได้: %v", err)
	}

	return &RealEmailService{
		Config:  config,
		LogFile: logFile,
	}, nil
}

func (res *RealEmailService) SendEmail(to, subject, body string) error {
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", to, subject, body)

	auth := smtp.PlainAuth("", res.Config.Username, res.Config.Password, res.Config.Host)

	err := res.sendEmailWithTLS(res.Config.Host+":"+res.Config.Port, auth, res.Config.From, []string{to}, []byte(message))

	return err
}

func (res *RealEmailService) SendEmailWithAttachments(to, subject, body string, attachmentFiles []string) error {
	boundary := "boundary123456789"

	// สร้าง MIME headers
	headers := fmt.Sprintf("To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n\r\n", to, subject, boundary)

	// เริ่มต้น message body
	var message strings.Builder
	message.WriteString(headers)

	// เพิ่ม HTML body
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	message.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
	message.WriteString(body)
	message.WriteString("\r\n")

	// เพิ่ม attachments
	for _, filename := range attachmentFiles {
		filePath := filepath.Join("attachments", filename)

		// อ่านไฟล์
		fileData, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("ไม่สามารถอ่านไฟล์ %s: %v", filePath, err)
			continue
		}

		// กำหนด MIME type
		mimeType := mime.TypeByExtension(filepath.Ext(filename))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		// เพิ่ม attachment header
		message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		message.WriteString(fmt.Sprintf("Content-Type: %s\r\n", mimeType))
		message.WriteString("Content-Transfer-Encoding: base64\r\n")
		message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", filename))

		// เข้ารหัส base64 และเพิ่มลงใน message
		encoded := base64.StdEncoding.EncodeToString(fileData)
		// แบ่งบรรทัดที่ 76 ตัวอักษร (ตาม RFC)
		for i := 0; i < len(encoded); i += 76 {
			end := i + 76
			if end > len(encoded) {
				end = len(encoded)
			}
			message.WriteString(encoded[i:end] + "\r\n")
		}
	}

	// ปิด boundary
	message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	auth := smtp.PlainAuth("", res.Config.Username, res.Config.Password, res.Config.Host)

	err := res.sendEmailWithTLS(res.Config.Host+":"+res.Config.Port, auth, res.Config.From, []string{to}, []byte(message.String()))

	return err
}

func (res *RealEmailService) sendEmailWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.StartTLS(&tls.Config{ServerName: strings.Split(addr, ":")[0]}); err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}

func (res *RealEmailService) SendPaymentNotificationEmail(payment PaymentNotification) error {
	res.WriteLog(payment.TransactionID, "เริ่มส่งเมลแจ้งเตือนการชำระเงิน")

	emailBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2c3e50; color: #f1c40f; padding: 20px; text-align: center; }
        .content { background-color: #f8f9fa; padding: 30px; border: 1px solid #dee2e6; }
        .info-table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
        .info-table th, .info-table td { padding: 12px; text-align: left; border-bottom: 1px solid #dee2e6; }
        .info-table th { background-color: #e9ecef; font-weight: bold; }
        .footer { background-color: #6c757d; color: white; padding: 15px; text-align: center; font-size: 12px; }
        .highlight { background-color: #fff3cd; padding: 10px; border: 1px solid #ffeaa7; margin: 15px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>💳 แจ้งเตือนการชำระเงิน</h1>
        </div>
        
        <div class="content">
            <div class="highlight">
                <strong>🔔 ได้รับการแจ้งชำระเงินใหม่!</strong><br>
                กรุณาตรวจสอบข้อมูลการชำระเงินด้านล่าง
            </div>
            
            <table class="info-table">
                <tr><th>รายการ</th><th>ข้อมูล</th></tr>
                <tr><td><strong>ชื่อบัญชี</strong></td><td>%s</td></tr>
                <tr><td><strong>Transaction ID</strong></td><td><span style="color: #e74c3c; font-weight: bold;">%s</span></td></tr>
                <tr><td><strong>วันที่โอน</strong></td><td>%s</td></tr>
                <tr><td><strong>เวลาที่โอน</strong></td><td>%s</td></tr>
                <tr><td><strong>เลขที่บัญชีที่โอนเข้า</strong></td><td>%s</td></tr>
                <tr><td><strong>จำนวนเงิน</strong></td><td><span style="color: #27ae60; font-weight: bold; font-size: 18px;">%.2f บาท</span></td></tr>
                <tr><td><strong>ข้อมูลเพิ่มเติม</strong></td><td>%s</td></tr>
                <tr><td><strong>ส่งข้อมูลเมื่อ</strong></td><td>%s</td></tr>
            </table>
            
            <div style="margin-top: 30px; padding: 15px; background-color: #d1ecf1; border: 1px solid #bee5eb;">
                <strong>📋 ขั้นตอนต่อไป:</strong><br>
                1. ตรวจสอบข้อมูลการโอนเงิน<br>
                2. ยืนยันการรับเงิน<br>
                3. อัปเดตสถานะคำสั่งซื้อ<br>
                4. แจ้งผลการตรวจสอบให้ลูกค้า
            </div>
        </div>
        
        <div class="footer">
            <p>ระบบแจ้งเตือนอัตโนมัติ | ส่งเมื่อ: %s</p>
            <p>กรุณาตรวจสอบและดำเนินการภายใน 24 ชั่วโมง</p>
        </div>
    </div>
</body>
</html>
`,
		payment.Account,
		payment.TransactionID,
		payment.TransferDate,
		payment.TransferTime,
		payment.RecipientAccount,
		payment.Amount,
		payment.Additional,
		payment.Timestamp.Format(EmailTimeFormat),
		payment.Timestamp.Format(EmailTimeFormat),
	)

	subject := fmt.Sprintf("💳 แจ้งชำระเงิน - Transaction ID: %s", payment.TransactionID)
	toEmail := getEmailEnv("NOTIFICATION_EMAIL", "admin@example.com")

	var err error
	// ตรวจสอบว่ามี attachment files หรือไม่
	if len(payment.AttachmentFiles) > 0 {
		res.WriteLog(payment.TransactionID, fmt.Sprintf("ส่งเมลพร้อม attachment %d ไฟล์: %v", len(payment.AttachmentFiles), payment.AttachmentFiles))
		err = res.SendEmailWithAttachments(toEmail, subject, emailBody, payment.AttachmentFiles)
	} else {
		err = res.SendEmail(toEmail, subject, emailBody)
	}

	if err != nil {
		res.WriteLog(payment.TransactionID, fmt.Sprintf("ส่งเมลไม่สำเร็จ: %v", err))
		return err
	}

	attachmentInfo := ""
	if len(payment.AttachmentFiles) > 0 {
		attachmentInfo = fmt.Sprintf(" (พร้อม %d ไฟล์แนบ)", len(payment.AttachmentFiles))
	}
	res.WriteLog(payment.TransactionID, fmt.Sprintf("ส่งเมลสำเร็จไปยัง: %s - จำนวนเงิน: %.2f บาท%s", toEmail, payment.Amount, attachmentInfo))

	return nil
}

func (res *RealEmailService) WriteLog(transactionID string, message string) {
	timestamp := time.Now().Format(EmailTimeFormat)
	logEntry := fmt.Sprintf("[%s] Transaction ID: %s - %s\n", timestamp, transactionID, message)

	if res.LogFile != nil {
		res.LogFile.WriteString(logEntry)
	}

	log.Print(logEntry)
}

func (res *RealEmailService) Close() {
	if res.LogFile != nil {
		res.LogFile.Close()
	}
}
