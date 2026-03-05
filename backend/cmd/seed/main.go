package main

import (
	"fmt"
	"log"
	"partitionlab/internal/app/ds"
	"partitionlab/internal/app/dsn"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load("../../.env")

	postgresString := dsn.FromEnv()
	db, err := gorm.Open(postgres.Open(postgresString), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Создаем пользователей
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	users := []ds.User{
		{Login: "user1", Password: string(hashedPassword), IsModerator: false},
		{Login: "moderator", Password: string(hashedPassword), IsModerator: true},
	}

	for _, user := range users {
		var existing ds.User
		if err := db.Where("login = ?", user.Login).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&user).Error; err != nil {
				log.Printf("failed to create user %s: %v", user.Login, err)
			} else {
				log.Printf("Created user: %s", user.Login)
			}
		} else {
			log.Printf("User %s already exists", user.Login)
		}
	}

	// Создаем типы перегородок
	partitions := []ds.Partition{
		{
			Title:          "Гипсокартон ГКЛ 12.5мм + минвата 50мм",
			Category:       "Легкие",
			Description:    "Каркасная перегородка из оцинкованного профиля, обшитая гипсокартоном с двух сторон, заполненная минеральной ватой. Оптимальна для офисов и жилых помещений.",
			NoiseReduction: "35-40 дБ",
			Thickness:      "10-12 см",
			Material:       "ГКЛ + минеральная вата",
			PricePerSqm:    "700-1000 руб/м²",
			ImageURL:       "/partitions/гипсокартонная.png",
			IsActive:       true,
		},
		{
			Title:          "Кирпичная кладка 120мм",
			Category:       "Тяжелые",
			Description:    "Кирпичная перегородка в полкирпича. Высокая звукоизоляция за счет массивности конструкции. Требует прочного основания.",
			NoiseReduction: "45-50 дБ",
			Thickness:      "12 см",
			Material:       "Керамический кирпич",
			PricePerSqm:    "1500-2000 руб/м²",
			ImageURL:       "/partitions/кирпичная.png",
			IsActive:       true,
		},
		{
			Title:          "Газобетонные блоки D500 100мм",
			Category:       "Средние",
			Description:    "Перегородка из газобетонных блоков. Легче кирпича, хорошая звукоизоляция. Простой монтаж на клеевой раствор.",
			NoiseReduction: "40-44 дБ",
			Thickness:      "10 см",
			Material:       "Газобетон D500",
			PricePerSqm:    "900-1300 руб/м²",
			ImageURL:       "/partitions/газобетонная.jpg",
			IsActive:       true,
		},
		{
			Title:          "Сэндвич-панели акустические",
			Category:       "Специализированные",
			Description:    "Многослойные панели с акустическим наполнителем. Отличная звукоизоляция при малой толщине. Подходят для студий звукозаписи.",
			NoiseReduction: "50-55 дБ",
			Thickness:      "8-10 см",
			Material:       "Сталь/алюминий + акустический наполнитель",
			PricePerSqm:    "2000-3500 руб/м²",
			ImageURL:       "/partitions/сэндвич-панель.jpg",
			IsActive:       true,
		},
		{
			Title:          "Двойная каркасная перегородка с виброподвесами",
			Category:       "Профессиональные",
			Description:    "Двойная каркасная система с раздельными профилями, виброизоляцией и многослойной обшивкой. Максимальная звукоизоляция для домашних кинотеатров.",
			NoiseReduction: "55-65 дБ",
			Thickness:      "18-20 см",
			Material:       "ГКЛ + акустический гипсокартон + виброподвесы + минвата",
			PricePerSqm:    "3000-5000 руб/м²",
			ImageURL:       "/partitions/двойная.png",
			IsActive:       true,
		},
	}

	// Обновляем существующие перегородки по Title или создаём новые
	for _, partition := range partitions {
		var existing ds.Partition
		err := db.Where("title = ?", partition.Title).First(&existing).Error
		if err == nil {
			updates := map[string]any{
				"category":        partition.Category,
				"description":     partition.Description,
				"noise_reduction": partition.NoiseReduction,
				"thickness":       partition.Thickness,
				"material":        partition.Material,
				"price_per_sqm":   partition.PricePerSqm,
				"image_url":       partition.ImageURL,
				"is_active":       partition.IsActive,
			}
			if err := db.Model(&existing).Updates(updates).Error; err != nil {
				log.Printf("failed to update partition %s: %v", partition.Title, err)
			} else {
				log.Printf("Updated partition: %s", partition.Title)
			}
			continue
		}

		if err := db.Create(&partition).Error; err != nil {
			log.Printf("failed to create partition %s: %v", partition.Title, err)
		} else {
			log.Printf("Created partition: %s", partition.Title)
		}
	}

	fmt.Println("\n✅ Seeding completed successfully!")
	fmt.Println("Created users: user1, moderator (password: 123456)")
	fmt.Println("Created 5 partition types")
}
