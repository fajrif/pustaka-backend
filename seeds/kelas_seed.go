package seeds

import (
	"fmt"
	"pustaka-backend/models"
	"gorm.io/gorm"
)

// KelasSeeder seeds the kelas table with initial data
// Uses FirstOrCreate for conflict resolution - won't create duplicates
func KelasSeeder(db *gorm.DB) error {
	fmt.Println("📝 Seeding kelas table...")

	// Helper function to create string pointer
	// strPtr := func(s string) *string {
	// 	return &s
	// }

	kelasData := []models.Kelas{
		{
				Code: "ALL",
				Name: "SEMUA TINGKAT",
		},
		{
				Code: "A",
				Name: "USIA 4-5 TAHUN",
		},
		{
				Code: "B",
				Name: "USIA 5-6 TAHUN",
		},
		{
				Code: "K1",
				Name: "KELAS I",
		},
		{
				Code: "K2",
				Name: "KELAS II",
		},
		{
				Code: "K3",
				Name: "KELAS III",
		},
		{
				Code: "K4",
				Name: "KELAS IV",
		},
		{
				Code: "K5",
				Name: "KELAS V",
		},
		{
				Code: "K6",
				Name: "KELAS VI",
		},
		{
				Code: "K7",
				Name: "KELAS VII",
		},
		{
				Code: "K8",
				Name: "KELAS VIII",
		},
		{
				Code: "K9",
				Name: "KELAS IX",
		},
		{
				Code: "K10",
				Name: "KELAS X",
		},
		{
				Code: "K11",
				Name: "KELAS XI",
		},
		{
				Code: "K12",
				Name: "KELAS XII",
		},
	}

	created := 0
	skipped := 0

	// delete all existing records before seeding
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

	fmt.Printf("✓ Kelas seeding completed: %d created, %d skipped (already exist)\n", created, skipped)
	return nil
}
