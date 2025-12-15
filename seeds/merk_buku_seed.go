package seeds

import (
	"fmt"
	"os"
	"io"
	"encoding/csv"
	"pustaka-backend/models"
	"gorm.io/gorm"
)

func ReadMerkBukuCSV(filePath string) ([]models.MerkBuku, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // <-- IMPORTANT
	reader.FieldsPerRecord = -1 // allow variable columns

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var results []models.MerkBuku

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		results = append(results, models.MerkBuku{
			Code:            record[0],
			Name:            record[1],
		})
	}

	return results, nil
}

// MerkBukuSeeder seeds the merk_buku table with initial data
// Uses FirstOrCreate for conflict resolution - won't create duplicates
func MerkBukuSeeder(db *gorm.DB) error {
	fmt.Println("📝 Seeding merk_buku table...")

	merkBukuData, err := ReadMerkBukuCSV("./seeds/files/merk_buku.csv")
	if err != nil {
		return fmt.Errorf("failed to load merk_buku.csv: %w", err.Error())
	}

	fmt.Println("Loaded:", len(merkBukuData))
	// for _, salesAssociate := range salesAssociatesData {
	// 	fmt.Printf("%+v\n", salesAssociate.Code)
	// }

	created := 0
	skipped := 0

	// delete all existing records before seeding
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

	fmt.Printf("✓ Merk Buku seeding completed: %d created, %d skipped (already exist)\n", created, skipped)
	return nil
}
