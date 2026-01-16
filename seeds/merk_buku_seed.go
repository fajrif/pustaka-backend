package seeds

import (
	"fmt"
	"pustaka-backend/models"

	"gorm.io/gorm"
)

// MerkBukuSeeder seeds the merk_buku table with initial data
// Uses FirstOrCreate for conflict resolution - won't create duplicates
func MerkBukuSeeder(db *gorm.DB) error {
	fmt.Println("Seeding merk_buku table...")

	merkBukuData := []models.MerkBuku{
		{
			Code: "SKR",
			Name: "Sekar",
		},
		{
			Code: "WJR",
			Name: "Wajar",
		},
		{
			Code: "TUN",
			Name: "Tuntas",
		},
		{
			Code: "MEN",
			Name: "Mentari",
		},
		{
			Code: "FTH",
			Name: "Fattah",
		},
		{
			Code: "FJR",
			Name: "Fajar",
		},
		{
			Code: "FTR",
			Name: "Fitrah",
		},
		{
			Code: "HKM",
			Name: "HIkmah",
		},
		{
			Code: "PRI",
			Name: "Prima",
		},
		{
			Code: "KOD",
			Name: "Koding",
		},
		{
			Code: "INT",
			Name: "Intens",
		},
		{
			Code: "MAX",
			Name: "MAXXI",
		},
		{
			Code: "KAR",
			Name: "Kartika",
		},
		{
			Code: "TKA",
			Name: "Tes Kemampuan Akademi",
		},
		{
			Code: "PSAJ",
			Name: "Penilaian Sumatif Akhir Jenjang",
		},
		{
			Code: "PSAJM",
			Name: "PSAJ Madrasah",
		},
		{
			Code: "UN",
			Name: "Ujian Negara",
		},
		{
			Code: "US",
			Name: "Ujian Sekolah",
		},
	}

	created := 0
	skipped := 0

	// Delete all existing records before seeding
	db.Exec("DELETE FROM merk_buku")

	// Insert all records using FirstOrCreate
	for _, merkBuku := range merkBukuData {
		var result models.MerkBuku
		err := db.Where("code = ?", merkBuku.Code).FirstOrCreate(&result, merkBuku).Error
		if err != nil {
			return fmt.Errorf("failed to seed merk_buku with code %s: %w", merkBuku.Code, err)
		}

		// Check if it was created or already existed
		if result.CreatedAt.Equal(result.UpdatedAt) && result.Code == merkBuku.Code {
			created++
		} else {
			skipped++
		}
	}

	fmt.Printf("Merk Buku seeding completed: %d created, %d skipped (already exist)\n", created, skipped)
	return nil
}
