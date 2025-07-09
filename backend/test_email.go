package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	to := "ubekboo.1997@gmail.com" // 👉 сюда вставь email получателя

	msg := []byte("To: " + to + "\r\n" +
		"Subject: Тестовое письмо\r\n" +
		"\r\n" +
		"Привет! Это тестовое письмо от Go-программы через Gmail SMTP.\r\n")

	addr := os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT")
	auth := smtp.PlainAuth("", from, pass, os.Getenv("SMTP_HOST"))

	err = smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		log.Fatal("Не удалось отправить письмо: ", err)
	}

	fmt.Println("✅ Письмо успешно отправлено!")
}
