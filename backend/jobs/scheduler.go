package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"medislink-backend/config"
	"medislink-backend/controllers"
	"medislink-backend/models"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func StartCronJob() {
	c := cron.New()

	_, err := c.AddFunc("0 9 * * *", func() {
		fmt.Println("[CRON] Menjalankan pengecekan pengingat pengembalian...")
		CheckDueLoans()
	})

	if err != nil {
		log.Fatal("Gagal menjalankan cron job:", err)
	}

	c.Start()
	fmt.Println("[SYSTEM] Cron Job Berjalan: Cek pengembalian setiap jam 09:00")
}

func CheckDueLoans() {
	coll := config.DB.Collection("loans")

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	filter := bson.M{
		"status":     "active",
		"return_due": tomorrow,
	}

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		log.Println("[CRON ERROR] Gagal mengambil data loan:", err)
		return
	}
	defer cursor.Close(context.Background())

	var loans []models.Loan
	if err = cursor.All(context.Background(), &loans); err != nil {
		log.Println("[CRON ERROR] Gagal decode loan:", err)
		return
	}

	if len(loans) > 0 {
		fmt.Printf("[CRON] Ditemukan %d peminjaman yang harus kembali besok (%s)\n", len(loans), tomorrow)

		for _, loan := range loans {
			title := "Pengingat Pengembalian Alat"
			message := fmt.Sprintf("Halo! Jangan lupa alat medis pinjaman Anda harus dikembalikan besok (%s). Mohon persiapkan pengembalian.", tomorrow)

			controllers.CreateNotificationService(loan.UserID, title, message, "warning")
		}
	} else {
		fmt.Println("[CRON] Tidak ada peminjaman yang jatuh tempo besok.")
	}
}
