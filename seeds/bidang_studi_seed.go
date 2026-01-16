package seeds

import (
	"fmt"
	"pustaka-backend/models"

	"gorm.io/gorm"
)

// BidangStudiSeeder seeds the bidang_studi table with initial data
// Uses FirstOrCreate for conflict resolution - won't create duplicates
func BidangStudiSeeder(db *gorm.DB) error {
	fmt.Println("📝 Seeding bidang_studi table...")

	// Helper function to create string pointer
	// strPtr := func(s string) *string {
	// 	return &s
	// }

	bidangStudiData := []models.BidangStudi{
		{
			Code: "PPKN",
			Name: "PPKN",
		},
		{
			Code: "PNC",
			Name: "Pendidikan Pancasila",
		},
		{
			Code: "IND",
			Name: "Bahasa Indonesia",
		},
		{
			Code: "INDL",
			Name: "Bahasa Indonesia TK Lanjut",
		},
		{
			Code: "INDL1",
			Name: "Bahasa Indonesia TK Lanjut (1 Tahun)",
		},
		{
			Code: "ING",
			Name: "Bahasa Inggris",
		},
		{
			Code: "INGL",
			Name: "Bahasa Inggris TK Lanjut",
		},
		{
			Code: "INGL1",
			Name: "Bahasa Inggris TK Lanjut (1 Tahun)",
		},
		{
			Code: "ARB",
			Name: "Bahasa Arab",
		},
		{
			Code: "MTK",
			Name: "Matematika",
		},
		{
			Code: "MTKL",
			Name: "Matematika TK Lanjut",
		},
		{
			Code: "PJOK",
			Name: "Pendidikan Jasmani, Olahraga, dan Kesehatan",
		},
		{
			Code: "IPL",
			Name: "IPAS GASAL",
		},
		{
			Code: "IPN",
			Name: "IPAS GENAP",
		},
		{
			Code: "PAI",
			Name: "Pendidikan Agama Islam",
		},
		{
			Code: "SM",
			Name: "Seni Musik",
		},
		{
			Code: "SR",
			Name: "Seni Rupa",
		},
		{
			Code: "ST",
			Name: "Seni Tari",
		},
		{
			Code: "STR",
			Name: "Seni Teater",
		},
		{
			Code: "STR1",
			Name: "Seni Teater (1 Tahun)",
		},
		{
			Code: "STP",
			Name: "Seni Terpadu",
		},
		{
			Code: "STP1",
			Name: "Seni Terpadu (1 Tahun)",
		},
		{
			Code: "IPA",
			Name: "IPA",
		},
		{
			Code: "IPT",
			Name: "IPA Terpadu",
		},
		{
			Code: "IPS",
			Name: "IPS",
		},
		{
			Code: "IPST",
			Name: "IPS Terpadu",
		},
		{
			Code: "PJ",
			Name: "Pendidikan Jasmani",
		},
		{
			Code: "PRK",
			Name: "Prakarya Kerajinan",
		},
		{
			Code: "PRKP",
			Name: "Prakarya Pengolahan",
		},
		{
			Code: "PRKB",
			Name: "Prakarya Budidaya",
		},
		{
			Code: "PRKR",
			Name: "Prakarya Rekayasa",
		},
		{
			Code: "PRKT",
			Name: "Prakarya Terpadu",
		},
		{
			Code: "PRKT1",
			Name: "Prakarya Terpadu (1 Tahun)",
		},
		{
			Code: "PKWUB",
			Name: "Prakarya & KWU - Budidaya",
		},
		{
			Code: "PKWUB1",
			Name: "Prakarya & KWU - Budidaya (1 Tahun)",
		},
		{
			Code: "PKWUR",
			Name: "Prakarya & KWU - Rekayasa",
		},
		{
			Code: "PKWUR1",
			Name: "Prakarya & KWU - Rekayasa (1 Tahun)",
		},
		{
			Code: "INF",
			Name: "Informatika",
		},
		{
			Code: "BK",
			Name: "Bimbingan Konseling",
		},
		{
			Code: "SJR",
			Name: "Sejarah",
		},
		{
			Code: "IPA1",
			Name: "IPAS (1 Tahun)",
		},
		{
			Code: "SJR1",
			Name: "SEJARAH (1 Tahun)",
		},
		{
			Code: "IPB",
			Name: "IPAS Bisnis",
		},
		{
			Code: "IPK",
			Name: "IPAS Kesehatan",
		},
		{
			Code: "IPG",
			Name: "IPAS Teknologi",
		},
		{
			Code: "AKU",
			Name: "Akutansi dan Keuangan Lembaga",
		},
		{
			Code: "DSG",
			Name: "Desain Komunikasi Visual",
		},
		{
			Code: "MPLB",
			Name: "Manajemen Perkantoran dan Layanan Bisnis",
		},
		{
			Code: "PMS",
			Name: "Pemasaran",
		},
		{
			Code: "PPLG",
			Name: "Pengembangan Perangkat Lunak dan Gim",
		},
		{
			Code: "PGDR",
			Name: "Pengembangan Diri",
		},
		{
			Code: "HTL",
			Name: "Perhotelan",
		},
		{
			Code: "TJKT",
			Name: "Teknik Jaringan Komputer dan Telekomunikasi",
		},
		{
			Code: "TKM",
			Name: "Teknik Mesin",
		},
		{
			Code: "TKO",
			Name: "Teknik Otomotif",
		},
		{
			Code: "SDG",
			Name: "SD GABUNGAN",
		},
		{
			Code: "SDT",
			Name: "SD TERPADU",
		},
		{
			Code: "MIT",
			Name: "MI TERPADU",
		},
		{
			Code: "PAIBP",
			Name: "Agama Islam & BP",
		},
		{
			Code: "USSDT",
			Name: "USBN / US SD Terpadu",
		},
		{
			Code: "TPKWU",
			Name: "Produk / Projek Kreatif & KWU - SMK",
		},
		{
			Code: "QRH",
			Name: "Quran hadist",
		},
		{
			Code: "FIQ",
			Name: "Fiqih",
		},
		{
			Code: "AQA",
			Name: "Aqidah Akhlak",
		},
		{
			Code: "SKI",
			Name: "Sejarah Kebudayaan Islam",
		},
		{
			Code: "FIS",
			Name: "Fisika",
		},
		{
			Code: "BIO",
			Name: "Biologi",
		},
		{
			Code: "KIM",
			Name: "Kimia",
		},
		{
			Code: "EKO",
			Name: "Ekonomi",
		},
		{
			Code: "GEO",
			Name: "Geografi",
		},
		{
			Code: "SOS",
			Name: "Sosiologi",
		},
		{
			Code: "ANT",
			Name: "Antropologi",
		},
		{
			Code: "ANT1",
			Name: "Antropologi (1 Tahun)",
		},
		{
			Code: "KOD",
			Name: "Koding dan Kecerdasan Artifisial",
		},
	}

	created := 0
	skipped := 0

	// delete all existing records before seeding
	db.Exec("DELETE FROM bidang_studi")

	//// Insert all records using FirstOrCreate
	for _, bidangStudi := range bidangStudiData {
		var result models.BidangStudi
		err := db.Where("code = ?", bidangStudi.Code).FirstOrCreate(&result, bidangStudi).Error
		if err != nil {
			return fmt.Errorf("failed to seed bidang_studi with code %s: %w", bidangStudi.Code, err)
		}

		// Check if it was created or already existed
		if result.CreatedAt.Equal(result.UpdatedAt) && result.Code == bidangStudi.Code {
			created++
		} else {
			skipped++
		}
	}

	fmt.Printf("✓ Bidang Studi seeding completed: %d created, %d skipped (already exist)\n", created, skipped)
	return nil
}
