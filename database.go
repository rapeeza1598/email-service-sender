package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Database simulation using JSON file
type PaymentDatabase struct {
	FilePath string
	Data     map[string]PaymentRecord `json:"payments"`
}

type PaymentRecord struct {
	PaymentNotification
	Status      string    `json:"status"`
	SubmittedAt time.Time `json:"submitted_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewPaymentDatabase(filePath string) *PaymentDatabase {
	db := &PaymentDatabase{
		FilePath: filePath,
		Data:     make(map[string]PaymentRecord),
	}

	// โหลดข้อมูลจากไฟล์ (ถ้ามี)
	db.Load()

	return db
}

func (db *PaymentDatabase) Load() error {
	if _, err := os.Stat(db.FilePath); os.IsNotExist(err) {
		// ไฟล์ไม่มี สร้างใหม่
		return db.Save()
	}

	data, err := os.ReadFile(db.FilePath)
	if err != nil {
		return err
	}

	var dbData struct {
		Payments map[string]PaymentRecord `json:"payments"`
	}

	if err := json.Unmarshal(data, &dbData); err != nil {
		return err
	}

	db.Data = dbData.Payments
	return nil
}

func (db *PaymentDatabase) Save() error {
	dbData := struct {
		Payments map[string]PaymentRecord `json:"payments"`
		SavedAt  time.Time                `json:"saved_at"`
	}{
		Payments: db.Data,
		SavedAt:  time.Now(),
	}

	data, err := json.MarshalIndent(dbData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(db.FilePath, data, 0644)
}

func (db *PaymentDatabase) AddPayment(payment PaymentNotification) error {
	record := PaymentRecord{
		PaymentNotification: payment,
		Status:              "รอตรวจสอบ",
		SubmittedAt:         payment.Timestamp,
		UpdatedAt:           payment.Timestamp,
	}

	db.Data[payment.TransactionID] = record
	return db.Save()
}

func (db *PaymentDatabase) GetPayment(transactionID string) (*PaymentRecord, bool) {
	record, exists := db.Data[transactionID]
	return &record, exists
}

func (db *PaymentDatabase) UpdatePaymentStatus(transactionID, status string) error {
	if record, exists := db.Data[transactionID]; exists {
		record.Status = status
		record.UpdatedAt = time.Now()
		db.Data[transactionID] = record
		return db.Save()
	}
	return fmt.Errorf("transaction ID %s not found", transactionID)
}

func (db *PaymentDatabase) GetAllPayments() map[string]PaymentRecord {
	return db.Data
}

func (db *PaymentDatabase) GetPaymentsByStatus(status string) map[string]PaymentRecord {
	result := make(map[string]PaymentRecord)
	for id, record := range db.Data {
		if record.Status == status {
			result[id] = record
		}
	}
	return result
}

// Log management
type LogManager struct {
	LogFilePath string
}

func NewLogManager(filePath string) *LogManager {
	return &LogManager{
		LogFilePath: filePath,
	}
}

func (lm *LogManager) WriteLog(transactionID, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// เขียน log แบบง่าย: เฉพาะ Transaction ID และ timestamp
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, transactionID)

	// เขียน log ลงไฟล์
	file, err := os.OpenFile(lm.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error writing log: %v", err)
		return
	}
	defer file.Close()

	file.WriteString(logEntry)

	// แสดง log ใน console ด้วย
	log.Printf("[%s] %s", timestamp, transactionID)
}

func (lm *LogManager) GetLogs(transactionID string) ([]string, error) {
	data, err := os.ReadFile(lm.LogFilePath)
	if err != nil {
		return nil, err
	}

	lines := []string{}
	for _, line := range strings.Split(string(data), "\n") {
		// เนื่องจาก log ใหม่มีแค่ timestamp และ transaction ID
		if strings.Contains(line, transactionID) && strings.TrimSpace(line) != "" {
			lines = append(lines, strings.TrimSpace(line))
		}
	}

	return lines, nil
}

// ฟังก์ชันทำความสะอาด log เก่าที่เกิน 5 วัน
func (lm *LogManager) CleanOldLogs() error {
	content, err := os.ReadFile(lm.LogFilePath)
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
				logTime, err := time.Parse("2006-01-02 15:04:05", timestampStr)
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

	return os.WriteFile(lm.LogFilePath, []byte(newContent), 0666)
}
