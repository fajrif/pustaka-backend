package seeds

import (
	"fmt"
	"pustaka-backend/models"
	"gorm.io/gorm"
)

// JenjangStudiSeeder seeds the jenjang_studi table with initial data
// Uses FirstOrCreate for conflict resolution - won't create duplicates
func JenjangStudiSeeder(db *gorm.DB) error {
	fmt.Println("📝 Seeding jenjang_studi table...")

	// Helper function to create string pointer
	// strPtr := func(s string) *string {
	// 	return &s
	// }

	jenjangStudiData := []models.JenjangStudi{
		{
			Code: "SMACL",
			Name: "SMA COVER LAMA",
		},
		{
			Code: "SMKCL",
			Name: "SMK COVER LAMA",
		},
		{
			Code: "MJL",
			Name: "SEKOLAH MENENGAH PERTAMA MAJELIS",
		},
		{
			Code: "MJL/P",
			Name: "SEKOLAH MENENGAH PERTAMA PRIBADI",
		},
		{
			Code: "SMPCL",
			Name: "SEKOLAH MENENGAH PERTAMA CL",
		},
		{
			Code: "SMABK",
			Name: "SMA",
		},
		{
			Code: "SMP/U",
			Name: "SEKOLAH MENENGAH PERTAMA (UMUM)",
		},
		{
			Code: "SMA/U",
			Name: "SEKOLAH MENENGAH ATAS (UMUM)",
		},
		{
			Code: "STM",
			Name: "SEKOLAH MENENGAH TEHNIK",
		},
		{
			Code: "SD/S",
			Name: "SEKOLAH DASAR SWASTA",
		},
		{
			Code: "SD/G",
			Name: "SEKOLAH DASAR (GENERAL)",
		},
		{
			Code: "SMK",
			Name: "SEKOLAH MENENGAH KEJURUAN",
		},
		{
			Code: "SMAX",
			Name: "SEKOLAH MENENGAH ATAS X",
		},
		{
			Code: "SMK/U",
			Name: "SEKOLAH MENENGAH KEJURUAN (UMUM)",
		},
		{
			Code: "SMPM",
			Name: "SMPM",
		},
		{
			Code: "SMPKB",
			Name: "SMPKB",
		},
		{
			Code: "SD/M",
			Name: "SEKOLAH DASAR M",
		},
		{
			Code: "SD",
			Name: "SEKOLAH DASAR",
		},
		{
			Code: "UMUM",
			Name: "SD / SMP",
		},
		{
			Code: "SMPX",
			Name: "SMP KARTIKA CL",
		},
		{
			Code: "SMP",
			Name: "SEKOLAH MENENGAH PERTAMA",
		},
		{
			Code: "MTS",
			Name: "MADRASAH TSANAWIYAH",
		},
		{
			Code: "IT",
			Name: "SMK INFORMATIKA",
		},
		{
			Code: "BIS",
			Name: "SMK BISMEN",
		},
		{
			Code: "SMA",
			Name: "SEKOLAH MENENGAH ATAS",
		},
		{
			Code: "SMP/M",
			Name: "SEKOLAH MENENGAH PERTAMA M",
		},
		{
			Code: "SMP/P",
			Name: "SMP PRIBADI",
		},
		{
			Code: "SMA/P",
			Name: "SMA PRIBADI",
		},
		{
			Code: "SD/P",
			Name: "SD PRIBADI",
		},
		{
			Code: "MI",
			Name: "MADRASAH IBTIDAIYAH",
		},
		{
			Code: "PRO",
			Name: "SMK PRODUKTIF",
		},
		{
			Code: "TKN",
			Name: "SMK (TEKNOLOGI  & REKAYASA )",
		},
		{
			Code: "SMP1",
			Name: "SEKOLAH MENENGAH PERTAMA 1",
		},
		{
			Code: "MI/G",
			Name: "MI",
		},
		{
			Code: "SD/CL",
			Name: "SD COVER LAMA",
		},
		{
			Code: "KEA",
			Name: "SMK KEAHLIAN",
		},
		{
			Code: "UM",
			Name: "UMUM (SD/SMP/SMA)",
		},
		{
			Code: "PROP",
			Name: "PRODUKTIF PG",
		},
		{
			Code: "TK A",
			Name: "TK USIA 4-5 TAHUN",
		},
		{
			Code: "TK B",
			Name: "TK USIA 5-6 TAHUN",
		},
	}

	created := 0
	skipped := 0

	// delete all existing records before seeding
	db.Exec("DELETE FROM jenjang_studi")

	// Insert all records using FirstOrCreate
	for _, jenjang_studi := range jenjangStudiData {
		var result models.JenjangStudi
		err := db.Where("code = ?", jenjang_studi.Code).FirstOrCreate(&result, jenjang_studi).Error
		if err != nil {
			return fmt.Errorf("failed to seed jenjang_studi with code %s: %w", jenjang_studi.Code, err)
		}

		// Check if it was created or already existed
		if result.CreatedAt.Equal(result.UpdatedAt) && result.Code == jenjang_studi.Code {
			created++
		} else {
			skipped++
		}
	}

	fmt.Printf("✓ Jenjang Studi seeding completed: %d created, %d skipped (already exist)\n", created, skipped)
	return nil
}
