package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"sukimise/internal/config"
	"sukimise/internal/models"
	"sukimise/internal/repositories"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// getStandardCategories returns a list of commonly used categories in Japan
func getStandardCategories() []models.CategoryCustomization {
	return []models.CategoryCustomization{
		{ CategoryName: "和食", Icon: stringPtr("🍱"), Color: stringPtr("#8B4513")},
		{ CategoryName: "洋食", Icon: stringPtr("🍽️"), Color: stringPtr("#FF6B6B")},
		{ CategoryName: "中華", Icon: stringPtr("🥢"), Color: stringPtr("#FF4444")},
		{ CategoryName: "カフェ", Icon: stringPtr("☕"), Color: stringPtr("#8B4513")},
		{ CategoryName: "居酒屋", Icon: stringPtr("🍺"), Color: stringPtr("#FF6347")},
		{ CategoryName: "ラーメン", Icon: stringPtr("🍜"), Color: stringPtr("#FF4500")},
		{ CategoryName: "寿司", Icon: stringPtr("🍣"), Color: stringPtr("#4682B4")},
		{ CategoryName: "焼肉", Icon: stringPtr("🥩"), Color: stringPtr("#DC143C")},
		{ CategoryName: "スイーツ", Icon: stringPtr("🍰"), Color: stringPtr("#FF69B4")},
		{ CategoryName: "パン", Icon: stringPtr("🥖"), Color: stringPtr("#DEB887")},
		{ CategoryName: "スーパー", Icon: stringPtr("🛒"), Color: stringPtr("#228B22")},
		{ CategoryName: "本屋", Icon: stringPtr("📚"), Color: stringPtr("#8B4513")},
		{ CategoryName: "衣料品", Icon: stringPtr("👕"), Color: stringPtr("#FF69B4")},
		{ CategoryName: "雑貨", Icon: stringPtr("🎁"), Color: stringPtr("#DA70D6")},
		{ CategoryName: "温泉", Icon: stringPtr("♨️"), Color: stringPtr("#20B2AA")},
		{ CategoryName: "公園", Icon: stringPtr("🌳"), Color: stringPtr("#228B22")},
		{ CategoryName: "博物館", Icon: stringPtr("🏛️"), Color: stringPtr("#8B4513")},
		{ CategoryName: "図書館", Icon: stringPtr("📖"), Color: stringPtr("#4682B4")},
		{ CategoryName: "ホテル", Icon: stringPtr("🏨"), Color: stringPtr("#4169E1")},
		{ CategoryName: "ガソリンスタンド", Icon: stringPtr("⛽"), Color: stringPtr("#FF4500")},
		{ CategoryName: "駐車場", Icon: stringPtr("🅿️"), Color: stringPtr("#808080")},
		{ CategoryName: "駅", Icon: stringPtr("🚉"), Color: stringPtr("#4682B4")},
		{ CategoryName: "オフィス", Icon: stringPtr("🏢"), Color: stringPtr("#708090")},
		{ CategoryName: "工場", Icon: stringPtr("🏭"), Color: stringPtr("#696969")},
		{ CategoryName: "その他", Icon: stringPtr("📍"), Color: stringPtr("#808080")},
	}
}

func stringPtr(s string) *string {
	return &s
}

func main() {
	log.Println("Starting category seeding...")

	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Info: .env file not found (normal in Docker environment)")
	}

	// Load configuration
	cfg := config.LoadConfig()
	if cfg.Database.URL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create category customization repository
	categoryRepo := repositories.NewCategoryCustomizationRepository(db)

	// Get standard categories
	standardCategories := getStandardCategories()

	log.Printf("Found %d standard categories to seed", len(standardCategories))

	// Create categories
	createdCount := 0
	updatedCount := 0
	skippedCount := 0

	for _, stdCategory := range standardCategories {
		err := createOrUpdateCategory(categoryRepo, stdCategory)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("Category '%s' already exists, skipping", stdCategory.CategoryName)
				skippedCount++
			} else {
				log.Printf("Failed to create category '%s': %v", stdCategory.CategoryName, err)
				continue
			}
		} else {
			log.Printf("Successfully created category: %s", stdCategory.CategoryName)
			createdCount++
		}
	}

	log.Printf("Category seeding completed:")
	log.Printf("  Created: %d", createdCount)
	log.Printf("  Updated: %d", updatedCount)
	log.Printf("  Skipped: %d", skippedCount)
	log.Printf("  Total: %d", len(standardCategories))
}

// createOrUpdateCategory creates or updates a category customization
func createOrUpdateCategory(repo repositories.CategoryCustomizationRepositoryInterface, category models.CategoryCustomization) error {
	// Check if category already exists
	existing, err := repo.GetByCategoryName(category.CategoryName)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing category: %v", err)
	}

	if existing != nil {
		// Category exists, skip update
		return nil
	} else {
		// Category doesn't exist, create new one
		customization := &models.CategoryCustomization{
			CategoryName: category.CategoryName,
			Icon:        category.Icon,
			Color:       category.Color,
		}

		return repo.Create(customization)
	}
}