package seeds

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"pustaka-backend/models"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BookCSVRow represents a row from the books.csv file
type BookCSVRow struct {
	BidangStudiCode string
	JenjangCode     string
	CurriculumCode  string
	Kelas           string
	Pages           int
	LKSPrice        float64
	PGPrice         float64
	Periode         int
	Year            string
	MerkCode        string
	Stock           int
	PublisherCode   string // Optional, defaults to "GP" if empty
}

// BooksSeeder seeds the books table from CSV file
// Each CSV row generates 1 or 2 book records:
// - Always creates LKS book
// - Only creates PG book if pg_price > 0
func BooksSeeder(db *gorm.DB) error {
	fmt.Println("📚 Seeding books table from CSV...")

	// Get the path to the CSV file
	_, filename, _, _ := runtime.Caller(0)
	seedsDir := filepath.Dir(filename)
	csvPath := filepath.Join(seedsDir, "files", "books.csv")

	// Read CSV file
	rows, err := readBooksCSV(csvPath)
	if err != nil {
		return fmt.Errorf("failed to read books CSV: %w", err)
	}

	if len(rows) == 0 {
		fmt.Println("⚠️  No data rows found in books.csv")
		return nil
	}

	// Load lookup maps for foreign keys
	bidangStudiMap, err := loadBidangStudiMap(db)
	if err != nil {
		return err
	}

	jenjangStudiMap, err := loadJenjangStudiMap(db)
	if err != nil {
		return err
	}

	curriculumMap, err := loadCurriculumMap(db)
	if err != nil {
		return err
	}

	jenisBukuMap, err := loadJenisBukuMap(db)
	if err != nil {
		return err
	}

	merkBukuMap, err := loadMerkBukuMap(db)
	if err != nil {
		return err
	}

	publisherMap, err := loadPublisherMap(db)
	if err != nil {
		return err
	}

	// Get LKS and PG UUIDs
	lksID, lksExists := jenisBukuMap["LKS"]
	pgID, pgExists := jenisBukuMap["PG"]

	if !lksExists || !pgExists {
		return fmt.Errorf("LKS or PG not found in jenis_buku table. Please run jenis_buku seeder first")
	}

	// Validate default publisher 'GP' exists (used when CSV doesn't specify publisher)
	if _, gpExists := publisherMap["GP"]; !gpExists {
		return fmt.Errorf("default publisher 'GP' not found in publishers table")
	}

	// Delete existing books (optional - you can comment this out if you want to keep existing data)
	fmt.Println("🗑️  Deleting existing books...")
	if err := db.Exec("DELETE FROM books").Error; err != nil {
		return fmt.Errorf("failed to delete existing books: %w", err)
	}

	created := 0
	skipped := 0
	errors := []string{}

	// Process each CSV row
	for i, row := range rows {
		lineNum := i + 2 // +2 because row 1 is header

		// Lookup foreign keys
		bidangStudiID, bidangStudiName, bsErr := lookupBidangStudi(bidangStudiMap, row.BidangStudiCode)
		if bsErr != nil {
			errors = append(errors, fmt.Sprintf("Line %d: %s", lineNum, bsErr.Error()))
			skipped += 2
			continue
		}

		jenjangStudiID, jsErr := lookupCode(jenjangStudiMap, row.JenjangCode, "jenjang_studi")
		if jsErr != nil {
			errors = append(errors, fmt.Sprintf("Line %d: %s", lineNum, jsErr.Error()))
			skipped += 2
			continue
		}

		curriculumID, cErr := lookupCode(curriculumMap, row.CurriculumCode, "curriculum")
		if cErr != nil {
			errors = append(errors, fmt.Sprintf("Line %d: %s", lineNum, cErr.Error()))
			skipped += 2
			continue
		}

		merkBukuID, mbErr := lookupCode(merkBukuMap, row.MerkCode, "merk_buku")
		if mbErr != nil {
			errors = append(errors, fmt.Sprintf("Line %d: %s", lineNum, mbErr.Error()))
			skipped += 2
			continue
		}

		// Determine publisher ID (from CSV or default 'GP')
		publisherCode := row.PublisherCode
		if publisherCode == "" {
			publisherCode = "GP" // Default publisher
		}
		publisherID, pubErr := lookupCode(publisherMap, publisherCode, "publisher")
		if pubErr != nil {
			errors = append(errors, fmt.Sprintf("Line %d: %s", lineNum, pubErr.Error()))
			skipped += 2
			continue
		}

		// Create LKS book
		lksBook := models.Book{
			Name:           bidangStudiName,
			Year:           row.Year,
			Periode:        row.Periode,
			Stock:          row.Stock,
			NoPages:        row.Pages,
			Kelas:          &row.Kelas,
			MerkBukuID:     &merkBukuID,
			JenisBukuID:    &lksID,
			JenjangStudiID: &jenjangStudiID,
			BidangStudiID:  &bidangStudiID,
			CurriculumID:   &curriculumID,
			PublisherID:    &publisherID,
			Price:          row.LKSPrice,
		}

		if err := db.Create(&lksBook).Error; err != nil {
			errors = append(errors, fmt.Sprintf("Line %d (LKS): %s", lineNum, err.Error()))
			skipped++
		} else {
			created++
		}

		// Create PG book only if pg_price is specified (> 0)
		if row.PGPrice > 0 {
			pgBook := models.Book{
				Name:           bidangStudiName,
				Year:           row.Year,
				Periode:        row.Periode,
				Stock:          row.Stock,
				NoPages:        row.Pages,
				Kelas:          &row.Kelas,
				MerkBukuID:     &merkBukuID,
				JenisBukuID:    &pgID,
				JenjangStudiID: &jenjangStudiID,
				BidangStudiID:  &bidangStudiID,
				CurriculumID:   &curriculumID,
				PublisherID:    &publisherID,
				Price:          row.PGPrice,
			}

			if err := db.Create(&pgBook).Error; err != nil {
				errors = append(errors, fmt.Sprintf("Line %d (PG): %s", lineNum, err.Error()))
				skipped++
			} else {
				created++
			}
		}
	}

	// Report errors if any
	if len(errors) > 0 {
		fmt.Println("⚠️  Errors encountered:")
		for _, e := range errors {
			fmt.Printf("   - %s\n", e)
		}
	}

	fmt.Printf("✅ Books seeding completed: %d created, %d skipped\n", created, skipped)
	return nil
}

// readBooksCSV reads and parses the books CSV file
func readBooksCSV(path string) ([]BookCSVRow, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Allow variable number of fields per record (some rows may not have publisher_code)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return []BookCSVRow{}, nil // Only header or empty
	}

	var rows []BookCSVRow
	for i, record := range records[1:] { // Skip header
		if len(record) < 11 {
			fmt.Printf("⚠️  Skipping row %d: not enough columns\n", i+2)
			continue
		}

		pages, _ := strconv.Atoi(strings.TrimSpace(record[4]))
		lksPrice, _ := strconv.ParseFloat(strings.TrimSpace(record[5]), 64)
		pgPrice, _ := strconv.ParseFloat(strings.TrimSpace(record[6]), 64)
		periode, _ := strconv.Atoi(strings.TrimSpace(record[7]))
		stock, _ := strconv.Atoi(strings.TrimSpace(record[10]))

		// Read optional publisher_code (12th column)
		publisherCode := ""
		if len(record) >= 12 {
			publisherCode = strings.TrimSpace(record[11])
		}

		rows = append(rows, BookCSVRow{
			BidangStudiCode: strings.TrimSpace(record[0]),
			JenjangCode:     strings.TrimSpace(record[1]),
			CurriculumCode:  strings.TrimSpace(record[2]),
			Kelas:           strings.TrimSpace(record[3]),
			Pages:           pages,
			LKSPrice:        lksPrice,
			PGPrice:         pgPrice,
			Periode:         periode,
			Year:            strings.TrimSpace(record[8]),
			MerkCode:        strings.TrimSpace(record[9]),
			Stock:           stock,
			PublisherCode:   publisherCode,
		})
	}

	return rows, nil
}

// BidangStudiInfo holds both ID and Name
type BidangStudiInfo struct {
	ID   uuid.UUID
	Name string
}

// loadBidangStudiMap loads bidang_studi code -> {ID, Name} mapping
func loadBidangStudiMap(db *gorm.DB) (map[string]BidangStudiInfo, error) {
	var records []models.BidangStudi
	if err := db.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to load bidang_studi: %w", err)
	}

	m := make(map[string]BidangStudiInfo)
	for _, r := range records {
		m[r.Code] = BidangStudiInfo{ID: r.ID, Name: r.Name}
	}
	return m, nil
}

// loadJenjangStudiMap loads jenjang_studi code -> UUID mapping
func loadJenjangStudiMap(db *gorm.DB) (map[string]uuid.UUID, error) {
	var records []models.JenjangStudi
	if err := db.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to load jenjang_studi: %w", err)
	}

	m := make(map[string]uuid.UUID)
	for _, r := range records {
		m[r.Code] = r.ID
	}
	return m, nil
}

// loadCurriculumMap loads curriculum code -> UUID mapping
func loadCurriculumMap(db *gorm.DB) (map[string]uuid.UUID, error) {
	var records []models.Curriculum
	if err := db.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to load curriculum: %w", err)
	}

	m := make(map[string]uuid.UUID)
	for _, r := range records {
		m[r.Code] = r.ID
	}
	return m, nil
}

// loadJenisBukuMap loads jenis_buku code -> UUID mapping
func loadJenisBukuMap(db *gorm.DB) (map[string]uuid.UUID, error) {
	var records []models.JenisBuku
	if err := db.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to load jenis_buku: %w", err)
	}

	m := make(map[string]uuid.UUID)
	for _, r := range records {
		m[r.Code] = r.ID
	}
	return m, nil
}

// loadMerkBukuMap loads merk_buku code -> UUID mapping
func loadMerkBukuMap(db *gorm.DB) (map[string]uuid.UUID, error) {
	var records []models.MerkBuku
	if err := db.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to load merk_buku: %w", err)
	}

	m := make(map[string]uuid.UUID)
	for _, r := range records {
		m[r.Code] = r.ID
	}
	return m, nil
}

// loadPublisherMap loads publisher code -> UUID mapping
func loadPublisherMap(db *gorm.DB) (map[string]uuid.UUID, error) {
	var records []models.Publisher
	if err := db.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to load publishers: %w", err)
	}

	m := make(map[string]uuid.UUID)
	for _, r := range records {
		m[r.Code] = r.ID
	}
	return m, nil
}

// lookupBidangStudi looks up bidang_studi by code and returns both ID and Name
func lookupBidangStudi(m map[string]BidangStudiInfo, code string) (uuid.UUID, string, error) {
	info, exists := m[code]
	if !exists {
		return uuid.Nil, "", fmt.Errorf("bidang_studi code '%s' not found", code)
	}
	return info.ID, info.Name, nil
}

// lookupCode is a generic lookup function for code -> UUID
func lookupCode(m map[string]uuid.UUID, code string, tableName string) (uuid.UUID, error) {
	id, exists := m[code]
	if !exists {
		return uuid.Nil, fmt.Errorf("%s code '%s' not found", tableName, code)
	}
	return id, nil
}
