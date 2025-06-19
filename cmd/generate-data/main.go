package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/baskint/bidding-analysis/internal/config"
	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/baskint/bidding-analysis/internal/store"
)

func main() {
	log.Println("Starting test data generation...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := store.NewPostgresDB(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize stores
	bidStore := store.NewBidStore(db)
	campaignStore := store.NewCampaignStore(db)

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	generator := &DataGenerator{
		bidStore:      bidStore,
		campaignStore: campaignStore,
	}

	// Generate test data
	if err := generator.GenerateAll(); err != nil {
		log.Fatalf("Failed to generate test data: %v", err)
	}

	log.Println("Test data generation completed successfully!")
}

type DataGenerator struct {
	bidStore      *store.BidStore
	campaignStore *store.CampaignStore
}

func (g *DataGenerator) GenerateAll() error {
	log.Println("Creating test user...")
	userID := g.createTestUser()

	log.Println("Generating campaigns...")
	campaigns, err := g.generateCampaigns(userID, 8)
	if err != nil {
		return fmt.Errorf("failed to generate campaigns: %w", err)
	}

	log.Println("Generating bid events...")
	for i, campaign := range campaigns {
		bidCount := 500 + rand.Intn(1000) // 500-1500 bids per campaign
		log.Printf("Generating %d bids for campaign %d/%d", bidCount, i+1, len(campaigns))

		if err := g.generateBidEvents(campaign.ID, bidCount); err != nil {
			return fmt.Errorf("failed to generate bid events for campaign %s: %w", campaign.ID, err)
		}
	}

	log.Println("All test data generated successfully!")
	return nil
}

func (g *DataGenerator) createTestUser() uuid.UUID {
	return uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
}

func (g *DataGenerator) createUser(user *models.User) error {
	query := `
		INSERT INTO users (id, username, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO NOTHING`

	_, err := g.campaignStore.GetDB().Exec(query, user.ID, user.Username, user.PasswordHash)
	return err
}

func (g *DataGenerator) generateCampaigns(userID uuid.UUID, count int) ([]*models.Campaign, error) {
	campaignTypes := []struct {
		name        string
		budget      float64
		targetCPA   float64
		performance string // high, medium, low
	}{
		{"Premium Mobile App Campaign", 50000, 15.0, "high"},
		{"E-commerce Holiday Sale", 75000, 25.0, "high"},
		{"Brand Awareness Desktop", 30000, 40.0, "medium"},
		{"Retargeting Campaign", 20000, 12.0, "high"},
		{"Local Services Marketing", 15000, 20.0, "medium"},
		{"SaaS Lead Generation", 40000, 35.0, "medium"},
		{"Gaming App Install", 25000, 8.0, "low"},
		{"Fashion Brand Promotion", 35000, 30.0, "low"},
	}

	var campaigns []*models.Campaign

	for i := 0; i < count && i < len(campaignTypes); i++ {
		ct := campaignTypes[i]

		campaign := &models.Campaign{
			Name:        ct.name,
			UserID:      userID,
			Status:      "active",
			Budget:      &ct.budget,
			DailyBudget: func() *float64 { db := ct.budget / 30; return &db }(),
			TargetCPA:   &ct.targetCPA,
		}

		if err := g.campaignStore.CreateCampaign(campaign); err != nil {
			return nil, err
		}

		campaigns = append(campaigns, campaign)
		log.Printf("Created campaign: %s (ID: %s)", campaign.Name, campaign.ID)
	}

	return campaigns, nil
}

func (g *DataGenerator) generateBidEvents(campaignID uuid.UUID, count int) error {
	// Define user segments with different behaviors
	segments := []struct {
		id                   string
		category             string
		engagementScore      float64
		conversionProb       float64
		winRateMultiplier    float64
		conversionMultiplier float64
	}{
		{"premium_users", "premium", 0.85, 0.18, 1.3, 2.0},
		{"engaged_users", "high_engagement", 0.75, 0.12, 1.1, 1.5},
		{"regular_users", "standard", 0.55, 0.08, 1.0, 1.0},
		{"price_sensitive", "budget", 0.45, 0.06, 0.8, 0.8},
		{"new_users", "acquisition", 0.35, 0.04, 0.7, 0.6},
	}

	deviceTypes := []string{"mobile", "desktop", "tablet"}
	browsers := []string{"Chrome", "Safari", "Firefox", "Edge"}
	osTypes := []string{"Windows", "macOS", "iOS", "Android", "Linux"}

	countries := []string{"US", "CA", "UK", "DE", "FR", "AU", "JP"}
	regions := map[string][]string{
		"US": {"California", "New York", "Texas", "Florida", "Illinois"},
		"CA": {"Ontario", "Quebec", "British Columbia", "Alberta"},
		"UK": {"England", "Scotland", "Wales", "Northern Ireland"},
		"DE": {"Bavaria", "North Rhine-Westphalia", "Baden-Württemberg"},
		"FR": {"Île-de-France", "Provence-Alpes-Côte d'Azur", "Auvergne-Rhône-Alpes"},
		"AU": {"New South Wales", "Victoria", "Queensland", "Western Australia"},
		"JP": {"Tokyo", "Osaka", "Kanagawa", "Aichi"},
	}

	// Generate events over the last 30 days
	startTime := time.Now().AddDate(0, 0, -30)

	for i := 0; i < count; i++ {
		// Pick random segment
		segment := segments[rand.Intn(len(segments))]

		// Pick random location
		country := countries[rand.Intn(len(countries))]
		regionList := regions[country]
		region := regionList[rand.Intn(len(regionList))]

		// Pick random device
		deviceType := deviceTypes[rand.Intn(len(deviceTypes))]
		browser := browsers[rand.Intn(len(browsers))]
		os := osTypes[rand.Intn(len(osTypes))]

		// Generate random timestamp in the last 30 days
		randomDuration := time.Duration(rand.Int63n(int64(30 * 24 * time.Hour)))
		timestamp := startTime.Add(randomDuration)

		// Generate floor price (typically $0.50 - $5.00)
		floorPrice := 0.5 + rand.Float64()*4.5

		// Generate bid price (floor * 1.1 to 2.5)
		bidMultiplier := 1.1 + rand.Float64()*1.4
		bidPrice := floorPrice * bidMultiplier

		// Determine if bid won (based on segment and bid aggressiveness)
		baseWinRate := 0.3 + (bidMultiplier-1.1)*0.4 // Higher bids win more
		winRate := baseWinRate * segment.winRateMultiplier
		won := rand.Float64() < winRate

		var winPrice *float64
		if won {
			// Win price is usually 85-95% of bid price
			wp := bidPrice * (0.85 + rand.Float64()*0.1)
			winPrice = &wp
		}

		// Determine conversion (only if won)
		var converted bool
		if won {
			conversionRate := segment.conversionProb * segment.conversionMultiplier
			converted = rand.Float64() < conversionRate
		}

		// Generate keywords
		keywords := g.generateKeywords()

		bidEvent := &models.BidEvent{
			CampaignID:            campaignID,
			UserID:                g.generateUserID(),
			BidPrice:              bidPrice,
			WinPrice:              winPrice,
			FloorPrice:            floorPrice,
			Won:                   won,
			Converted:             converted,
			SegmentID:             segment.id,
			SegmentCategory:       segment.category,
			EngagementScore:       &segment.engagementScore,
			ConversionProbability: &segment.conversionProb,
			Country:               country,
			Region:                region,
			City:                  g.generateCity(region),
			DeviceType:            deviceType,
			OS:                    os,
			Browser:               browser,
			IsMobile:              deviceType == "mobile",
			Keywords:              pq.StringArray(keywords),
			Timestamp:             timestamp,
		}

		if err := g.bidStore.StoreBidEvent(bidEvent); err != nil {
			return fmt.Errorf("failed to store bid event: %w", err)
		}

		// Progress indicator
		if (i+1)%100 == 0 {
			log.Printf("Generated %d/%d bid events", i+1, count)
		}
	}

	return nil
}

func (g *DataGenerator) generateUserID() string {
	return fmt.Sprintf("user_%d", rand.Intn(100000))
}

func (g *DataGenerator) generateKeywords() []string {
	keywordSets := [][]string{
		{"advertising", "marketing", "digital"},
		{"mobile", "app", "download"},
		{"ecommerce", "shopping", "sale"},
		{"software", "saas", "business"},
		{"gaming", "entertainment", "fun"},
		{"fashion", "style", "clothing"},
		{"travel", "vacation", "hotel"},
		{"food", "restaurant", "delivery"},
	}

	set := keywordSets[rand.Intn(len(keywordSets))]
	numKeywords := 1 + rand.Intn(3) // 1-3 keywords

	var keywords []string
	for i := 0; i < numKeywords && i < len(set); i++ {
		keywords = append(keywords, set[i])
	}

	return keywords
}

func (g *DataGenerator) generateCity(region string) string {
	cities := map[string][]string{
		"California": {"San Francisco", "Los Angeles", "San Diego", "Sacramento"},
		"New York":   {"New York City", "Buffalo", "Rochester", "Syracuse"},
		"Texas":      {"Houston", "Dallas", "Austin", "San Antonio"},
		"Ontario":    {"Toronto", "Ottawa", "Hamilton", "London"},
		"England":    {"London", "Manchester", "Birmingham", "Liverpool"},
		"Bavaria":    {"Munich", "Nuremberg", "Augsburg", "Regensburg"},
		"Tokyo":      {"Tokyo", "Shibuya", "Shinjuku", "Harajuku"},
	}

	if cityList, exists := cities[region]; exists {
		return cityList[rand.Intn(len(cityList))]
	}

	// Default cities
	defaultCities := []string{"City Center", "Downtown", "Suburb", "Metro Area"}
	return defaultCities[rand.Intn(len(defaultCities))]
}
