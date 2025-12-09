package seeds

import (
	"fmt"
	"os"
	"io"
	"time"
	"strconv"
	"encoding/csv"
	"pustaka-backend/models"
	"gorm.io/gorm"
)

// Helper function to create string pointer
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Helper function to create time pointer
func timePtr(s string) (*time.Time, error) {
    if s == "" {
        return nil, nil
    }

    layout := "2006-01-02 15:04:05.000"

    t, err := time.Parse(layout, s)
    if err != nil {
        return nil, err
    }

    return &t, nil
}

func ReadSalesAssociatesCSV(filePath string, db *gorm.DB) ([]models.SalesAssociate, error) {
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

	var results []models.SalesAssociate

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// fmt.Printf("%+v\n", record)

		// Get City ID
		var city models.City
		cityCode := record[4]
		if err := db.Where("code = ?", cityCode).First(&city).Error; err != nil {
			return nil, err
		}
		cityID := city.ID

		joinDate, err := time.Parse("2006-01-02 15:04:05.000", record[9])
		if err != nil {
			return nil, err
		}

		endJoinDate, err := timePtr(record[10])
		if err != nil {
			return nil, err
		}

		discount, err := strconv.ParseFloat(record[11], 64)
		if err != nil {
			return nil, err
		}

		results = append(results, models.SalesAssociate{
			Code:            record[0],
			Name:            record[1],
			Address:         record[2] + " " + record[3],
			CityID:          &cityID,
			Phone1:          record[5],
			Phone2:          strPtr(record[6]),
			JenisPembayaran: record[8],
			JoinDate:        joinDate,
			EndJoinDate:     endJoinDate,
			Discount:        discount,
		})
	}

	return results, nil
}


// SalesAssociateSeeder seeds the sales_associate table with initial data
// Uses FirstOrCreate for conflict resolution - won't create duplicates
func SalesAssociateSeeder(db *gorm.DB) error {
	fmt.Println("📝 Seeding sales_associates table...")

	salesAssociatesData, err := ReadSalesAssociatesCSV("./seeds/files/sales_associates.csv", db)
	if err != nil {
		return fmt.Errorf("failed to load sales_associates.csv: %w", err.Error())
	}

	fmt.Println("Loaded:", len(salesAssociatesData))
	// for _, salesAssociate := range salesAssociatesData {
	// 	fmt.Printf("%+v\n", salesAssociate.Code)
	// }

	created := 0
	skipped := 0

	// delete all existing records before seeding
	db.Exec("DELETE FROM sales_associates")

	// Insert all records using FirstOrCreate
	for _, sales_associate := range salesAssociatesData {
		var result models.SalesAssociate
		err := db.Where("code = ?", sales_associate.Code).FirstOrCreate(&result, sales_associate).Error
		if err != nil {
			return fmt.Errorf("failed to seed sales_associate with code %s: %w", sales_associate.Code, err)
		}

		// Check if it was created or already existed
		if result.CreatedAt.Equal(result.UpdatedAt) && result.Code == sales_associate.Code {
			created++
		} else {
			skipped++
		}
	}

	fmt.Printf("✓ SalesAssociate seeding completed: %d created, %d skipped (already exist)\n", created, skipped)
	return nil
}
