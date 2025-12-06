package seeds

import (
	"fmt"
	"pustaka-backend/models"
	"gorm.io/gorm"
)

// CitiesSeeder seeds the cities table with initial data
// Uses FirstOrCreate for conflict resolution - won't create duplicates
func CitiesSeeder(db *gorm.DB) error {
	fmt.Println("📝 Seeding cities table...")

	// Helper function to create string pointer
	// strPtr := func(s string) *string {
	// 	return &s
	// }

	citiesData := []models.City{
		{
			Code: "ACH",
			Name: "ACEH",
		},
		{
			Code: "BDG",
			Name: "BANDUNG",
		},
		{
			Code: "BB",
			Name: "BANGKA BELITUNG",
		},
		{
			Code: "BNT",
			Name: "BANTEN",
		},
		{
			Code: "BKS",
			Name: "BEKASI",
		},
		{
			Code: "BKL",
			Name: "BENGKULU",
		},
		{
			Code: "BGR",
			Name: "BOGOR",
		},
		{
			Code: "CJR",
			Name: "CIANJUR",
		},
		{
			Code: "CLG",
			Name: "CILEDUG",
		},
		{
			Code: "DKI",
			Name: "DAERAH KHUSUS IBUKOTA JAKARTA",
		},
		{
			Code: "DPK",
			Name: "DEPOK",
		},
		{
			Code: "JKB",
			Name: "JAKARTA BARAT",
		},
		{
			Code: "JKP",
			Name: "JAKARTA PUSAT",
		},
		{
			Code: "JKS",
			Name: "JAKARTA SELATAN",
		},
		{
			Code: "JKT",
			Name: "JAKARTA TIMUR",
		},
		{
			Code: "JKU",
			Name: "JAKARTA UTARA",
		},
		{
			Code: "JMB",
			Name: "JAMBI",
		},
		{
			Code: "JWB",
			Name: "JAWA BARAT",
		},
		{
			Code: "JWT",
			Name: "JAWA TENGAH",
		},
		{
			Code: "JTM",
			Name: "JAWA TIMUR",
		},
		{
			Code: "KB",
			Name: "KALIMANTAN BARAT",
		},
		{
			Code: "KSL",
			Name: "KALIMANTAN SELATAN",
		},
		{
			Code: "KAT",
			Name: "KALIMANTAN TENGAH",
		},
		{
			Code: "KTM",
			Name: "KALIMANTAN TIMUR",
		},
		{
			Code: "LMP",
			Name: "LAMPUNG",
		},
		{
			Code: "MDN",
			Name: "MEDAN",
		},
		{
			Code: "NTT",
			Name: "NUSA TENGGARA TIMUR",
		},
		{
			Code: "PDG",
			Name: "PADANG",
		},
		{
			Code: "PLG",
			Name: "PALEMBANG",
		},
		{
			Code: "PAP",
			Name: "PAPUA",
		},
		{
			Code: "PKU",
			Name: "PEKANBARU",
		},
		{
			Code: "PON",
			Name: "PONTIANAK",
		},
		{
			Code: "PRW",
			Name: "PURWAKARTA",
		},
		{
			Code: "RIU",
			Name: "RIAU",
		},
		{
			Code: "SLO",
			Name: "SOLO",
		},
		{
			Code: "SBG",
			Name: "SUBANG",
		},
		{
			Code: "SLT",
			Name: "SULAWESI TENGAH",
		},
		{
			Code: "ST",
			Name: "SULAWESI TENGGARA",
		},
		{
			Code: "SU",
			Name: "SULAWESI UTARA",
		},
		{
			Code: "SBR",
			Name: "SUMATERA BARAT",
		},
		{
			Code: "SSL",
			Name: "SUMATERA SELATAN",
		},
		{
			Code: "SUT",
			Name: "SUMATERA UTARA",
		},
		{
			Code: "TNG",
			Name: "TANGERANG",
		},
	}

	created := 0
	skipped := 0

	// delete all existing records before seeding
	db.Exec("DELETE FROM cities")

	// Insert all records using FirstOrCreate
	for _, city := range citiesData {
		var result models.City
		err := db.Where("code = ?", city.Code).FirstOrCreate(&result, city).Error
		if err != nil {
			return fmt.Errorf("failed to seed cities with code %s: %w", city.Code, err)
		}

		// Check if it was created or already existed
		if result.CreatedAt.Equal(result.UpdatedAt) && result.Code == city.Code {
			created++
		} else {
			skipped++
		}
	}

	fmt.Printf("✓ Cities seeding completed: %d created, %d skipped (already exist)\n", created, skipped)
	return nil
}
