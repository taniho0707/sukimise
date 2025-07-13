package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"sukimise-discord-bot/internal/models"
)

// IsValidGoogleMapsURL checks if the URL is a valid Google Maps place URL
func IsValidGoogleMapsURL(mapURL string) bool {
	return strings.HasPrefix(mapURL, "https://www.google.com/maps/place/") ||
		strings.HasPrefix(mapURL, "https://maps.google.com/maps/place/") ||
		strings.HasPrefix(mapURL, "https://goo.gl/maps/") ||
		strings.HasPrefix(mapURL, "https://maps.app.goo.gl/")
}

// ExtractStoreInfoFromURL extracts store information from Google Maps URL
func ExtractStoreInfoFromURL(mapURL string) (*models.StoreCreateRequest, error) {
	// If it's a shortened URL, expand it first
	if strings.HasPrefix(mapURL, "https://maps.app.goo.gl/") || strings.HasPrefix(mapURL, "https://goo.gl/maps/") {
		expandedURL, err := expandShortenedURL(mapURL)
		if err != nil {
			return nil, fmt.Errorf("failed to expand shortened URL: %v", err)
		}
		mapURL = expandedURL
		fmt.Printf("DEBUG: Expanded URL: %s\n", mapURL)
	}

	// Parse URL to extract information
	parsedURL, err := url.Parse(mapURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// Extract store name from URL path
	storeName := extractStoreNameFromPath(parsedURL.Path)
	fmt.Printf("DEBUG: Store name from URL path: '%s'\n", storeName)
	if storeName == "" {
		storeName = "Unknown Store"
	}

	// Try to get place details from Google Maps Places API first
	address, websiteURL, pageStoreName, apiLat, apiLng, businessHours, err := getPlaceDetailsFromAPI(mapURL)
	var lat, lng float64
	
	if err != nil {
		fmt.Printf("DEBUG: Google Maps API failed, falling back to HTML scraping and URL coordinates: %v\n", err)
		// Fallback to HTML scraping if API fails
		address, websiteURL, pageStoreName, err = extractAdditionalInfoFromPage(mapURL)
		if err != nil {
			// If we can't extract additional info, continue with basic info
			address = ""
			websiteURL = ""
			pageStoreName = ""
		}
		
		// Extract coordinates from URL as fallback
		lat, lng, err = extractCoordinatesFromURL(mapURL)
		if err != nil {
			return nil, fmt.Errorf("failed to extract coordinates: %v", err)
		}
		
		// Use default business hours if API failed
		businessHours = "営業時間: 11:00-22:00\nラストオーダー: 21:30\n定休日: 年中無休"
	} else {
		// Use coordinates from API if available
		lat, lng = apiLat, apiLng
		fmt.Printf("DEBUG: Using coordinates from Places API: %f, %f\n", lat, lng)
	}
	
	// If address extraction failed, try reverse geocoding as fallback
	if address == "" || address == "Address not available" {
		fmt.Printf("DEBUG: GoogleMaps address extraction failed, trying reverse geocoding as fallback\n")
		reverseAddress := reverseGeocodeFromCoordinates(lat, lng)
		if reverseAddress != "" {
			address = reverseAddress
		} else {
			address = "Address not available"
		}
	}
	
	// Use page store name if URL name is invalid or empty
	fmt.Printf("DEBUG: Page store name: '%s'\n", pageStoreName)
	fmt.Printf("DEBUG: Current store name: '%s'\n", storeName)
	if pageStoreName != "" && (storeName == "" || storeName == "Unknown Store" || !isValidStoreName(storeName)) {
		storeName = pageStoreName
		fmt.Printf("DEBUG: Using page store name: '%s'\n", storeName)
	}
	
	// Final fallback if still invalid
	if storeName == "" || !isValidStoreName(storeName) {
		storeName = "Unknown Store"
		fmt.Printf("DEBUG: Using fallback store name: '%s'\n", storeName)
	}
	
	fmt.Printf("DEBUG: Final store name: '%s'\n", storeName)

	// SNS URLとホームページURLを分類
	var snsUrls []string
	var homePageURL string
	
	if websiteURL != "" {
		if isSNSURL(websiteURL) {
			snsUrls = append(snsUrls, websiteURL)
		} else {
			homePageURL = websiteURL
		}
	}

	// Parse business hours from string to struct
	var businessHoursData models.BusinessHoursData
	if businessHours != "" {
		// If we got business hours from API (JSON string), parse it
		if strings.HasPrefix(businessHours, "{") {
			businessHoursData = parseBusinessHoursJSON(businessHours)
		} else {
			// If we got legacy format, use default
			businessHoursData = getDefaultBusinessHoursData()
		}
	} else {
		// Use default business hours
		businessHoursData = getDefaultBusinessHoursData()
	}

	storeInfo := &models.StoreCreateRequest{
		Name:          storeName,
		Address:       address,
		Latitude:      lat,
		Longitude:     lng,
		Categories:    []string{}, // カテゴリーは登録しない
		BusinessHours: businessHoursData, // Google Maps APIから取得した営業時間
		GoogleMapURL:  mapURL,
		WebsiteURL:    homePageURL,
		SNSUrls:       snsUrls,
		Tags:          []string{"discord"},
	}

	return storeInfo, nil
}

// isSNSURL checks if the URL is a social media URL
func isSNSURL(url string) bool {
	lowerURL := strings.ToLower(url)
	return strings.Contains(lowerURL, "instagram.com") ||
		strings.Contains(lowerURL, "facebook.com") ||
		strings.Contains(lowerURL, "twitter.com") ||
		strings.Contains(lowerURL, "x.com") ||
		strings.Contains(lowerURL, "tiktok.com") ||
		strings.Contains(lowerURL, "youtube.com") ||
		strings.Contains(lowerURL, "linkedin.com")
}

// expandShortenedURL expands shortened Google Maps URLs to full URLs
func expandShortenedURL(shortURL string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects automatically, we want to capture the final URL
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("HEAD", shortURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set user agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Get the Location header which contains the expanded URL
	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("no redirect location found")
	}

	// If it's still a redirect, follow it one more time
	if strings.Contains(location, "redirect") || strings.Contains(location, "url=") {
		// Extract the actual URL from redirect parameters
		if strings.Contains(location, "url=") {
			parts := strings.Split(location, "url=")
			if len(parts) > 1 {
				decodedURL, err := url.QueryUnescape(parts[1])
				if err == nil {
					location = decodedURL
				}
			}
		}
	}

	fmt.Printf("DEBUG: Shortened URL %s expanded to: %s\n", shortURL, location)
	return location, nil
}

// reverseGeocodeFromCoordinates attempts to get address from coordinates using OpenStreetMap Nominatim API
func reverseGeocodeFromCoordinates(lat, lng float64) string {
	// Use OpenStreetMap Nominatim API for reverse geocoding
	nominatimURL := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f&accept-language=ja", lat, lng)
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	req, err := http.NewRequest("GET", nominatimURL, nil)
	if err != nil {
		fmt.Printf("DEBUG: Failed to create nominatim request: %v\n", err)
		return ""
	}
	
	// Set user agent for Nominatim
	req.Header.Set("User-Agent", "SukimiseDiscordBot/1.0")
	
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("DEBUG: Failed to make nominatim request: %v\n", err)
		return ""
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("DEBUG: Nominatim returned status: %d\n", resp.StatusCode)
		return ""
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("DEBUG: Failed to read nominatim response: %v\n", err)
		return ""
	}
	
	// Parse Nominatim response
	var nominatimResp struct {
		DisplayName string `json:"display_name"`
		Address     struct {
			Country     string `json:"country"`
			State       string `json:"state"`
			City        string `json:"city"`
			Town        string `json:"town"`
			Village     string `json:"village"`
			Suburb      string `json:"suburb"`
			Road        string `json:"road"`
			HouseNumber string `json:"house_number"`
			Postcode    string `json:"postcode"`
		} `json:"address"`
	}
	
	if err := json.Unmarshal(body, &nominatimResp); err != nil {
		fmt.Printf("DEBUG: Failed to parse nominatim response: %v\n", err)
		return ""
	}
	
	// Use display_name if available
	if nominatimResp.DisplayName != "" {
		fmt.Printf("DEBUG: Nominatim address: %s\n", nominatimResp.DisplayName)
		return nominatimResp.DisplayName
	}
	
	return ""
}

// cleanStoreName removes invalid characters and patterns from store name
func cleanStoreName(name string) string {
	// Remove any JSON-like patterns
	if strings.Contains(name, "[[") || strings.Contains(name, "]]") || strings.Contains(name, "null") {
		return ""
	}
	
	// Remove any coordinate-like patterns
	coordPattern := regexp.MustCompile(`@-?\d+\.\d+,-?\d+\.\d+`)
	name = coordPattern.ReplaceAllString(name, "")
	
	// Remove any data patterns
	dataPattern := regexp.MustCompile(`data-[^=]*="[^"]*"`)
	name = dataPattern.ReplaceAllString(name, "")
	
	// Remove HTML tags
	htmlPattern := regexp.MustCompile(`<[^>]*>`)
	name = htmlPattern.ReplaceAllString(name, "")
	
	// Remove multiple spaces and trim
	name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")
	name = strings.TrimSpace(name)
	
	return name
}

// isValidStoreName checks if the store name is valid
func isValidStoreName(name string) bool {
	if name == "" || name == "Unknown Store" {
		return false
	}
	
	// Check for invalid patterns
	invalidPatterns := []string{
		"[[", "]]", "null", "{", "}", "[", "]",
		"data-", "aria-", "class=", "id=",
	}
	
	lowerName := strings.ToLower(name)
	for _, pattern := range invalidPatterns {
		if strings.Contains(lowerName, pattern) {
			return false
		}
	}
	
	// Check if it's mostly numbers or special characters
	validChars := 0
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '\u3040' && r <= '\u309F') || (r >= '\u30A0' && r <= '\u30FF') || (r >= '\u4E00' && r <= '\u9FAF') {
			validChars++
		}
	}
	
	// Name should have at least 30% valid characters
	return float64(validChars)/float64(len(name)) >= 0.3
}

// extractStoreNameFromHTML extracts store name from HTML content
func extractStoreNameFromHTML(html string) string {
	// Look for store name patterns in the HTML
	patterns := []string{
		`"name":"([^"]+)"`,                      // JSON-LD format
		`<title>([^<]+)</title>`,                // Page title
		`aria-label="([^"]+)"[^>]*role="heading"`, // Heading aria-label
		`<h1[^>]*>([^<]+)</h1>`,                 // H1 heading
		`data-value="([^"]+)"[^>]*aria-label="Name"`, // Name data attribute
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) >= 2 {
			name := strings.TrimSpace(matches[1])
			name = cleanStoreName(name)
			if name != "" && isValidStoreName(name) {
				// Clean up common title suffixes
				name = strings.TrimSuffix(name, " - Google Maps")
				name = strings.TrimSuffix(name, " - Google マップ")
				return strings.TrimSpace(name)
			}
		}
	}

	return ""
}

// isValidAddress checks if the address is valid
func isValidAddress(address string) bool {
	if address == "" || len(address) < 3 {
		return false
	}
	
	// Check for invalid patterns
	invalidPatterns := []string{
		"[[", "]]", "null", "{", "}", "[null",
		"google.com", "maps.google", "data-",
		"aria-", "class=", "id=", "onclick=",
		"javascript:", "function(", "var ",
		"undefined", "NaN", "Infinity",
	}
	
	lowerAddress := strings.ToLower(address)
	for _, pattern := range invalidPatterns {
		if strings.Contains(lowerAddress, pattern) {
			fmt.Printf("DEBUG: Address rejected due to invalid pattern '%s': %s\n", pattern, address)
			return false
		}
	}
	
	// Check for minimum length and meaningful content
	if len(strings.TrimSpace(address)) < 5 {
		fmt.Printf("DEBUG: Address too short: %s\n", address)
		return false
	}
	
	// Check if it contains recognizable address components
	addressIndicators := []string{
		"都", "道", "府", "県", "市", "区", "町", "村",
		"丁目", "番地", "号", "番", "条", "丁", "目",
		"街", "通", "路", "駅", "店", "ビル", "マンション",
		"Japan", "Prefecture", "City", "District",
	}
	
	hasIndicator := false
	for _, indicator := range addressIndicators {
		if strings.Contains(address, indicator) {
			hasIndicator = true
			break
		}
	}
	
	if !hasIndicator {
		fmt.Printf("DEBUG: Address lacks location indicators: %s\n", address)
		return false
	}
	
	// Check character composition
	validChars := 0
	for _, r := range address {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || 
		   (r >= '0' && r <= '9') || (r >= '\u3040' && r <= '\u309F') || 
		   (r >= '\u30A0' && r <= '\u30FF') || (r >= '\u4E00' && r <= '\u9FAF') ||
		   r == ' ' || r == '、' || r == '。' || r == '-' || r == '〒' || r == '／' || r == '−' {
			validChars++
		}
	}
	
	// Address should have at least 60% valid characters
	isValid := float64(validChars)/float64(len(address)) >= 0.6
	if !isValid {
		fmt.Printf("DEBUG: Address failed character composition test: %s (valid: %d/%d)\n", address, validChars, len(address))
	}
	
	return isValid
}

// extractStoreNameFromPath extracts store name from URL path
func extractStoreNameFromPath(path string) string {
	// Google Maps URLs format: /maps/place/Store+Name/...
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "place" && i+1 < len(parts) {
			// Decode URL-encoded store name
			storeName, err := url.QueryUnescape(parts[i+1])
			if err == nil {
				// Replace + with spaces and clean up
				storeName = strings.ReplaceAll(storeName, "+", " ")
				// Remove any @coordinates suffix
				if atIndex := strings.Index(storeName, "@"); atIndex != -1 {
					storeName = storeName[:atIndex]
				}
				// Clean up any invalid characters or patterns
				storeName = cleanStoreName(storeName)
				return strings.TrimSpace(storeName)
			}
		}
	}
	return ""
}

// extractCoordinatesFromURL extracts latitude and longitude from Google Maps URL
func extractCoordinatesFromURL(mapURL string) (float64, float64, error) {
	// Look for coordinates in various formats
	patterns := []string{
		`@(-?\d+\.\d+),(-?\d+\.\d+)`,           // @lat,lng
		`!3d(-?\d+\.\d+)!4d(-?\d+\.\d+)`,       // !3dlat!4dlng
		`ll=(-?\d+\.\d+),(-?\d+\.\d+)`,         // ll=lat,lng
		`center=(-?\d+\.\d+),(-?\d+\.\d+)`,     // center=lat,lng
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(mapURL)
		if len(matches) >= 3 {
			lat, err1 := strconv.ParseFloat(matches[1], 64)
			lng, err2 := strconv.ParseFloat(matches[2], 64)
			if err1 == nil && err2 == nil {
				return lat, lng, nil
			}
		}
	}

	return 0, 0, fmt.Errorf("coordinates not found in URL")
}

// extractAdditionalInfoFromPage attempts to extract additional information from the Google Maps page
func extractAdditionalInfoFromPage(mapURL string) (address, websiteURL, storeName string, err error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make request to Google Maps URL
	req, err := http.NewRequest("GET", mapURL, nil)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set user agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read response body: %v", err)
	}

	pageContent := string(body)
	
	// Debug: Print samples of HTML content to understand structure
	fmt.Printf("DEBUG: HTML length: %d characters\n", len(pageContent))
	
	// Look for postal code (〒) patterns in HTML
	keyword := "〒"
	keywordIndex := strings.Index(pageContent, keyword)
	if keywordIndex != -1 {
		start := keywordIndex - 100
		if start < 0 {
			start = 0
		}
		end := keywordIndex + 200
		if end > len(pageContent) {
			end = len(pageContent)
		}
		fmt.Printf("DEBUG: HTML sample around '%s': %s\n", keyword, pageContent[start:end])
	} else {
		fmt.Printf("DEBUG: No postal code '〒' found in HTML\n")
	}

	// Extract address using regex patterns
	address = extractAddressFromHTML(pageContent)
	fmt.Printf("DEBUG: Extracted address: '%s'\n", address)
	if address == "" {
		address = "Address not available"
	}

	// Extract website URL using regex patterns
	websiteURL = extractWebsiteFromHTML(pageContent)
	
	// Extract store name from HTML
	storeName = extractStoreNameFromHTML(pageContent)

	return address, websiteURL, storeName, nil
}

// extractAddressFromHTML extracts address from HTML content
func extractAddressFromHTML(html string) string {
	// Look for addresses starting with 〒 (postal code) - highest priority
	postalPatterns := []string{
		// Addresses starting with 〒 in various HTML contexts
		`"([^"]*〒[^"]*)"`,                      // JSON format with postal code
		`>([^<]*〒[^<]*)</[^>]*>`,               // Any HTML tag containing postal address
		`<div[^>]*>([^<]*〒[^<]*)</div>`,        // Div containing postal address
		`<span[^>]*>([^<]*〒[^<]*)</span>`,      // Span containing postal address
		`<button[^>]*>.*?([^<]*〒[^<]*).*?</button>`, // Button containing postal address
		`data-value="([^"]*〒[^"]*)"`,           // Data attribute with postal address
		`aria-label="[^"]*([^"]*〒[^"]*)"`,      // Aria label with postal address
		`content="([^"]*〒[^"]*)"`,              // Meta content with postal address
		
		// More flexible patterns for postal addresses
		`([^>\s"]*〒\s*\d{3}-?\d{4}[^<"]*[都道府県][^<"]*[市区町村][^<"]*)`, // Complete postal address
		`([^>\s"]*〒\s*\d{3}-?\d{4}[^<"]*)`,     // Any postal code pattern
	}
	
	fmt.Printf("DEBUG: Searching for addresses starting with 〒\n")
	
	for i, pattern := range postalPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(html, -1) // Find all matches
		
		for _, match := range matches {
			if len(match) >= 2 {
				address := strings.TrimSpace(match[1])
				fmt.Printf("DEBUG: Postal pattern %d found potential address: '%s'\n", i, address)
				
				// Check if address starts with 〒
				if strings.HasPrefix(address, "〒") && isValidPostalAddress(address) {
					fmt.Printf("DEBUG: Valid postal address found: '%s'\n", address)
					return address
				}
			}
		}
	}
	
	fmt.Printf("DEBUG: No addresses starting with 〒 found\n")
	return ""
}

// isValidPostalAddress checks if the postal address starting with 〒 is valid
func isValidPostalAddress(address string) bool {
	if address == "" || !strings.HasPrefix(address, "〒") {
		return false
	}
	
	// Check for minimum length
	if len(address) < 8 { // 〒 + 3 digits + - + 4 digits minimum
		fmt.Printf("DEBUG: Postal address too short: %s\n", address)
		return false
	}
	
	// Check for invalid patterns
	invalidPatterns := []string{
		"[[", "]]", "null", "{", "}", "[null",
		"google.com", "maps.google", "data-",
		"aria-", "class=", "id=", "onclick=",
		"javascript:", "function(", "var ",
		"undefined", "NaN", "Infinity",
	}
	
	lowerAddress := strings.ToLower(address)
	for _, pattern := range invalidPatterns {
		if strings.Contains(lowerAddress, pattern) {
			fmt.Printf("DEBUG: Postal address rejected due to invalid pattern '%s': %s\n", pattern, address)
			return false
		}
	}
	
	// Check for postal code pattern: 〒XXX-XXXX or 〒XXXXXXX
	postalCodePattern := `〒\s*\d{3}-?\d{4}`
	matched, _ := regexp.MatchString(postalCodePattern, address)
	if !matched {
		fmt.Printf("DEBUG: Invalid postal code format: %s\n", address)
		return false
	}
	
	// Check if it contains meaningful address components
	addressIndicators := []string{
		"都", "道", "府", "県", "市", "区", "町", "村",
		"丁目", "番地", "号", "番", "条", "丁", "目",
		"街", "通", "路", "駅", "店", "ビル", "マンション",
	}
	
	hasIndicator := false
	for _, indicator := range addressIndicators {
		if strings.Contains(address, indicator) {
			hasIndicator = true
			break
		}
	}
	
	if !hasIndicator {
		fmt.Printf("DEBUG: Postal address lacks location indicators: %s\n", address)
		// For postal addresses, we're more lenient as they might be abbreviated
		return len(address) >= 10 // Allow if it's reasonably long
	}
	
	fmt.Printf("DEBUG: Postal address validation passed: %s\n", address)
	return true
}

// Google Maps Places API (New) response structures
type PlaceDetailNew struct {
	ID               string                  `json:"id"`
	DisplayName      *PlaceDisplayName       `json:"displayName,omitempty"`
	FormattedAddress string                  `json:"formattedAddress,omitempty"`
	Location         *PlaceLocationNew       `json:"location,omitempty"`
	WebsiteUri       string                  `json:"websiteUri,omitempty"`
	PhoneNumber      string                  `json:"nationalPhoneNumber,omitempty"`
	OpeningHours     *PlaceOpeningHoursNew   `json:"regularOpeningHours,omitempty"`
	PriceLevel       string                  `json:"priceLevel,omitempty"`
	Rating           float64                 `json:"rating,omitempty"`
	Types            []string                `json:"types,omitempty"`
	BusinessStatus   string                  `json:"businessStatus,omitempty"`
}

type PlaceDisplayName struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
}

type PlaceLocationNew struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PlaceOpeningHoursNew struct {
	OpenNow     bool     `json:"openNow"`
	WeekdayText []string `json:"weekdayDescriptions"`
}

// getPlaceDetailsFromAPI gets place details using Google Maps Places API
func getPlaceDetailsFromAPI(mapURL string) (string, string, string, float64, float64, string, error) {
	
	// Extract Place ID from URL
	placeID, err := extractPlaceIDFromURL(mapURL)
	if err != nil {
		return "", "", "", 0, 0, "", fmt.Errorf("failed to extract place ID: %v", err)
	}
	
	// If the extracted ID is in CID format (0x...), we MUST convert it to standard Place ID
	// because Places API (New) doesn't support CID format directly
	if strings.HasPrefix(placeID, "0x") {
		fmt.Printf("DEBUG: Detected CID format: %s, MUST convert to standard Place ID via Text Search\n", placeID)
		standardPlaceID, err := findPlaceIDByTextSearch(mapURL, placeID)
		if err != nil {
			fmt.Printf("DEBUG: Text Search failed: %v, cannot proceed with CID format\n", err)
			return "", "", "", 0, 0, "", fmt.Errorf("CID format not supported by new API, and Text Search failed: %v", err)
		} else {
			placeID = standardPlaceID
			fmt.Printf("DEBUG: Successfully converted CID to standard Place ID: %s\n", placeID)
		}
	}

	// Get API key from environment
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	fmt.Printf("DEBUG: Retrieved API key from environment: '%s'\n", apiKey)
	if apiKey == "" || apiKey == "your_google_maps_api_key_here" {
		fmt.Printf("DEBUG: API key is empty or default value. Available env vars:\n")
		for _, env := range os.Environ() {
			if strings.Contains(env, "GOOGLE") || strings.Contains(env, "API") {
				fmt.Printf("DEBUG: %s\n", env)
			}
		}
		return "", "", "", 0, 0, "", fmt.Errorf("google Maps API key not configured")
	}

	// Build Places API (New) URL
	placesURL := fmt.Sprintf("https://places.googleapis.com/v1/places/%s?languageCode=ja", placeID)
	fmt.Printf("DEBUG: Places API (New) URL: %s\n", placesURL)

	// Make API request with POST method for Places API (New)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", placesURL, nil)
	if err != nil {
		return "", "", "", 0, 0, "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers for Places API (New)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", apiKey)
	req.Header.Set("X-Goog-FieldMask", "id,displayName,formattedAddress,location,websiteUri,nationalPhoneNumber,regularOpeningHours,priceLevel,rating,types,businessStatus")
	
	fmt.Printf("DEBUG: Making GET request to %s with Place ID: %s\n", placesURL, placeID)
	fmt.Printf("DEBUG: API Key prefix: %s...\n", apiKey[:10])

	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", 0, 0, "", fmt.Errorf("failed to make Places API request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read error response body for debugging
		errorBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("DEBUG: Places API error response (status %d): %s\n", resp.StatusCode, string(errorBody))
		return "", "", "", 0, 0, "", fmt.Errorf("places API returned status: %d", resp.StatusCode)
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", 0, 0, "", fmt.Errorf("failed to read Places API response: %v", err)
	}

	fmt.Printf("DEBUG: Raw Places API (New) response: %s\n", string(body))

	var place PlaceDetailNew
	if err := json.Unmarshal(body, &place); err != nil {
		return "", "", "", 0, 0, "", fmt.Errorf("failed to parse Places API response: %v", err)
	}

	// Extract information from the new API response format
	var storeName, address, websiteURL, businessHours string
	var lat, lng float64
	
	// Safely extract display name
	if place.DisplayName != nil && place.DisplayName.Text != "" {
		storeName = place.DisplayName.Text
	}
	
	// Safely extract address
	if place.FormattedAddress != "" {
		address = place.FormattedAddress
	}
	
	// Safely extract website
	if place.WebsiteUri != "" {
		websiteURL = place.WebsiteUri
	}
	
	// Safely extract coordinates
	if place.Location != nil {
		lat = place.Location.Latitude
		lng = place.Location.Longitude
	}
	
	// Safely extract business hours
	if place.OpeningHours != nil && len(place.OpeningHours.WeekdayText) > 0 {
		fmt.Printf("DEBUG: Found %d weekday descriptions\n", len(place.OpeningHours.WeekdayText))
		for i, desc := range place.OpeningHours.WeekdayText {
			fmt.Printf("DEBUG: WeekdayText[%d]: %s\n", i, desc)
		}
		businessHoursData := parseBusinessHoursFromWeekdayDescriptions(place.OpeningHours.WeekdayText)
		// Convert to JSON string for temporary compatibility
		businessHoursJSON, _ := json.Marshal(businessHoursData)
		businessHours = string(businessHoursJSON)
		fmt.Printf("DEBUG: Parsed business hours JSON: %s\n", businessHours)
	} else {
		businessHours = "営業時間: 11:00-22:00\nラストオーダー: 21:30\n定休日: 年中無休"
		fmt.Printf("DEBUG: No business hours found, using default legacy format\n")
	}

	fmt.Printf("DEBUG: Places API (New) returned - Name: '%s', Address: '%s', Website: '%s', Location: (%f, %f), Business Hours: %s\n", 
		storeName, address, websiteURL, lat, lng, businessHours)

	// Check if we got meaningful data
	if storeName == "" && address == "" {
		return "", "", "", 0, 0, "", fmt.Errorf("API returned empty data for place ID: %s", placeID)
	}

	fmt.Printf("DEBUG: Places API SUCCESS - returning data\n")
	return address, websiteURL, storeName, lat, lng, businessHours, nil
}

// parseBusinessHoursFromWeekdayDescriptions parses Google Maps weekdayDescriptions into detailed business hours struct
func parseBusinessHoursFromWeekdayDescriptions(weekdayDescriptions []string) models.BusinessHoursData {
	if len(weekdayDescriptions) == 0 {
		return getDefaultBusinessHoursData()
	}

	// Day mapping from English and Japanese to internal keys
	dayMapping := map[string]string{
		"Monday":    "monday",
		"Tuesday":   "tuesday", 
		"Wednesday": "wednesday",
		"Thursday":  "thursday",
		"Friday":    "friday",
		"Saturday":  "saturday",
		"Sunday":    "sunday",
		"月曜日":      "monday",
		"火曜日":      "tuesday",
		"水曜日":      "wednesday",
		"木曜日":      "thursday",
		"金曜日":      "friday",
		"土曜日":      "saturday",
		"日曜日":      "sunday",
	}

	// Initialize business hours structure
	businessHours := models.BusinessHoursData{
		Monday:    models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{}},
		Tuesday:   models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{}},
		Wednesday: models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{}},
		Thursday:  models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{}},
		Friday:    models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{}},
		Saturday:  models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{}},
		Sunday:    models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{}},
	}
	
	for _, desc := range weekdayDescriptions {
		fmt.Printf("DEBUG: Processing weekday description: %s\n", desc)
		
		// Extract day name
		var dayKey string
		for dayName, key := range dayMapping {
			if strings.HasPrefix(desc, dayName) || strings.Contains(desc, dayName) {
				dayKey = key
				fmt.Printf("DEBUG: Found day '%s' -> '%s' in description: %s\n", dayName, key, desc)
				break
			}
		}
		
		if dayKey == "" {
			fmt.Printf("DEBUG: No day key found for description: %s\n", desc)
			continue
		}
		
		// Get day schedule reference
		daySchedule := getDayScheduleRef(&businessHours, dayKey)
		if daySchedule == nil {
			continue
		}
		
		// Check if closed (English or Japanese)
		if strings.Contains(desc, "Closed") || strings.Contains(desc, "定休日") {
			daySchedule.IsClosed = true
			daySchedule.TimeSlots = []models.TimeSlot{}
			fmt.Printf("DEBUG: Day %s is closed (found 'Closed' or '定休日')\n", dayKey)
			continue
		}
		
		// Check if 24 hours
		if strings.Contains(desc, "Open 24 hours") {
			daySchedule.IsClosed = false
			daySchedule.TimeSlots = []models.TimeSlot{{
				OpenTime:      "00:00",
				CloseTime:     "00:00",
				LastOrderTime: "00:00",
			}}
			continue
		}
		
		// Parse time ranges (handle multiple time ranges per day)
		timeRanges := parseDetailedTimeRanges(desc)
		if len(timeRanges) > 0 {
			daySchedule.IsClosed = false
			daySchedule.TimeSlots = timeRanges
		}
	}
	
	return businessHours
}

// getDayScheduleRef returns a reference to the day schedule
func getDayScheduleRef(businessHours *models.BusinessHoursData, dayKey string) *models.DaySchedule {
	switch dayKey {
	case "monday":
		return &businessHours.Monday
	case "tuesday":
		return &businessHours.Tuesday
	case "wednesday":
		return &businessHours.Wednesday
	case "thursday":
		return &businessHours.Thursday
	case "friday":
		return &businessHours.Friday
	case "saturday":
		return &businessHours.Saturday
	case "sunday":
		return &businessHours.Sunday
	default:
		return nil
	}
}

// getDefaultBusinessHoursData returns default business hours struct
func getDefaultBusinessHoursData() models.BusinessHoursData {
	defaultSlot := models.TimeSlot{
		OpenTime:      "11:00",
		CloseTime:     "22:00",
		LastOrderTime: "21:30",
	}
	
	return models.BusinessHoursData{
		Monday:    models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{defaultSlot}},
		Tuesday:   models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{defaultSlot}},
		Wednesday: models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{defaultSlot}},
		Thursday:  models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{defaultSlot}},
		Friday:    models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{defaultSlot}},
		Saturday:  models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{defaultSlot}},
		Sunday:    models.DaySchedule{IsClosed: false, TimeSlots: []models.TimeSlot{defaultSlot}},
	}
}

// parseDetailedTimeRanges extracts detailed time ranges from description text
func parseDetailedTimeRanges(desc string) []models.TimeSlot {
	var timeSlots []models.TimeSlot
	
	fmt.Printf("DEBUG: Parsing time ranges from description: %s\n", desc)
	
	// Handle Japanese format: "月曜日: 11時00分～20時00分" or "Tuesday: 11時00分～20時00分"
	// Also handle English format: "6:00 AM – 10:00 PM"
	// Split by comma first to handle multiple ranges
	rangeParts := strings.Split(desc, ",")
	
	for _, part := range rangeParts {
		part = strings.TrimSpace(part)
		fmt.Printf("DEBUG: Processing time range part: %s\n", part)
		
		var openTime, closeTime string
		
		// Pattern 1: Japanese format with 時分 - "11時00分～20時00分"
		jaTimePattern := regexp.MustCompile(`(\d{1,2})時(\d{2})分[～〜~\-–]+(\d{1,2})時(\d{2})分`)
		jaMatches := jaTimePattern.FindStringSubmatch(part)
		
		if len(jaMatches) >= 5 {
			startHour, _ := strconv.Atoi(jaMatches[1])
			startMin, _ := strconv.Atoi(jaMatches[2])
			endHour, _ := strconv.Atoi(jaMatches[3])
			endMin, _ := strconv.Atoi(jaMatches[4])
			
			openTime = fmt.Sprintf("%02d:%02d", startHour, startMin)
			closeTime = fmt.Sprintf("%02d:%02d", endHour, endMin)
			
			fmt.Printf("DEBUG: Japanese time pattern matched - Open: %s, Close: %s\n", openTime, closeTime)
		} else {
			// Pattern 2: English AM/PM format - "6:00 AM – 10:00 PM"
			enTimePattern := regexp.MustCompile(`(\d{1,2}):(\d{2})\s*(AM|PM)?\s*[–\-～〜~]+\s*(\d{1,2}):(\d{2})\s*(AM|PM)?`)
			enMatches := enTimePattern.FindStringSubmatch(part)
			
			if len(enMatches) >= 7 {
				startHour, _ := strconv.Atoi(enMatches[1])
				startMin, _ := strconv.Atoi(enMatches[2])
				startPeriod := enMatches[3]
				endHour, _ := strconv.Atoi(enMatches[4])
				endMin, _ := strconv.Atoi(enMatches[5])
				endPeriod := enMatches[6]
				
				// Convert to 24-hour format
				openTime = convertTo24Hour(startHour, startMin, startPeriod)
				closeTime = convertTo24Hour(endHour, endMin, endPeriod)
				
				fmt.Printf("DEBUG: English time pattern matched - Open: %s, Close: %s\n", openTime, closeTime)
			} else {
				fmt.Printf("DEBUG: No time pattern matched for: %s\n", part)
				continue
			}
		}
		
		if openTime != "" && closeTime != "" {
			// Calculate last order time (30 minutes before close)
			lastOrderTime := calculateLastOrderTime(closeTime)
			
			timeSlots = append(timeSlots, models.TimeSlot{
				OpenTime:      openTime,
				CloseTime:     closeTime,
				LastOrderTime: lastOrderTime,
			})
			
			fmt.Printf("DEBUG: Added time slot - Open: %s, Close: %s, LastOrder: %s\n", 
				openTime, closeTime, lastOrderTime)
		}
	}
	
	// Limit to maximum 3 time slots
	if len(timeSlots) > 3 {
		timeSlots = timeSlots[:3]
	}
	
	fmt.Printf("DEBUG: Total time slots parsed: %d\n", len(timeSlots))
	return timeSlots
}


// convertTo24Hour converts 12-hour format to 24-hour format
func convertTo24Hour(hour, minute int, period string) string {
	if period == "PM" && hour != 12 {
		hour += 12
	} else if period == "AM" && hour == 12 {
		hour = 0
	}
	
	return fmt.Sprintf("%02d:%02d", hour, minute)
}

// calculateLastOrderTime calculates last order time (30 minutes before close)
func calculateLastOrderTime(closeTime string) string {
	// Parse close time
	parts := strings.Split(closeTime, ":")
	if len(parts) != 2 {
		return "21:30"
	}
	
	hour, err1 := strconv.Atoi(parts[0])
	minute, err2 := strconv.Atoi(parts[1])
	
	if err1 != nil || err2 != nil {
		return "21:30"
	}
	
	// Handle 24-hour operation
	if hour == 0 && minute == 0 {
		return "00:00"
	}
	
	// Subtract 30 minutes
	minute -= 30
	if minute < 0 {
		minute += 60
		hour -= 1
		if hour < 0 {
			hour = 23
		}
	}
	
	return fmt.Sprintf("%02d:%02d", hour, minute)
}

// parseBusinessHoursJSON parses business hours JSON string to struct
func parseBusinessHoursJSON(jsonStr string) models.BusinessHoursData {
	var businessHours models.BusinessHoursData
	if err := json.Unmarshal([]byte(jsonStr), &businessHours); err != nil {
		fmt.Printf("DEBUG: Failed to parse business hours JSON: %v\n", err)
		return getDefaultBusinessHoursData()
	}
	return businessHours
}

// TextSearchResponse represents the response from Places API Text Search
type TextSearchResponse struct {
	Places []PlaceBasic `json:"places"`
}

type PlaceBasic struct {
	ID          string           `json:"id"`
	DisplayName PlaceDisplayName `json:"displayName"`
	Location    PlaceLocationNew `json:"location"`
}

// findPlaceIDByTextSearch attempts to find a standard Place ID using Text Search API
func findPlaceIDByTextSearch(mapURL, cid string) (string, error) {
	// Extract place name and coordinates from URL for search
	placeName := extractPlaceNameFromURL(mapURL)
	lat, lng, err := extractCoordinatesFromURL(mapURL)
	if err != nil || placeName == "" {
		return "", fmt.Errorf("insufficient data for text search")
	}

	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API key not available")
	}

	// Prepare Text Search request
	searchURL := "https://places.googleapis.com/v1/places:searchText"
	
	requestBody := map[string]interface{}{
		"textQuery": placeName,
		"languageCode": "ja",
		"locationBias": map[string]interface{}{
			"circle": map[string]interface{}{
				"center": map[string]float64{
					"latitude":  lat,
					"longitude": lng,
				},
				"radius": 1000.0, // 1km radius
			},
		},
		"maxResultCount": 5,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal search request: %v", err)
	}

	req, err := http.NewRequest("POST", searchURL, strings.NewReader(string(requestJSON)))
	if err != nil {
		return "", fmt.Errorf("failed to create search request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", apiKey)
	req.Header.Set("X-Goog-FieldMask", "places.id,places.displayName,places.location")
	req.Header.Set("Accept-Language", "ja")
	req.Header.Set("X-Goog-LanguageCode", "ja")
	
	fmt.Printf("DEBUG: Text Search request - placeName: '%s', coordinates: (%f, %f)\n", placeName, lat, lng)
	fmt.Printf("DEBUG: Text Search request body: %s\n", string(requestJSON))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make search request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("DEBUG: Text Search API error response (status %d): %s\n", resp.StatusCode, string(errorBody))
		return "", fmt.Errorf("search API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read search response: %v", err)
	}

	var searchResp TextSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return "", fmt.Errorf("failed to parse search response: %v", err)
	}

	if len(searchResp.Places) == 0 {
		return "", fmt.Errorf("no places found in search")
	}

	// Return the first result's Place ID
	bestPlace := searchResp.Places[0]
	fmt.Printf("DEBUG: Text Search found Place ID: %s for '%s'\n", bestPlace.ID, bestPlace.DisplayName.Text)
	return bestPlace.ID, nil
}

// extractPlaceIDFromURL extracts Google Place ID from Google Maps URL
func extractPlaceIDFromURL(mapURL string) (string, error) {
	fmt.Printf("DEBUG: Attempting to extract place_id from URL: %s\n", mapURL)
	
	// Method 1: Direct place_id parameter in URL
	if strings.Contains(mapURL, "place_id=") {
		re := regexp.MustCompile(`place_id=([A-Za-z0-9_-]+)`)
		matches := re.FindStringSubmatch(mapURL)
		if len(matches) >= 2 {
			fmt.Printf("DEBUG: Found place_id in URL parameter: %s\n", matches[1])
			return matches[1], nil
		}
	}

	// Method 2: Extract from data parameter - multiple place references
	// This URL contains multiple places, extract all possible place IDs
	placePatterns := []string{
		`!1s(0x[a-f0-9]+:0x[a-f0-9]+)!`,              // Hex format place IDs
		`!3m5!1s(0x[a-f0-9]+:0x[a-f0-9]+)!`,          // Specific to the main place
		`data=[^&]*!3m5!1s(0x[a-f0-9]+:0x[a-f0-9]+)`, // In data parameter
		`/place/[^/]+/@[^/]+/data=[^!]*!3m5!1s(0x[a-f0-9]+:0x[a-f0-9]+)`, // Full path context
	}
	
	// Find all potential place IDs
	var foundPlaceIDs []string
	for i, pattern := range placePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(mapURL, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				placeID := match[1]
				fmt.Printf("DEBUG: Found potential place_id with pattern %d: %s\n", i, placeID)
				foundPlaceIDs = append(foundPlaceIDs, placeID)
			}
		}
	}
	
	// Method 3: For complex URLs, try to identify the main place by context
	// The URL path contains the place name, try to match the most relevant place ID
	if len(foundPlaceIDs) > 0 {
		// Extract place name from URL path for context
		placeName := extractPlaceNameFromURL(mapURL)
		fmt.Printf("DEBUG: Extracted place name from URL: '%s'\n", placeName)
		
		// For now, use the last found place ID (often the main target)
		// In this case: 0x35546d98282af541:0xbe9e600b1755010b
		bestPlaceID := foundPlaceIDs[len(foundPlaceIDs)-1]
		fmt.Printf("DEBUG: Using place_id: %s\n", bestPlaceID)
		return bestPlaceID, nil
	}

	// Method 4: Extract from place URL path and fetch to get place_id from page
	if strings.Contains(mapURL, "/place/") {
		fmt.Printf("DEBUG: Attempting to extract place_id from page content\n")
		placeID, err := extractPlaceIDFromPage(mapURL)
		if err == nil && placeID != "" {
			fmt.Printf("DEBUG: Extracted place_id from page: %s\n", placeID)
			return placeID, nil
		}
		fmt.Printf("DEBUG: Failed to extract from page: %v\n", err)
	}

	return "", fmt.Errorf("could not extract place_id from URL: %s", mapURL)
}

// extractPlaceNameFromURL extracts the place name from Google Maps URL path
func extractPlaceNameFromURL(mapURL string) string {
	// Extract place name from URL path like /place/日本一のだがし売場/@...
	re := regexp.MustCompile(`/place/([^/@]+)`)
	matches := re.FindStringSubmatch(mapURL)
	if len(matches) >= 2 {
		// URL decode the place name
		decoded, err := url.QueryUnescape(matches[1])
		if err == nil {
			return decoded
		}
		return matches[1]
	}
	return ""
}

// extractPlaceIDFromPage extracts place ID by fetching the Google Maps page
func extractPlaceIDFromPage(mapURL string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", mapURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read page: %v", err)
	}

	pageContent := string(body)

	// Look for place_id patterns in the page content
	patterns := []string{
		`"place_id":"([A-Za-z0-9_-]+)"`,     // JSON format
		`place_id=([A-Za-z0-9_-]+)`,            // URL parameter
		`data-place-id="([A-Za-z0-9_-]+)"`,   // Data attribute
		`"placeId":"([A-Za-z0-9_-]+)"`,     // Alternative JSON format
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(pageContent)
		if len(matches) >= 2 {
			placeID := matches[1]
			if len(placeID) > 10 { // Place IDs are typically longer than 10 characters
				return placeID, nil
			}
		}
	}

	return "", fmt.Errorf("place_id not found in page content")
}

// extractWebsiteFromHTML extracts website URL from HTML content
func extractWebsiteFromHTML(html string) string {
	// Look for website patterns in the HTML
	patterns := []string{
		`"url":"(https?://[^"]+)"`,               // JSON-LD format
		`href="(https?://[^"]+)"[^>]*aria-label="Website"`, // Website link
		`data-value="(https?://[^"]+)"[^>]*aria-label="Website"`, // Data attribute
		`aria-label="Website: (https?://[^"]+)"`,  // Website aria-label
		`"sameAs":\["(https?://[^"]+)"\]`,         // JSON-LD sameAs
		`website.*?"(https?://[^"]+)"`,            // General website pattern
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) >= 2 {
			website := strings.TrimSpace(matches[1])
			if website != "" && !strings.Contains(website, "google.com") {
				return website
			}
		}
	}

	return ""
}