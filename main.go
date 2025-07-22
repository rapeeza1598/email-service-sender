package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

const (
	TimeFormat  = "2006-01-02 15:04:05"
	LogFileName = "payment_logs.txt"
)

type PaymentNotification struct {
	Account          string    `json:"account" form:"account"`
	TransactionID    string    `json:"transactionId" form:"transactionId"`
	TransferDate     string    `json:"transferDate" form:"transferDate"`
	TransferTime     string    `json:"transferTime" form:"transferTime"`
	RecipientAccount string    `json:"recipientAccount" form:"recipientAccount"`
	Amount           float64   `json:"amount" form:"amount"`
	Additional       string    `json:"additional" form:"additional"`
	Timestamp        time.Time `json:"timestamp"`
	AttachmentFiles  []string  `json:"attachmentFiles,omitempty"` // รายชื่อไฟล์ที่แนบมา
}

type EmailService struct {
	LogFile *os.File
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func NewEmailService() (*EmailService, error) {
	logFile, err := os.OpenFile(LogFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถสร้างไฟล์ log ได้: %v", err)
	}

	return &EmailService{
		LogFile: logFile,
	}, nil
}

func (es *EmailService) Close() {
	if es.LogFile != nil {
		es.LogFile.Close()
	}
}

func (es *EmailService) WriteLog(transactionID string, message string) {
	timestamp := time.Now().Format(TimeFormat)
	// เขียน log แบบง่าย: เฉพาะ Transaction ID และ timestamp
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, transactionID)

	if es.LogFile != nil {
		es.LogFile.WriteString(logEntry)
	}

	log.Printf("[%s] %s", timestamp, transactionID)
}

// ฟังก์ชันทำความสะอาด log เก่าที่เกิน 5 วัน
func (es *EmailService) CleanOldLogs() error {
	// อ่านไฟล์ log ทั้งหมด
	content, err := os.ReadFile(LogFileName)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var validLines []string
	cutoffDate := time.Now().AddDate(0, 0, -5) // 5 วันที่แล้ว

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// แยก timestamp จาก log line
		if strings.HasPrefix(line, "[") {
			endBracket := strings.Index(line, "]")
			if endBracket > 0 {
				timestampStr := line[1:endBracket]
				logTime, err := time.Parse(TimeFormat, timestampStr)
				if err == nil && logTime.After(cutoffDate) {
					validLines = append(validLines, line)
				}
			}
		}
	}

	// เขียนกลับไฟล์ด้วย log ที่ยังไม่เก่าเกินไป
	newContent := strings.Join(validLines, "\n")
	if len(validLines) > 0 {
		newContent += "\n"
	}

	return os.WriteFile(LogFileName, []byte(newContent), 0666)
}

func (es *EmailService) SendPaymentNotificationEmail(payment PaymentNotification) error {
	es.WriteLog(payment.TransactionID, "")

	time.Sleep(1 * time.Second)

	es.WriteLog(payment.TransactionID, "")

	return nil
}

func setupRoutes(app *fiber.App, db *PaymentDatabase, logManager *LogManager, emailSender interface {
	SendPaymentNotificationEmail(PaymentNotification) error
	WriteLog(string, string)
}) {
	// Static files - เสิร์ฟไฟล์ HTML
	app.Static("/", "./")

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now(),
			"service":   "email-service-sender",
		})
	})

	// API endpoint สำหรับอัปโหลดไฟล์ attachment
	app.Post("/api/upload-attachment", func(c *fiber.Ctx) error {
		// สร้างโฟลเดอร์ attachments หากยังไม่มี
		os.MkdirAll("attachments", 0755)

		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "ไม่สามารถอ่านไฟล์ได้",
			})
		}

		files := form.File["files"]
		var uploadedFiles []string

		for _, file := range files {
			// สร้างชื่อไฟล์ใหม่เพื่อหลีกเลี่ยงการชนกัน
			timestamp := time.Now().Format("20060102_150405")
			filename := fmt.Sprintf("%s_%s", timestamp, file.Filename)
			filepath := fmt.Sprintf("attachments/%s", filename)

			// บันทึกไฟล์
			if err := c.SaveFile(file, filepath); err != nil {
				return c.Status(500).JSON(fiber.Map{
					"success": false,
					"message": fmt.Sprintf("ไม่สามารถบันทึกไฟล์ %s ได้", file.Filename),
				})
			}

			uploadedFiles = append(uploadedFiles, filename)
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "อัปโหลดไฟล์สำเร็จ",
			"files":   uploadedFiles,
		})
	})

	// API endpoint สำหรับรับข้อมูลการชำระเงิน
	app.Post("/api/payment-notification", func(c *fiber.Ctx) error {
		var payment PaymentNotification

		// Parse ข้อมูลจาก form
		if err := c.BodyParser(&payment); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "ข้อมูลไม่ถูกต้อง",
			})
		}

		// ตั้งค่า timestamp
		payment.Timestamp = time.Now()

		// เขียน log เฉพาะ Transaction ID
		logManager.WriteLog(payment.TransactionID, "")

		// ตรวจสอบข้อมูลที่จำเป็น
		if payment.TransactionID == "" || payment.Account == "" || payment.Amount <= 0 {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "กรุณากรอกข้อมูลให้ครบถ้วน",
			})
		}

		// บันทึกข้อมูลลงฐานข้อมูล
		if err := db.AddPayment(payment); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "ไม่สามารถบันทึกข้อมูลได้",
			})
		}

		// ส่งเมลแจ้งเตือน
		if err := emailSender.SendPaymentNotificationEmail(payment); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "ไม่สามารถส่งเมลแจ้งเตือนได้",
			})
		}

		// ส่งผลลัพธ์กลับ
		return c.JSON(fiber.Map{
			"success": true,
			"message": "ได้รับข้อมูลการชำระเงินเรียบร้อยแล้ว เราจะตรวจสอบและแจ้งผลภายใน 24 ชั่วโมง",
			"data": fiber.Map{
				"transactionId": payment.TransactionID,
				"timestamp":     payment.Timestamp,
			},
		})
	})

	// API endpoint สำหรับตรวจสอบข้อมูลการชำระเงินที่มีอยู่
	app.Get("/api/payment-check/:transactionId", func(c *fiber.Ctx) error {
		transactionId := c.Params("transactionId")

		// ตรวจสอบข้อมูลจากฐานข้อมูล
		record, exists := db.GetPayment(transactionId)

		if exists {
			return c.JSON(fiber.Map{
				"exists": true,
				"paymentData": fiber.Map{
					"transactionId":    record.TransactionID,
					"account":          record.Account,
					"transferDate":     record.TransferDate,
					"transferTime":     record.TransferTime,
					"recipientAccount": record.RecipientAccount,
					"amount":           record.Amount,
					"additional":       record.Additional,
					"status":           record.Status,
					"submittedDate":    record.SubmittedAt.Format(TimeFormat),
				},
			})
		}

		return c.JSON(fiber.Map{
			"exists": false,
		})
	})

	// API endpoint สำหรับดู logs
	app.Get("/api/logs/:transactionId", func(c *fiber.Ctx) error {
		transactionId := c.Params("transactionId")

		// อ่าน logs จาก LogManager
		logs, err := logManager.GetLogs(transactionId)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "ไม่สามารถอ่าน logs ได้",
			})
		}

		return c.JSON(fiber.Map{
			"success":       true,
			"transactionId": transactionId,
			"logs":          logs,
		})
	})

	// API endpoint สำหรับดูรายการการชำระเงินทั้งหมด
	app.Get("/api/payments", func(c *fiber.Ctx) error {
		status := c.Query("status") // กรองตามสถานะ

		var payments map[string]PaymentRecord
		if status != "" {
			payments = db.GetPaymentsByStatus(status)
		} else {
			payments = db.GetAllPayments()
		}

		return c.JSON(fiber.Map{
			"success":  true,
			"payments": payments,
			"count":    len(payments),
		})
	})

	// API endpoint สำหรับอัปเดตสถานะการชำระเงิน
	app.Put("/api/payment/:transactionId/status", func(c *fiber.Ctx) error {
		transactionId := c.Params("transactionId")

		var body struct {
			Status string `json:"status"`
		}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "ข้อมูลไม่ถูกต้อง",
			})
		}

		if err := db.UpdatePaymentStatus(transactionId, body.Status); err != nil {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบ Transaction ID ที่ระบุ",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "อัปเดตสถานะสำเร็จ",
		})
	})
}

func main() {
	// โหลด environment variables จากไฟล์ .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// สร้าง Database
	db := NewPaymentDatabase("payments.json")

	// สร้าง Log Manager
	logManager := NewLogManager(LogFileName)

	// ทำความสะอาด log เก่าที่เกิน 5 วัน
	if err := logManager.CleanOldLogs(); err != nil {
		log.Printf("⚠️ ไม่สามารถทำความสะอาด log เก่าได้: %v", err)
	} else {
		log.Println("✅ ทำความสะอาด log เก่าเรียบร้อยแล้ว")
	}

	// สร้าง Email Service (เลือกใช้ Real หรือ Mock)
	useRealEmail := getEnv("USE_REAL_EMAIL", "false") == "true"

	var emailSender interface {
		SendPaymentNotificationEmail(PaymentNotification) error
		WriteLog(string, string)
	}

	if useRealEmail {
		realEmailService, err := NewRealEmailService()
		if err != nil {
			log.Fatal("ไม่สามารถเริ่ม Real Email Service ได้:", err)
		}
		defer realEmailService.Close()
		emailSender = realEmailService
		log.Println("✅ ใช้งาน Real Email Service")
	} else {
		mockEmailService, err := NewEmailService()
		if err != nil {
			log.Fatal("ไม่สามารถเริ่ม Mock Email Service ได้:", err)
		}
		defer mockEmailService.Close()
		emailSender = mockEmailService
		log.Println("✅ ใช้งาน Mock Email Service")
	}

	// สร้าง Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Setup routes
	setupRoutes(app, db, logManager, emailSender)

	// Start server
	port := getEnv("PORT", "3000")

	log.Printf("🚀 Email Service Server เริ่มทำงานที่พอร์ต %s", port)
	log.Printf("📧 URL: http://localhost:%s", port)
	log.Printf("📋 Health Check: http://localhost:%s/health", port)
	log.Printf("📊 API Docs: http://localhost:%s/api/payments", port)

	log.Fatal(app.Listen(":" + port))
}
