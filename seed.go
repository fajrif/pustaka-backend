package main

import (
	"fmt"
	"log"
	"os"
	"pustaka-backend/config"
	"pustaka-backend/seeds"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Connect to database
	fmt.Println("🔌 Connecting to database...")
	config.ConnectDB()
	fmt.Println("✓ Database connected")
	fmt.Println("")

	// Get command line argument to determine which seeder to run
	args := os.Args[1:]

	if len(args) == 0 {
		// Run all seeders
		fmt.Println("🌱 Running all seeders...")
		fmt.Println("==========================================")
		runAllSeeders()
	} else {
		// Run specific seeder
		seederName := args[0]
		fmt.Printf("🌱 Running seeder: %s\n", seederName)
		fmt.Println("==========================================")
		runSpecificSeeder(seederName)
	}

	fmt.Println("")
	fmt.Println("==========================================")
	fmt.Println("✅ Seeding completed successfully!")
}

// runAllSeeders runs all available seeders
func runAllSeeders() {
	seeders := []struct {
		name string
		fn   func() error
	}{
		{"jenis_buku", func() error { return seeds.JenisBukuSeeder(config.DB) }},
		// Add more seeders here as you create them
		// {"cities", func() error { return seeds.CitiesSeeder(config.DB) }},
		// {"expeditions", func() error { return seeds.ExpeditionsSeeder(config.DB) }},
	}

	for _, seeder := range seeders {
		fmt.Printf("\n▶ Running %s seeder...\n", seeder.name)
		if err := seeder.fn(); err != nil {
			log.Fatalf("❌ Error running %s seeder: %v", seeder.name, err)
		}
	}
}

// runSpecificSeeder runs a specific seeder by name
func runSpecificSeeder(name string) {
	var err error

	switch name {
	case "jenis_buku":
		err = seeds.JenisBukuSeeder(config.DB)
	// Add more cases here as you create more seeders
	// case "cities":
	// 	err = seeds.CitiesSeeder(config.DB)
	// case "expeditions":
	// 	err = seeds.ExpeditionsSeeder(config.DB)
	default:
		log.Fatalf("❌ Unknown seeder: %s\n\nAvailable seeders:\n  - jenis_buku\n", name)
	}

	if err != nil {
		log.Fatalf("❌ Error running seeder: %v", err)
	}
}
