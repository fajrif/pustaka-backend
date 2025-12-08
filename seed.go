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
		{"kelas", func() error { return seeds.KelasSeeder(config.DB) }},
		{"bidang_studi", func() error { return seeds.BidangStudiSeeder(config.DB) }},
		{"cities", func() error { return seeds.CitiesSeeder(config.DB) }},
		{"jenjang_studi", func() error { return seeds.JenjangStudiSeeder(config.DB) }},
		{"expeditions", func() error { return seeds.ExpeditionsSeeder(config.DB) }},
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
	case "kelas":
		err = seeds.KelasSeeder(config.DB)
	case "bidang_studi":
		err = seeds.BidangStudiSeeder(config.DB)
	case "jenjang_studi":
		err = seeds.JenjangStudiSeeder(config.DB)
	case "cities":
		err = seeds.CitiesSeeder(config.DB)
	case "expeditions":
	  err = seeds.ExpeditionsSeeder(config.DB)
	default:
		log.Fatalf("❌ Unknown seeder: %s\n\nAvailable seeders:\n  - jenis_buku, kelas, bidang_studi, jenjang_studi, cities, expeditions\n", name)
	}

	if err != nil {
		log.Fatalf("❌ Error running seeder: %v", err)
	}
}
