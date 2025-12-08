package seeds

import (
	"fmt"
	"pustaka-backend/models"
	"gorm.io/gorm"
)

// ExpeditionSeeder seeds the expedition table with initial data
// Uses FirstOrCreate for conflict resolution - won't create duplicates
func ExpeditionsSeeder(db *gorm.DB) error {
	fmt.Println("📝 Seeding expedition table...")

	// Helper function to create string pointer
	strPtr := func(s string) *string {
		return &s
	}

	expeditionData := []models.Expedition{
		{
			Code: "ARP",
			Name: "ARYA PRIMA",
			Address: "JL.JEND.SUDIRMAN BY PASS NO . 61",
			Area: strPtr("CIKOKOL"),
			CityCode: "TNG",
			Phone1: "021 55740384",
		},
		{
			Code: "AKS",
			Name: "ANGKASA",
			Address: "RUKO MANGGA DUA PLAZA BLOK B NO 6",
			Area: strPtr("MANGGA DUA RAYA"),
			CityCode: "JKB",
			Phone1: "021 6120709",
		},
		{
			Code: "ALS",
			Name: "ANTAR LINTAS SUMATERA",
			Address: "JL.DAAN MOGOT KM 24 NO 25",
			Area: strPtr("TANAH TINGGI"),
			CityCode: "TNG",
		},
		{
			Code: "AP",
			Name: "ATLAS PRIMA",
			Address: "RUKO TAMAN PALEM FANTASI W-21",
			Area: strPtr("CENGKARENG"),
			CityCode: "JKB",
			Phone1: "021 68072588",
		},
		{
			Code: "AWR",
			Name: "AWR CARGO",
			Address: "JL. KEBON KACANG 1 NO52",
			Area: strPtr("TANAH ABANG"),
			CityCode: "JKP",
			Phone1: "081382274868",
			Phone2: strPtr("81311436280"),
		},
		{
			Code: "BE",
			Name: "BERLIAN EXPRES",
			Address: "JATI BUNDER ( ANEKA BETON PANGKALAN TRUCK)",
			Area: strPtr("TANAH ABANG"),
			CityCode: "JKP",
			Phone1: "021 6881 8963",
		},
		{
			Code: "BTN",
			Name: "BETELNUTS",
			Address: "KOMP. RUKO FANTASI TAMAN PALEM LESTARI BLOK U 26",
			CityCode: "JKB",
			Phone1: "33099836",
		},
		{
			Code: "CEN",
			Name: "CENTRAL EXPRESS",
			Address: "JL.T.B.ANGKE NO 6 AL-AM (SEBELAH KOMP.DUTAMAS )",
			Area: strPtr("JAKARTA BARAT"),
			CityCode: "DKI",
			Phone1: "021 5676210",
		},
		{
			Code: "CJE",
			Name: "CITRAMAS JAYA EXPRES",
			Address: "JL.KAMPUNG BANDAN NO.1 PINTU 25",
			CityCode: "JKP",
			Phone1: "691 5335",
		},
		{
			Code: "CS",
			Name: "CINTA SAUDARA",
			Address: "JL KS TUBUN IV NO.1 RT 001/07",
			Area: strPtr("SLIPI"),
			CityCode: "JKB",
			Phone1: "081296535400",
		},
		{
			Code: "DJ",
			Name: "DJ EXPRES",
			Address: "JL. RAYA BOGOR KM 22 NO 51",
			Area: strPtr("KRAMAT JATI"),
			CityCode: "JKT",
			Phone1: "081280979190",
		},
		{
			Code: "DKT",
			Name: "DAKOTA",
			Address: "JL PAHLAWAN REVOLUSI 10.123",
			Area: strPtr("PONDOK BAMBU"),
			CityCode: "JKT",
			Phone1: "021 8616987",
		},
		{
			Code: "FI",
			Name: "FAJAR INDAH",
			Address: "JL.JATI BARU BENGKEL NO : 20",
			Area: strPtr("TANAH ABANG"),
			CityCode: "JKP",
			Phone1: "021 3449212",
		},
		{
			Code: "FRC",
			Name: "FAMILY RAYA CERIA",
			Address: "TERMINAL PORIS  PLAWAD",
			CityCode: "TNG",
			Phone1: "021 70515248",
		},
		{
			Code: "GSN",
			Name: "GOSEND",
			Address: "CIPUTAT",
			CityCode: "TNG",
			Phone1: "085890663454",
		},
		{
			Code: "HSE",
			Name: "HASRAT SAMUDERA EXPRESS",
			Address: "JL. JATI BARU RAYA NO 66",
			Area: strPtr("TANAH ABANG"),
			CityCode: "JKP",
			Phone1: "021 3862385",
		},
		{
			Code: "IC",
			Name: "INDAH LOGISTIK CARGO CIPUTAT",
			Address: "JL IR H JUANDA NO 129B CEMPAKA PUTIH",
			Area: strPtr("KEC CIPUTAT TIMUR"),
			CityCode: "TNG",
			Phone1: "081119641178",
		},
		{
			Code: "JNE",
			Name: "JNE",
			Address: "JL. WR SUPRATMAN NO. 4 RENGAS",
			Area: strPtr("KEC CIPUTAT TIMUR"),
			CityCode: "TNG",
			Phone1: "(022) 736398",
		},
		{
			Code: "JI",
			Name: "CV JASA IBU",
			Address: "JL. KH. MAS MANSYUR NO.23 A",
			CityCode: "JKP",
			Phone1: "(021) 3190 4377",
			Phone2: strPtr("31904957"),
		},
		{
			Code: "JK",
			Name: "JASA KITA",
			Address: "JL.JATI BARU RAYA NO: 59",
			Area: strPtr("TANAH ABANG"),
			CityCode: "JKP",
			Phone1: "(021) 3919407",
		},
		{
			Code: "JMT",
			Name: "JAYA MUSTIKA TRANSPORT",
			Address: "TAMAN PALEM LESTARI RUKO FANTASI BLOK Z3 NO 10",
			Area: strPtr("CENGKARENG"),
			CityCode: "JKB",
			Phone1: "55963815",
			Phone2: strPtr("71283991"),
		},
		{
			Code: "JNT",
			Name: "JNT EXPRESS CP KAMPUNG UTAN",
			Address: "JL W.R SUPRATMAN RT 004 RW 006 CEMPAKA PUTIH",
			Area: strPtr("KEC. CIPUTAT TIMUR"),
			CityCode: "TNG",
			Phone1: "7448762",
		},
		{
			Code: "KS",
			Name: "KS",
			Address: "JL.PANGERAN JAYA KARTA",
			CityCode: "JKP",
			Phone1: "(021) 6596817",
		},
		{
			Code: "LTR",
			Name: "LANTRA",
			Address: "CIKOKOL",
			CityCode: "TNG",
		},
		{
			Code: "MPS",
			Name: "METRO PARCEL SERVICE",
			Address: "JL.RAWA GELAM I NO 6",
			Area: strPtr("PULO GADUNG"),
			CityCode: "JKT",
		},
		{
			Code: "MPX",
			Name: "MAKMUR PRIMA XPRESS",
			Address: "TAMAN PALEM LESTARI BLOK J NO 28A",
			Area: strPtr("KPMPLEK GALAKSI"),
			CityCode: "JKB",
			Phone1: "(021) 55959750",
		},
		{
			Code: "PCS",
			Name: "PANCA KOBRA SAKTI",
			Address: "JL ALAYDRUS NO 13",
			CityCode: "JKP",
		},
		{
			Code: "PI",
			Name: "PARUNG INDAH",
			Address: "DEPAN SEKPOLWAN",
			CityCode: "JKS",
		},
		{
			Code: "PIE",
			Name: "PRIMA INDAH EXPRES",
			Address: "JL.JATI BARU NO.87 RT 03/01",
			Area: strPtr("TANAH ABANG"),
			CityCode: "JKP",
			Phone1: "(021) 3159078",
			Phone2: strPtr("0815 9506108"),
		},
		{
			Code: "PJ",
			Name: "PRIMA JASA",
			Address: "DEPAN SEKPOLWAN",
			CityCode: "JKS",
		},
		{
			Code: "PMT",
			Name: "PMTOH/CV KHARISMA",
			Address: "JL.KH MANSYUR TANAH ABANG",
			CityCode: "JKP",
			Phone1: "(021) 3143563",
		},
		{
			Code: "RON",
			Name: "RONA INDAH TRANS",
			Address: "CIKOKOL",
			CityCode: "TNG",
			Phone2: strPtr("081369492614"),
		},
		{
			Code: "SD",
			Name: "SINAR DAGANG",
			Address: "JL.JATI BARU NO 1 SAMPING HOTEL PARMIN",
			Area: strPtr("TANAH ABANG"),
			CityCode: "JKP",
			Phone1: "(021) 31924038",
			Phone2: strPtr("082260636664"),
		},
		{
			Code: "SUT",
			Name: "SUMBAR AMANDA TRANS",
			Address: "TAMAN KEDOYA BARU BLOK A 15",
			CityCode: "JKB",
			Phone1: "(021) 70774348",
		},
		{
			Code: "TIA",
			Name: "TELAGA BIRU GROUP",
			Address: "JL. JENDRAL SUDIRMAN BY PASS TANGERANG",
			CityCode: "TNG",
			Phone1: "5544666",
			Phone2: strPtr("081386885051"),
		},
		{
			Code: "TJS",
			Name: "TARUNA JAYASARANA SEMPURNA",
			Address: "RUKO GALAXI BLOK K NO 1 TAMAN PALEM LESTARI",
			Area: strPtr("JL KAMAL RAYA OUTER RING ROAD"),
			CityCode: "JKB",
			Phone1: "021 55959755",
		},
		{
			Code: "TK",
			Name: "TIKI",
			Address: "DSJK",
			Area: strPtr("CDSJFDS"),
			CityCode: "DKI",
		},
		{
			Code: "TRU",
			Name: "TRUCK PASAR INDUK",
			Address: "PASAR INDUK KRAMAT JATI",
			Area: strPtr("JAKARTA TIMUR"),
			CityCode: "BKS",
		},
		{
			Code: "UB",
			Name: "USAHA BERLIAN",
			Address: "RUKO TAMAN PALEM LESTARI BLOK A30 NO 11",
			Area: strPtr("CENGKARENG"),
			CityCode: "JKB",
			Phone1: "55952201",
		},
		{
			Code: "UD",
			Name: "UDAYANA",
			Address: "JL.JEND.SUDIRMAN NO 69 CIKOKOL",
			CityCode: "TNG",
			Phone1: "(021) 70287505",
		},
		{
			Code: "US",
			Name: "USOR2",
			Address: "DEPAN TPI",
			CityCode: "JKT",
		},
		{
			Code: "WBW",
			Name: "WIBOWO",
			Address: "JL.P.JAYAKARTA NO.117 BLOK.C-27",
			CityCode: "JKP",
			Phone1: "600 7383",
		},
		{
			Code: "WHY",
			Name: "WAHYU EXPRESS",
			Address: "JL.KARET TENGSIN NO: 9 ( SAMPING MENARA BATAVIA )",
			Area: strPtr("POOL TRUCK KARET"),
			CityCode: "JKP",
			Phone1: "(021) 7071870",
			Phone2: strPtr("0818798761"),
		},
		{
			Code: "WP",
			Name: "WALI PITUE",
			Address: "JL.KH MAS MANSYUR NO.104",
			Area: strPtr("TANAH ABANG"),
			CityCode: "JKP",
			Phone1: "3152991",
		},
		{
			Code: "POS",
			Name: "KANTOR POS INDONESIA",
			Address: "JL W R SUPRATMAN NO 13 RT 3 RW 2 PD RANJI",
			Area: strPtr("KEC. CIPUTAT TIMUR"),
			CityCode: "TNG",
			Phone1: "1500161",
		},
		{
			Code: "JC",
			Name: "JNT CARGO",
			Address: "JL W R SUPRATMAN NO. 30 010, RT/RW 006",
			Area: strPtr("RENGAS KEC CIPUTAT TIMUR TGR SELATAN"),
			CityCode: "TNG",
			Phone1: "87828520997",
		},
		{
			Code: "KE",
			Name: "KOBRA EXPRESS",
			Address: "JL MOCH SUYUDI S.H NO. 5",
			Area: strPtr("BSD"),
			CityCode: "TNG",
		},
		{
			Code: "LM",
			Name: "LALAMOVE",
			Address: "TANGERANG SELATAN",
			CityCode: "DKI",
		},
		{
			Code: "EPM",
			Name: "EFISIEN PUTRA MANDIRI",
			Address: "TERMINAL PORIS PLAWAD TANGERANG",
			CityCode: "TNG",
		},
		{
			Code: "BRK",
			Name: "BARAKA EKSPEDISI",
			Address: "JL IR H JUANDA NO 001 CIPUTAT TIMUR",
			CityCode: "TNG",
		},
	}

	created := 0
	skipped := 0

	// delete all existing records before seeding
	db.Exec("DELETE FROM expedition")

	// Insert all records using FirstOrCreate
	for _, expedition := range expeditionData {
		var result models.Expedition
		// Get City ID
		var city models.City
		if err := db.Where("code = ?", expedition.CityCode).First(&city).Error; err == nil {
				expedition.CityID = &city.ID
		}

		err := db.Where("code = ?", expedition.Code).FirstOrCreate(&result, expedition).Error
		if err != nil {
			return fmt.Errorf("failed to seed expedition with code %s: %w", expedition.Code, err)
		}

		// Check if it was created or already existed
		if result.CreatedAt.Equal(result.UpdatedAt) && result.Code == expedition.Code {
			created++
		} else {
			skipped++
		}
	}

	fmt.Printf("✓ Expeditions seeding completed: %d created, %d skipped (already exist)\n", created, skipped)
	return nil
}
