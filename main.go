package main

import (
	"context"
	"fmt"
	//"time"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

// RedisClient represents a Redis client
var RedisClient *redis.Client
var ctx = context.Background()

func initRedis() {
	// Connect to Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Ping Redis to check connectivity
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

func sendEmail() {
	// Simulate fetching data from Redis
	data, err := RedisClient.Get(ctx, "email_data").Result()
	if err != nil {
		fmt.Println("Error fetching data from Redis:", err)
		return
	}

	// Initialize email message
	mail := gomail.NewMessage()
	mail.SetHeader("From", "katherin.shibuki@gmail.com")
	mail.SetHeader("To", "katherin.shibuki@gmail.com")
	mail.SetHeader("Subject", "Scheduled Email")
	mail.SetBody("text/plain", data)

	// Initialize SMTP dialer
	dialer := gomail.NewDialer("smtp.example.com", 587, "username", "password")

	// Send email
	if err := dialer.DialAndSend(mail); err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}

func main() {
	// Initialize Redis
	initRedis()

	// Simulate caching data into Redis
	err := RedisClient.Set(ctx, "email_data", "Hello, this is a scheduled email!", 0).Err()
	if err != nil {
		fmt.Println("Error caching data into Redis:", err)
		return
	}

	// Initialize CRON scheduler
	c := cron.New()

	// Schedule email sending every minute
	c.AddFunc("@every 1m", sendEmail)

	// Start CRON scheduler
	c.Start()

	// Keep the program running
	select {}
}
