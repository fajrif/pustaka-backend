package seeds

import (
	"fmt"
	"pustaka-backend/models"

	"gorm.io/gorm"
)

// JenjangStudiSeeder seeds the jenjang_studi table with initial data
// Uses FirstOrCreate for conflict resolution - won't create duplicates
func JenjangStudiSeeder(db *gorm.DB) error {
	fmt.Println("Seeding jenjang_studi table...")

	jenjangStudiData := []models.JenjangStudi{
		{
			Code: "TK",
			Name: "Taman Kanak",
		},
		{
			Code: "SD",
			Name: "Sekolah Dasar",
		},
		{
			Code: "SMP",
			Name: "Sekolah Menengah Pertama",
		},
		{
			Code: "SMA",
			Name: "Sekolah Menengah Atas",
		},
		{
			Code: "SMK",
			Name: "Sekolah Menengah Kejuruan",
		},
		{
			Code: "MTS",
			Name: "Madrasah Tsanawiyah",
		},
		{
			Code: "MI",
			Name: "Madrasah Ibtidaiyah",
		},
		{
			Code: "MA",
			Name: "Madrasah Aliyah",
		},
	}

	created := 0
	skipped := 0

	// Delete all existing records before seeding
	db.Exec("DELETE FROM jenjang_studi")

	// Insert all records using FirstOrCreate
	for _, jenjangStudi := range jenjangStudiData {
		var result models.JenjangStudi
		err := db.Where("code = ?", jenjangStudi.Code).FirstOrCreate(&result, jenjangStudi).Error
		if err != nil {
			return fmt.Errorf("failed to seed jenjang_studi with code %s: %w", jenjangStudi.Code, err)
		}

		// Check if it was created or already existed
		if result.CreatedAt.Equal(result.UpdatedAt) && result.Code == jenjangStudi.Code {
			created++
		} else {
			skipped++
		}
	}

	fmt.Printf("Jenjang Studi seeding completed: %d created, %d skipped (already exist)\n", created, skipped)
	return nil
}
