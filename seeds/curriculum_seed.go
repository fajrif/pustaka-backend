package seeds

import (
	"fmt"
	"pustaka-backend/models"

	"gorm.io/gorm"
)

// CurriculumSeeder seeds the curriculum table with initial data
// Uses FirstOrCreate for conflict resolution - won't create duplicates
func CurriculumSeeder(db *gorm.DB) error {
	fmt.Println("Seeding curriculum table...")

	curriculumData := []models.Curriculum{
		{
			Code: "K13",
			Name: "K13",
		},
		{
			Code: "MER",
			Name: "Merdeka",
		},
		{
			Code: "NAS",
			Name: "Nasional",
		},
	}

	created := 0
	skipped := 0

	// Delete all existing records before seeding
	db.Exec("DELETE FROM curriculum")

	// Insert all records using FirstOrCreate
	for _, curriculum := range curriculumData {
		var result models.Curriculum
		err := db.Where("code = ?", curriculum.Code).FirstOrCreate(&result, curriculum).Error
		if err != nil {
			return fmt.Errorf("failed to seed curriculum with code %s: %w", curriculum.Code, err)
		}

		// Check if it was created or already existed
		if result.CreatedAt.Equal(result.UpdatedAt) && result.Code == curriculum.Code {
			created++
		} else {
			skipped++
		}
	}

	fmt.Printf("Curriculum seeding completed: %d created, %d skipped (already exist)\n", created, skipped)
	return nil
}
