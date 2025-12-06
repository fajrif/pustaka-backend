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
	// strPtr := func(s string) *string {
	// 	return &s
	// }

	jenisBukuData := []models.JenisBuku{
		{
				Code: "AKM",
				Name: "ASESMEN KOMP MINIMUN",
		},
		{
				Code: "BAS",
				Name: "BAHAN AJAR SISWA",
		},
		{
				Code: "PK1",
				Name: "BSE",
		},
		{
				Code: "BG",
				Name: "BUKU GURU",
		},
		{
				Code: "PKT",
				Name: "BUKU MATERI",
		},
		{
				Code: "KBR",
				Name: "K. BULAN RAMADHAN",
		},
		{
				Code: "KA",
				Name: "KOMPUTER AKUNTANSI",
		},
		{
				Code: "KB",
				Name: "KOMUNIKASI BISNIS",
		},
		{
				Code: "KM",
				Name: "KURIKULUM MERDEKA",
		},
		{
				Code: "LKS",
				Name: "LEMBAR KERJA SISWA",
		},
		{
				Code: "MTS",
				Name: "MADRASAH TSANAWIYAH",
		},
		{
				Code: "MJL",
				Name: "MAJALAH",
		},
		{
				Code: "MDL",
				Name: "MODUL",
		},
		{
				Code: "NK",
				Name: "NON KURIKULUM",
		},
		{
				Code: "PAR",
				Name: "PARIWISATA",
		},
		{
				Code: "PGP",
				Name: "PEG GURU PSAJ",
		},
		{
				Code: "PG",
				Name: "PEGANGAN GURU",
		},
		{
				Code: "PGU",
				Name: "PEGANGAN GURU US",
		},
		{
				Code: "PSA",
				Name: "PEN SUMATIF AKHIR",
		},
		{
				Code: "PPP",
				Name: "PENG & PEND PRO",
		},
		{
				Code: "PGB",
				Name: "PG K13",
		},
		{
				Code: "PGM",
				Name: "PG KUR MERDEKA",
		},
		{
				Code: "PRP",
				Name: "PG PRODUKTIF",
		},
		{
				Code: "PGT",
				Name: "PG TEMATIK",
		},
		{
				Code: "PGA",
				Name: "PG TKA",
		},
		{
				Code: "PRO",
				Name: "PRODUKTIF",
		},
		{
				Code: "SEK",
				Name: "SENI DAN EKONOMI KRE",
		},
		{
				Code: "TMT",
				Name: "TEMATIK",
		},
		{
				Code: "TKA",
				Name: "TES KEMAMPUAN AKADEM",
		},
		{
				Code: "UAN",
				Name: "UAS",
		},
		{
				Code: "UM",
				Name: "UJIAN MADRASAH",
		},
	}

	created := 0
	skipped := 0

	// delete all existing records before seeding
	db.Exec("DELETE FROM jenis_buku")

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
