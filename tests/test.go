package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func testEmailService() {
	// ทดสอบ API
	fmt.Println("🧪 เริ่มทดสอบ Email Service...")

	// ข้อมูลทดสอบ
	testPayment := map[string]interface{}{
		"account":          "นาย ทดสอบ ระบบ",
		"transactionId":    "TEST123456789",
		"transferDate":     time.Now().Format("2006-01-02"),
		"transferTime":     time.Now().Format("15:04"),
		"recipientAccount": "123-456-7890",
		"amount":           1500.00,
		"additional":       "ทดสอบระบบส่งเมล",
	}

	jsonData, _ := json.Marshal(testPayment)

	// ส่งคำขอไปยัง API
	resp, err := http.Post("http://localhost:3000/api/payment-notification",
		"application/json",
		bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Printf("❌ ไม่สามารถเชื่อมต่อกับเซิร์ฟเวอร์ได้: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("✅ Response Status: %s\n", resp.Status)

	// ทดสอบการดู logs
	logResp, err := http.Get("http://localhost:3000/api/logs/TEST123456789")
	if err != nil {
		fmt.Printf("❌ ไม่สามารถดู logs ได้: %v\n", err)
		return
	}
	defer logResp.Body.Close()

	fmt.Printf("✅ Log Response Status: %s\n", logResp.Status)
	fmt.Println("🎉 การทดสอบเสร็จสิ้น!")
}

func main() {
	// รอให้เซิร์ฟเวอร์เริ่มทำงาน
	time.Sleep(2 * time.Second)
	testEmailService()
}
