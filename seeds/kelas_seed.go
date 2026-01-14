package seeds

import (
	"fmt"
	"pustaka-backend/models"

	"gorm.io/gorm"
)

// KelasSeeder seeds the kelas table with initial data
// Uses FirstOrCreate for conflict resolution - won't create duplicates
func KelasSeeder(db *gorm.DB) error {
	fmt.Println("Seeding kelas table...")

	kelasData := []models.Kelas{
		{
			Code: "1",
			Name: "Kelas 1",
		},
		{
			Code: "2",
			Name: "Kelas 2",
		},
		{
			Code: "3",
			Name: "Kelas 3",
		},
		{
			Code: "4",
			Name: "Kelas 4",
		},
		{
			Code: "5",
			Name: "Kelas 5",
		},
		{
			Code: "6",
			Name: "Kelas 6",
		},
		{
			Code: "7",
			Name: "Kelas 7",
		},
		{
			Code: "8",
			Name: "Kelas 8",
		},
		{
			Code: "9",
			Name: "Kelas 9",
		},
		{
			Code: "10",
			Name: "Kelas 10",
		},
		{
			Code: "11",
			Name: "Kelas 11",
		},
		{
			Code: "12",
			Name: "Kelas 12",
		},
		{
			Code: "A",
			Name: "TK A (Usia 4-5)",
		},
		{
			Code: "B",
			Name: "TK B (Usia 5-6)",
		},
		{
			Code: "ALL",
			Name: "Semua Tingkat",
		},
	}

	created := 0
	skipped := 0

	// Delete all existing records before seeding
	db.Exec("DELETE FROM kelas")

	// Insert all records using FirstOrCreate
	for _, kelas := range kelasData {
		var result models.Kelas
		err := db.Where("code = ?", kelas.Code).FirstOrCreate(&result, kelas).Error
		if err != nil {
			return fmt.Errorf("failed to seed kelas with code %s: %w", kelas.Code, err)
		}

		// Check if it was created or already existed
		if result.CreatedAt.Equal(result.UpdatedAt) && result.Code == kelas.Code {
			created++
		} else {
			skipped++
		}
	}

	fmt.Printf("Kelas seeding completed: %d created, %d skipped (already exist)\n", created, skipped)
	return nil
}
