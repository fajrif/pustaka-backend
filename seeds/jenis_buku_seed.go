package seeds

import (
	"fmt"
	"pustaka-backend/models"
	"gorm.io/gorm"
)

// JenisBukuSeeder seeds the jenis_buku table with initial data
// Uses FirstOrCreate for conflict resolution - won't create duplicates
func JenisBukuSeeder(db *gorm.DB) error {
	fmt.Println("📝 Seeding jenis_buku table...")

	// Helper function to create string pointer
	strPtr := func(s string) *string {
		return &s
	}

	jenisBukuData := []models.JenisBuku{
		{
			Code:        "FIK",
			Name:        "Fiksi",
			Description: strPtr("Buku fiksi dan novel"),
		},
		{
			Code:        "NON-FIK",
			Name:        "Non-Fiksi",
			Description: strPtr("Buku non-fiksi dan referensi"),
		},
		{
			Code:        "PEND",
			Name:        "Pendidikan",
			Description: strPtr("Buku pendidikan dan akademik"),
		},
		{
			Code:        "ANAK",
			Name:        "Anak-Anak",
			Description: strPtr("Buku untuk anak-anak"),
		},
		{
			Code:        "REF",
			Name:        "Referensi",
			Description: strPtr("Buku referensi seperti kamus dan ensiklopedia"),
		},
	}

	created := 0
	skipped := 0

	// Insert all records using FirstOrCreate
	for _, jenisBuku := range jenisBukuData {
		var result models.JenisBuku
		err := db.Where("code = ?", jenisBuku.Code).FirstOrCreate(&result, jenisBuku).Error
		if err != nil {
			return fmt.Errorf("failed to seed jenis_buku with code %s: %w", jenisBuku.Code, err)
		}

		// Check if it was created or already existed
		if result.CreatedAt.Equal(result.UpdatedAt) && result.Code == jenisBuku.Code {
			created++
		} else {
			skipped++
		}
	}

	fmt.Printf("✓ Jenis Buku seeding completed: %d created, %d skipped (already exist)\n", created, skipped)
	return nil
}
