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
	validPrefixes := []string{
		"https://www.google.com/maps/place/",
		"https://maps.google.com/maps/place/",
		"https://www.google.co.jp/maps/place/",
		"https://maps.google.co.jp/maps/place/",
		"https://goo.gl/maps/",
		"https://maps.app.goo.gl/",
		"https://www.google.com/maps/@",
		"https://maps.google.com/maps/@",
	}
	
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(mapURL, prefix) {
			return true
		}
	}
	
	// Also check for URLs that contain place information even if they don't start with typical prefixes
	if strings.Contains(mapURL, "google.com/maps") && 
	   (strings.Contains(mapURL, "place/") || strings.Contains(mapURL, "@") || strings.Contains(mapURL, "data=")) {
		return true
	}
	
	return false
}

// ExtractStoreInfoFromURL extracts store information from Google Maps URL using simplified API-based approach
func ExtractStoreInfoFromURL(mapURL string) (*models.StoreCreateRequest, error) {
	fmt.Printf("DEBUG: Starting URL processing: %s\n", mapURL)
	
	// Step 1: Expand shortened URL if needed
	processedURL := mapURL
	if strings.HasPrefix(mapURL, "https://maps.app.goo.gl/") || strings.HasPrefix(mapURL, "https://goo.gl/maps/") {
		expandedURL, err := expandShortenedURL(mapURL)
		if err != nil {
			return nil, fmt.Errorf("failed to expand shortened URL: %v", err)
		}
		processedURL = expandedURL
		fmt.Printf("DEBUG: Expanded URL: %s\n", processedURL)
	}
	
	// Step 2: URL decode if needed
	if decodedURL, err := url.QueryUnescape(processedURL); err == nil {
		processedURL = decodedURL
		fmt.Printf("DEBUG: Decoded URL: %s\n", processedURL)
	}
	
	// Step 3: Extract store name from URL
	storeName := extractStoreNameFromURL(processedURL)
	if storeName == "" {
		return nil, fmt.Errorf("failed to extract store name from URL")
	}
	fmt.Printf("DEBUG: Extracted store name: '%s'\n", storeName)
	
	// Step 4: Try to extract coordinates from URL
	lat, lng, hasCoordinates := tryExtractCoordinatesFromURL(processedURL)
	
	var placeID string
	var err error
	
	if strings.Contains(storeName, "〒") || !hasCoordinates { // 郵便番号が含まれていれば Android、そうでないならブラウザからのリンクと仮定、最適な検索をかける
		// Step 5b: Use Text Search API with store name only
		fmt.Printf("DEBUG: Using Text Search API with store name: %s\n", storeName)
		placeID, err = findPlaceIDByTextSearch(storeName)
	} else {
		// Step 5a: Use Nearby Search API with coordinates and store name
		fmt.Printf("DEBUG: Using Nearby Search API with coordinates: %f, %f, %s\n", lat, lng, storeName)
		placeID, err = findPlaceIDByNearbySearch(storeName, lat, lng)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to find place ID: %v", err)
	}
	
	if placeID == "" {
		return nil, fmt.Errorf("no place found for the given store")
	}
	
	fmt.Printf("DEBUG: Found Place ID: %s\n", placeID)
	
	// Step 6: Get detailed information using Place Details API
	storeInfo, err := getStoreDetailsFromPlaceID(placeID, mapURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get store details: %v", err)
	}
	
	fmt.Printf("DEBUG: Successfully extracted store info - Name: '%s', Address: '%s'\n", 
		storeInfo.Name, storeInfo.Address)
	
	return storeInfo, nil
}

// extractStoreNameFromURL extracts store name from URL using both simple and complex methods
func extractStoreNameFromURL(processedURL string) string {
	// Try simple path extraction first
	if parsedURL, err := url.Parse(processedURL); err == nil {
		storeName := extractStoreNameFromPath(parsedURL.Path)
		if storeName != "" && isValidStoreName(storeName) {
			return storeName
		}
	}
	
	// Try complex extraction from decoded URL
	storeName := extractStoreNameFromDecodedURL(processedURL)
	if storeName != "" && isValidStoreName(storeName) {
		return storeName
	}
	
	return ""
}

// tryExtractCoordinatesFromURL tries to extract coordinates from URL and returns whether coordinates were found
func tryExtractCoordinatesFromURL(processedURL string) (lat, lng float64, hasCoordinates bool) {
	lat, lng, err := extractCoordinatesFromURL(processedURL)
	if err != nil {
		return 0, 0, false
	}
	return lat, lng, true
}

// findPlaceIDByNearbySearch finds Place ID using Nearby Search API with coordinates and store name
func findPlaceIDByNearbySearch(storeName string, lat, lng float64) (string, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("Google Maps API key not found")
	}
	
	// Use Places API (New) Nearby Search
	searchURL := fmt.Sprintf("https://places.googleapis.com/v1/places:searchNearby?key=%s&languageCode=ja", apiKey)
	
	requestBody := map[string]interface{}{
		"maxResultCount": 10,
		"locationRestriction": map[string]interface{}{
			"circle": map[string]interface{}{
				"center": map[string]interface{}{
					"latitude":  lat,
					"longitude": lng,
				},
				"radius": 100.0, // 100m radius
			},
		},
		"languageCode": "ja",
	}
	
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal nearby search request: %v", err)
	}
	
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", searchURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("failed to create nearby search request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-FieldMask", "places.id,places.displayName")
	
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make nearby search request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("nearby search API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var searchResp struct {
		Places []struct {
			ID          string `json:"id"`
			DisplayName struct {
				Text string `json:"text"`
			} `json:"displayName"`
		} `json:"places"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return "", fmt.Errorf("failed to decode nearby search response: %v", err)
	}
	
	if len(searchResp.Places) == 0 {
		return "", fmt.Errorf("no places found nearby")
	}
	
	// Find the best match by name similarity
	bestMatch := ""
	bestScore := 0.0
	
	for _, place := range searchResp.Places {
		score := calculateSimilarity(storeName, place.DisplayName.Text)
		fmt.Printf("DEBUG: Nearby place '%s' similarity score: %f\n", place.DisplayName.Text, score)
		if score > bestScore {
			bestScore = score
			bestMatch = place.ID
		}
	}
	
	if bestScore < 0.3 { // Minimum similarity threshold
		return "", fmt.Errorf("no similar places found nearby (best score: %f)", bestScore)
	}
	
	fmt.Printf("DEBUG: Best nearby match with score %f\n", bestScore)
	return bestMatch, nil
}

// findPlaceIDByTextSearch finds Place ID using Text Search API with store name
func findPlaceIDByTextSearch(storeName string) (string, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("Google Maps API key not found")
	}
	
	// Use Places API (New) Text Search
	searchURL := fmt.Sprintf("https://places.googleapis.com/v1/places:searchText?key=%s&languageCode=ja", apiKey)
	
	requestBody := map[string]interface{}{
		"textQuery": storeName,
		"languageCode": "ja",
		"maxResultCount": 5,
	}
	
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal text search request: %v", err)
	}
	
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", searchURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("failed to create text search request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-FieldMask", "places.id,places.displayName")
	
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make text search request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("text search API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var searchResp struct {
		Places []struct {
			ID          string `json:"id"`
			DisplayName struct {
				Text string `json:"text"`
			} `json:"displayName"`
		} `json:"places"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return "", fmt.Errorf("failed to decode text search response: %v", err)
	}
	
	if len(searchResp.Places) == 0 {
		return "", fmt.Errorf("no places found by text search")
	}
	
	// Return the first result (best match from Google)
	return searchResp.Places[0].ID, nil
}

// getStoreDetailsFromPlaceID gets detailed store information using Place Details API
func getStoreDetailsFromPlaceID(placeID, originalURL string) (*models.StoreCreateRequest, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Google Maps API key not found")
	}
	
	// Use Places API (New) Place Details
	detailsURL := fmt.Sprintf("https://places.googleapis.com/v1/places/%s?key=%s&languageCode=ja", placeID, apiKey)
	
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", detailsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create details request: %v", err)
	}
	req.Header.Set("X-Goog-FieldMask", "id,displayName,formattedAddress,location,websiteUri,regularOpeningHours")
	req.Header.Set("X-Goog-Api-Key", apiKey)
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make details request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("details API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var detailsResp struct {
		ID              string `json:"id"`
		DisplayName     struct {
			Text string `json:"text"`
		} `json:"displayName"`
		FormattedAddress string `json:"formattedAddress"`
		Location        struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
		WebsiteUri      string `json:"websiteUri"`
		RegularOpeningHours struct {
			WeekdayDescriptions []string `json:"weekdayDescriptions"`
		} `json:"regularOpeningHours"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&detailsResp); err != nil {
		return nil, fmt.Errorf("failed to decode details response: %v", err)
	}
	
	// Parse business hours
	var businessHoursData models.BusinessHoursData
	if len(detailsResp.RegularOpeningHours.WeekdayDescriptions) > 0 {
		businessHoursData = parseBusinessHoursFromWeekdayDescriptions(detailsResp.RegularOpeningHours.WeekdayDescriptions)
	} else {
		businessHoursData = getDefaultBusinessHoursData()
	}
	
	// Clean up address (remove Japan prefix if present)
	address := detailsResp.FormattedAddress
	if strings.HasPrefix(address, "日本、") {
		address = strings.TrimPrefix(address, "日本、")
	}
	
	// Classify URL as SNS or website
	var snsUrls []string
	var websiteURL string
	if detailsResp.WebsiteUri != "" {
		if isSNSURL(detailsResp.WebsiteUri) {
			snsUrls = append(snsUrls, detailsResp.WebsiteUri)
		} else {
			websiteURL = detailsResp.WebsiteUri
		}
	}
	
	return &models.StoreCreateRequest{
		Name:          detailsResp.DisplayName.Text,
		Address:       address,
		Latitude:      detailsResp.Location.Latitude,
		Longitude:     detailsResp.Location.Longitude,
		Categories:    []string{}, // Don't set categories
		BusinessHours: businessHoursData,
		GoogleMapURL:  originalURL,
		WebsiteURL:    websiteURL,
		SNSUrls:       snsUrls,
		Tags:          []string{"discord"},
	}, nil
}

// calculateSimilarity calculates similarity between two strings (simple implementation)
func calculateSimilarity(str1, str2 string) float64 {
	// Simple Jaccard similarity using character n-grams
	str1 = strings.ToLower(str1)
	str2 = strings.ToLower(str2)
	
	// Create character sets
	set1 := make(map[rune]bool)
	set2 := make(map[rune]bool)
	
	for _, r := range str1 {
		set1[r] = true
	}
	for _, r := range str2 {
		set2[r] = true
	}
	
	// Calculate intersection and union
	intersection := 0
	union := make(map[rune]bool)
	
	for r := range set1 {
		union[r] = true
		if set2[r] {
			intersection++
		}
	}
	for r := range set2 {
		union[r] = true
	}
	
	if len(union) == 0 {
		return 0.0
	}
	
	return float64(intersection) / float64(len(union))
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
	
	// // Check if it's mostly numbers or special characters
	// validChars := 0
	// for _, r := range name {
	// 	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '\u3040' && r <= '\u309F') || (r >= '\u30A0' && r <= '\u30FF') || (r >= '\u4E00' && r <= '\u9FAF') {
	// 		validChars++
	// 	}
	// }
	
	// // Name should have at least 30% valid characters
	// return float64(validChars)/float64(len(name)) >= 0.3

	return true
}

// // isValidAddress checks if the address is valid
// func isValidAddress(address string) bool {
// 	if address == "" || len(address) < 3 {
// 		return false
// 	}
	
// 	// Check for invalid patterns
// 	invalidPatterns := []string{
// 		"[[", "]]", "null", "{", "}", "[null",
// 		"google.com", "maps.google", "data-",
// 		"aria-", "class=", "id=", "onclick=",
// 		"javascript:", "function(", "var ",
// 		"undefined", "NaN", "Infinity",
// 	}
	
// 	lowerAddress := strings.ToLower(address)
// 	for _, pattern := range invalidPatterns {
// 		if strings.Contains(lowerAddress, pattern) {
// 			fmt.Printf("DEBUG: Address rejected due to invalid pattern '%s': %s\n", pattern, address)
// 			return false
// 		}
// 	}
	
// 	// Check for minimum length and meaningful content
// 	if len(strings.TrimSpace(address)) < 5 {
// 		fmt.Printf("DEBUG: Address too short: %s\n", address)
// 		return false
// 	}
	
// 	// Check if it contains recognizable address components
// 	addressIndicators := []string{
// 		"都", "道", "府", "県", "市", "区", "町", "村",
// 		"丁目", "番地", "号", "番", "条", "丁", "目",
// 		"街", "通", "路", "駅", "店", "ビル", "マンション",
// 		"Japan", "Prefecture", "City", "District",
// 	}
	
// 	hasIndicator := false
// 	for _, indicator := range addressIndicators {
// 		if strings.Contains(address, indicator) {
// 			hasIndicator = true
// 			break
// 		}
// 	}
	
// 	if !hasIndicator {
// 		fmt.Printf("DEBUG: Address lacks location indicators: %s\n", address)
// 		return false
// 	}
	
// 	// Check character composition
// 	validChars := 0
// 	for _, r := range address {
// 		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || 
// 		   (r >= '0' && r <= '9') || (r >= '\u3040' && r <= '\u309F') || 
// 		   (r >= '\u30A0' && r <= '\u30FF') || (r >= '\u4E00' && r <= '\u9FAF') ||
// 		   r == ' ' || r == '、' || r == '。' || r == '-' || r == '〒' || r == '／' || r == '−' {
// 			validChars++
// 		}
// 	}
	
// 	// Address should have at least 60% valid characters
// 	isValid := float64(validChars)/float64(len(address)) >= 0.6
// 	if !isValid {
// 		fmt.Printf("DEBUG: Address failed character composition test: %s (valid: %d/%d)\n", address, validChars, len(address))
// 	}
	
// 	return isValid
// }

// extractStoreNameFromPath extracts store name from URL path
func extractStoreNameFromPath(path string) string {
	// Google Maps URLs format: /maps/place/Store+Name/... or with URL encoding
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "place" && i+1 < len(parts) {
			storeNamePart := parts[i+1]
			
			// First, try to URL decode the store name part
			if decodedPart, err := url.QueryUnescape(storeNamePart); err == nil {
				storeNamePart = decodedPart
			}
			
			// Replace + with spaces (common in URL encoding)
			storeNamePart = strings.ReplaceAll(storeNamePart, "+", " ")
			
			// Remove any @coordinates suffix
			if atIndex := strings.Index(storeNamePart, "@"); atIndex != -1 {
				storeNamePart = storeNamePart[:atIndex]
			}
			
			fmt.Printf("DEBUG: Decoded store name part: '%s'\n", storeNamePart)
			
			// For complex URLs with address + store name, try to extract the actual store name
			storeName := extractStoreNameFromComplexPath(storeNamePart)
			if storeName == "" {
				storeName = storeNamePart
			}
			
			// Clean up any invalid characters or patterns
			storeName = cleanStoreName(storeName)
			cleanedName := strings.TrimSpace(storeName)
			
			// Validate the extracted store name
			if cleanedName != "" && isValidStoreName(cleanedName) {
				fmt.Printf("DEBUG: Extracted store name from path: '%s'\n", cleanedName)
				return cleanedName
			} else {
				fmt.Printf("DEBUG: Invalid store name extracted: '%s'\n", cleanedName)
			}
		}
	}
	return ""
}

// extractStoreNameFromComplexPath extracts store name from complex path that may contain address + store name
func extractStoreNameFromComplexPath(decodedPath string) string {
	// For Japanese addresses, the store name is often at the end after the address
	// Common patterns:
	// 1. "〒123-4567 都道府県市区町村 店舗名"
	// 2. "住所情報 店舗名"
	
	// 郵便番号は後の処理で使うため削除しない
	// // Remove postal code at the beginning if present
	// if strings.HasPrefix(decodedPath, "〒") {
	// 	// Find the end of postal code pattern: 〒xxx-xxxx
	// 	if match := regexp.MustCompile(`^〒\d{3}-\d{4}\s*`).FindString(decodedPath); match != "" {
	// 		decodedPath = strings.TrimSpace(decodedPath[len(match):])
	// 	}
	// }
	// fmt.Printf("DEBUG: After postal code removal: '%s'\n", decodedPath)
	
	// Split by common separators and try to identify the store name
	parts := strings.Fields(decodedPath)
	if len(parts) == 0 {
		return ""
	}
	
	// For Japanese addresses, address components often end with specific kanji
	// Store name is usually the part that doesn't end with these characters
	addressEndings := []string{"都", "道", "府", "県", "市", "区", "町", "村", "丁目", "番地", "号"}
	
	var potentialStoreName string
	
	// Look for parts that don't end with address indicators
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		isAddressPart := false
		
		// Check if this part ends with address indicators
		for _, ending := range addressEndings {
			if strings.HasSuffix(part, ending) {
				isAddressPart = true
				break
			}
		}
		
		// Also check for number-only parts (likely address numbers)
		if regexp.MustCompile(`^[\d０-９]+$`).MatchString(part) {
			isAddressPart = true
		}
		
		if !isAddressPart && len(part) > 1 {
			// This could be the store name
			if potentialStoreName == "" {
				potentialStoreName = part
			} else {
				potentialStoreName = part + " " + potentialStoreName
			}
		} else if potentialStoreName != "" {
			// We found address parts after finding a potential store name, so we can stop
			break
		}
	}
	
	fmt.Printf("DEBUG: Potential store name from complex path: '%s'\n", potentialStoreName)
	
	if potentialStoreName != "" && len(potentialStoreName) > 1 {
		return potentialStoreName
	}
	
	// Fallback: if we can't identify the store name, take the last meaningful part
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	
	return ""
}

// extractStoreNameFromDecodedURL extracts store name from fully decoded URL
func extractStoreNameFromDecodedURL(decodedURL string) string {
	// Look for /place/ in the decoded URL
	if !strings.Contains(decodedURL, "/place/") {
		return ""
	}
	
	// Extract the part after /place/
	parts := strings.Split(decodedURL, "/place/")
	if len(parts) < 2 {
		return ""
	}
	
	// Get the path part after /place/
	pathPart := parts[1]
	
	// Remove any parameters that come after the place name
	if paramIndex := strings.Index(pathPart, "/data="); paramIndex != -1 {
		pathPart = pathPart[:paramIndex]
	}
	if paramIndex := strings.Index(pathPart, "?"); paramIndex != -1 {
		pathPart = pathPart[:paramIndex]
	}
	
	fmt.Printf("DEBUG: Path part for store name extraction: '%s'\n", pathPart)
	
	// Use the complex path extraction logic
	storeName := extractStoreNameFromComplexPath(pathPart)
	if storeName != "" && isValidStoreName(storeName) {
		return storeName
	}
	
	return ""
}

// extractPlaceIDFromURL extracts Place ID from Google Maps URL
func extractPlaceIDFromURL(mapURL string) string {
	fmt.Printf("DEBUG: Extracting Place ID from URL: %s\n", mapURL)
	
	// Look for Place ID in data parameter (newer format)
	if strings.Contains(mapURL, "1s0x") && strings.Contains(mapURL, ":0x") {
		// Format: 1s0x6001078b6e352ab1:0x1869084cc73893fc
		pattern := regexp.MustCompile(`1s(0x[a-fA-F0-9]+:0x[a-fA-F0-9]+)`)
		matches := pattern.FindStringSubmatch(mapURL)
		if len(matches) >= 2 {
			cid := matches[1]
			fmt.Printf("DEBUG: Found CID format: %s\n", cid)
			// Convert CID to Place ID using Places API Text Search
			if placeID, err := convertCIDToPlaceID(cid); err == nil && placeID != "" {
				fmt.Printf("DEBUG: Successfully converted CID to Place ID: %s\n", placeID)
				return placeID
			} else {
				fmt.Printf("DEBUG: Failed to convert CID to Place ID: %v\n", err)
			}
		}
	}
	
	// Look for ftid parameter (Place ID)
	if strings.Contains(mapURL, "ftid=") {
		pattern := regexp.MustCompile(`ftid=([a-zA-Z0-9_-]+)`)
		matches := pattern.FindStringSubmatch(mapURL)
		if len(matches) >= 2 {
			return matches[1]
		}
	}
	
	// Look for Place ID in other formats
	patterns := []string{
		`place_id:([a-zA-Z0-9_-]+)`,     // place_id:ID
		`!1s([a-zA-Z0-9_-]+)!`,          // !1sID!
		`data=.*!1m.*!1s([a-zA-Z0-9_-]+)`, // data parameter with Place ID
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(mapURL)
		if len(matches) >= 2 {
			return matches[1]
		}
	}
	
	return ""
}

// convertCIDToPlaceID converts CID format to coordinates by decoding the hexadecimal values
func convertCIDToPlaceID(cid string) (string, error) {
	// CID format: 0x6001078ce753f549:0x51a6cbb1594b20db
	// The first part encodes latitude, the second part encodes longitude
	
	parts := strings.Split(cid, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid CID format: %s", cid)
	}
	
	// For simplicity, we'll fallback to coordinate extraction from the URL
	// since CID decoding requires complex mathematical conversion
	fmt.Printf("DEBUG: CID found but using coordinate extraction fallback for: %s\n", cid)
	return "", fmt.Errorf("using coordinate extraction fallback for CID: %s", cid)
}

// getCoordinatesFromPlaceID gets coordinates using Places API (New)
func getCoordinatesFromPlaceID(placeID string) (float64, float64, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return 0, 0, fmt.Errorf("Google Maps API key not found")
	}
	
	// Use Places API (New) to get place details
	detailsURL := fmt.Sprintf("https://places.googleapis.com/v1/places/%s?fields=location&languageCode=ja&key=%s", placeID, apiKey)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(detailsURL)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to make details request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("details API returned status: %d", resp.StatusCode)
	}
	
	var detailsResp struct {
		Location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&detailsResp); err != nil {
		return 0, 0, fmt.Errorf("failed to decode details response: %v", err)
	}
	
	return detailsResp.Location.Latitude, detailsResp.Location.Longitude, nil
}

// extractCoordinatesFromURL extracts latitude and longitude from Google Maps URL
func extractCoordinatesFromURL(mapURL string) (float64, float64, error) {
	// First, try to extract Place ID from the URL and use Places API
	if placeID := extractPlaceIDFromURL(mapURL); placeID != "" {
		fmt.Printf("DEBUG: Found Place ID in URL: %s\n", placeID)
		lat, lng, err := getCoordinatesFromPlaceID(placeID)
		if err == nil {
			fmt.Printf("DEBUG: Successfully got coordinates from Place ID: %f, %f\n", lat, lng)
			return lat, lng, nil
		}
		fmt.Printf("DEBUG: Failed to get coordinates from Place ID: %v\n", err)
	}

	// Look for coordinates in various formats
	patterns := []string{
		`!3d(-?\d+\.\d+)!4d(-?\d+\.\d+)`,           // !3dlat!4dlng
		`ll=(-?\d+\.\d+),(-?\d+\.\d+)`,             // ll=lat,lng
	}

	for i, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(mapURL)
		if len(matches) >= 3 {
			lat, err1 := strconv.ParseFloat(matches[1], 64)
			lng, err2 := strconv.ParseFloat(matches[2], 64)
			if err1 == nil && err2 == nil {
				fmt.Printf("DEBUG: Coordinate pattern %d matched - Lat: %f, Lng: %f\n", i, lat, lng)
				return lat, lng, nil
			}
		}
	}

	// If URL patterns fail, try to extract from HTML content (for HTMLスクレイピング cases)
	lat, lng, err := extractCoordinatesFromHTML(mapURL)
	if err == nil {
		fmt.Printf("DEBUG: Coordinates extracted from HTML - Lat: %f, Lng: %f\n", lat, lng)
		return lat, lng, nil
	}

	fmt.Printf("DEBUG: All coordinate extraction methods failed for URL: %s\n", mapURL)
	return 0, 0, fmt.Errorf("coordinates not found in URL")
}

// extractCoordinatesFromHTML extracts coordinates from HTML page content
func extractCoordinatesFromHTML(mapURL string) (float64, float64, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make request to Google Maps URL
	req, err := http.NewRequest("GET", mapURL, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create request: %v", err)
	}

	// Set user agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read response body: %v", err)
	}

	htmlContent := string(body)
	fmt.Printf("DEBUG: HTML content length for coordinate extraction: %d characters\n", len(htmlContent))

	// Pattern 1: JSON array format like [[209103.37989011017,135.7217792,35.0257152],[0,0,0],[1024,768],13.1]
	jsonPattern := `\[\[[\d\.]+,(-?\d+\.\d+),(-?\d+\.\d+)\]`
	re := regexp.MustCompile(jsonPattern)
	matches := re.FindStringSubmatch(htmlContent)
	if len(matches) >= 3 {
		lng, err1 := strconv.ParseFloat(matches[1], 64)
		lat, err2 := strconv.ParseFloat(matches[2], 64)
		if err1 == nil && err2 == nil {
			fmt.Printf("DEBUG: JSON pattern matched - Lng: %f, Lat: %f\n", lng, lat)
			return lat, lng, nil
		}
	}

	// Pattern 2: Look for coordinates in script tags or data attributes
	scriptPatterns := []string{
		`"lat":(-?\d+\.\d+),"lng":(-?\d+\.\d+)`,              // {"lat":35.123,"lng":135.456}
		`latitude["\s]*[:=]\s*(-?\d+\.\d+)[\s\S]*?longitude["\s]*[:=]\s*(-?\d+\.\d+)`, // latitude:35.123 longitude:135.456
		`center.*?(-?\d+\.\d+)[,\s]+(-?\d+\.\d+)`,            // center: 35.123, 135.456
	}

	for i, pattern := range scriptPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(htmlContent)
		if len(matches) >= 3 {
			lat, err1 := strconv.ParseFloat(matches[1], 64)
			lng, err2 := strconv.ParseFloat(matches[2], 64)
			if err1 == nil && err2 == nil {
				fmt.Printf("DEBUG: Script pattern %d matched - Lat: %f, Lng: %f\n", i, lat, lng)
				return lat, lng, nil
			}
		}
	}

	return 0, 0, fmt.Errorf("coordinates not found in HTML content")
}

// // extractAdditionalInfoFromPage attempts to extract additional information from the Google Maps page
// func extractAdditionalInfoFromPage(mapURL string) (address, websiteURL, storeName string, err error) {
// 	// Create HTTP client with timeout
// 	client := &http.Client{
// 		Timeout: 10 * time.Second,
// 	}

// 	// Make request to Google Maps URL
// 	req, err := http.NewRequest("GET", mapURL, nil)
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to create request: %v", err)
// 	}

// 	// Set user agent to avoid blocking
// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to fetch page: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return "", "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// 	}

// 	// Read response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to read response body: %v", err)
// 	}

// 	pageContent := string(body)
	
// 	// Debug: Print samples of HTML content to understand structure
// 	fmt.Printf("DEBUG: HTML length: %d characters\n", len(pageContent))
	
// 	// Look for postal code (〒) patterns in HTML
// 	keyword := "〒"
// 	keywordIndex := strings.Index(pageContent, keyword)
// 	if keywordIndex != -1 {
// 		start := keywordIndex - 100
// 		if start < 0 {
// 			start = 0
// 		}
// 		end := keywordIndex + 200
// 		if end > len(pageContent) {
// 			end = len(pageContent)
// 		}
// 		fmt.Printf("DEBUG: HTML sample around '%s': %s\n", keyword, pageContent[start:end])
// 	} else {
// 		fmt.Printf("DEBUG: No postal code '〒' found in HTML\n")
// 	}

// 	// Extract address using regex patterns
// 	address = extractAddressFromHTML(pageContent)
// 	fmt.Printf("DEBUG: Extracted address: '%s'\n", address)
// 	if address == "" {
// 		address = "Address not available"
// 	}

// 	// Extract website URL using regex patterns
// 	websiteURL = extractWebsiteFromHTML(pageContent)
	
// 	// Extract store name from HTML
// 	storeName = extractStoreNameFromHTML(pageContent)

// 	return address, websiteURL, storeName, nil
// }

// // extractAddressFromHTML extracts address from HTML content
// func extractAddressFromHTML(html string) string {
// 	// Look for addresses starting with 〒 (postal code) - highest priority
// 	postalPatterns := []string{
// 		// Addresses starting with 〒 in various HTML contexts
// 		`"([^"]*〒[^"]*)"`,                      // JSON format with postal code
// 		`>([^<]*〒[^<]*)</[^>]*>`,               // Any HTML tag containing postal address
// 		`<div[^>]*>([^<]*〒[^<]*)</div>`,        // Div containing postal address
// 		`<span[^>]*>([^<]*〒[^<]*)</span>`,      // Span containing postal address
// 		`<button[^>]*>.*?([^<]*〒[^<]*).*?</button>`, // Button containing postal address
// 		`data-value="([^"]*〒[^"]*)"`,           // Data attribute with postal address
// 		`aria-label="[^"]*([^"]*〒[^"]*)"`,      // Aria label with postal address
// 		`content="([^"]*〒[^"]*)"`,              // Meta content with postal address
		
// 		// More flexible patterns for postal addresses
// 		`([^>\s"]*〒\s*\d{3}-?\d{4}[^<"]*[都道府県][^<"]*[市区町村][^<"]*)`, // Complete postal address
// 		`([^>\s"]*〒\s*\d{3}-?\d{4}[^<"]*)`,     // Any postal code pattern
// 	}
	
// 	fmt.Printf("DEBUG: Searching for addresses starting with 〒\n")
	
// 	for i, pattern := range postalPatterns {
// 		re := regexp.MustCompile(pattern)
// 		matches := re.FindAllStringSubmatch(html, -1) // Find all matches
		
// 		for _, match := range matches {
// 			if len(match) >= 2 {
// 				address := strings.TrimSpace(match[1])
// 				fmt.Printf("DEBUG: Postal pattern %d found potential address: '%s'\n", i, address)
				
// 				// Check if address starts with 〒
// 				if strings.HasPrefix(address, "〒") && isValidPostalAddress(address) {
// 					fmt.Printf("DEBUG: Valid postal address found: '%s'\n", address)
// 					return address
// 				}
// 			}
// 		}
// 	}
	
// 	fmt.Printf("DEBUG: No addresses starting with 〒 found\n")
// 	return ""
// }

// // isValidPostalAddress checks if the postal address starting with 〒 is valid
// func isValidPostalAddress(address string) bool {
// 	if address == "" || !strings.HasPrefix(address, "〒") {
// 		return false
// 	}
	
// 	// Check for minimum length
// 	if len(address) < 8 { // 〒 + 3 digits + - + 4 digits minimum
// 		fmt.Printf("DEBUG: Postal address too short: %s\n", address)
// 		return false
// 	}
	
// 	// Check for invalid patterns
// 	invalidPatterns := []string{
// 		"[[", "]]", "null", "{", "}", "[null",
// 		"google.com", "maps.google", "data-",
// 		"aria-", "class=", "id=", "onclick=",
// 		"javascript:", "function(", "var ",
// 		"undefined", "NaN", "Infinity",
// 	}
	
// 	lowerAddress := strings.ToLower(address)
// 	for _, pattern := range invalidPatterns {
// 		if strings.Contains(lowerAddress, pattern) {
// 			fmt.Printf("DEBUG: Postal address rejected due to invalid pattern '%s': %s\n", pattern, address)
// 			return false
// 		}
// 	}
	
// 	// Check for postal code pattern: 〒XXX-XXXX or 〒XXXXXXX
// 	postalCodePattern := `〒\s*\d{3}-?\d{4}`
// 	matched, _ := regexp.MatchString(postalCodePattern, address)
// 	if !matched {
// 		fmt.Printf("DEBUG: Invalid postal code format: %s\n", address)
// 		return false
// 	}
	
// 	// Check if it contains meaningful address components
// 	addressIndicators := []string{
// 		"都", "道", "府", "県", "市", "区", "町", "村",
// 		"丁目", "番地", "号", "番", "条", "丁", "目",
// 		"街", "通", "路", "駅", "店", "ビル", "マンション",
// 	}
	
// 	hasIndicator := false
// 	for _, indicator := range addressIndicators {
// 		if strings.Contains(address, indicator) {
// 			hasIndicator = true
// 			break
// 		}
// 	}
	
// 	if !hasIndicator {
// 		fmt.Printf("DEBUG: Postal address lacks location indicators: %s\n", address)
// 		// For postal addresses, we're more lenient as they might be abbreviated
// 		return len(address) >= 10 // Allow if it's reasonably long
// 	}
	
// 	fmt.Printf("DEBUG: Postal address validation passed: %s\n", address)
// 	return true
// }

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

// parseBusinessHoursFromWeekdayDescriptions parses Google Maps weekdayDescriptions into detailed business hours struct
func parseBusinessHoursFromWeekdayDescriptions(weekdayDescriptions []string) models.BusinessHoursData {
	if len(weekdayDescriptions) == 0 {
		return getDefaultBusinessHoursData()
	}

	businessHours := models.BusinessHoursData{}
	
	// Map of weekday keys to their corresponding day schedule
	dayMap := map[string]*models.DaySchedule{
		"monday":    &businessHours.Monday,
		"tuesday":   &businessHours.Tuesday,
		"wednesday": &businessHours.Wednesday,
		"thursday":  &businessHours.Thursday,
		"friday":    &businessHours.Friday,
		"saturday":  &businessHours.Saturday,
		"sunday":    &businessHours.Sunday,
	}
	
	// Parse each weekday description
	for i, desc := range weekdayDescriptions {
		dayKey := getDayKeyFromIndex(i)
		if daySchedule, exists := dayMap[dayKey]; exists {
			if strings.Contains(desc, "定休日") || strings.Contains(desc, "休業") {
				daySchedule.IsClosed = true
				daySchedule.TimeSlots = []models.TimeSlot{}
			} else {
				daySchedule.IsClosed = false
				// Parse time ranges from description
				timeRanges := parseDetailedTimeRanges(desc)
				if len(timeRanges) == 0 {
					// Default time slot if parsing fails
					timeRanges = []models.TimeSlot{{
						OpenTime:  "11:00",
						CloseTime: "22:00",
					}}
				}
				daySchedule.TimeSlots = timeRanges
			}
		}
	}
	
	return businessHours
}

// getDefaultBusinessHoursData returns default business hours struct
func getDefaultBusinessHoursData() models.BusinessHoursData {
	defaultSlot := models.TimeSlot{
		OpenTime:  "11:00",
		CloseTime: "22:00",
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

// getDayKeyFromIndex returns day key from weekday index
func getDayKeyFromIndex(index int) string {
	days := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
	if index >= 0 && index < len(days) {
		return days[index]
	}
	return "monday"
}

// parseDetailedTimeRanges extracts detailed time ranges from description text
func parseDetailedTimeRanges(desc string) []models.TimeSlot {
	var timeSlots []models.TimeSlot
	
	// Simple time pattern matching for Japanese format
	// Example: "11時00分～20時00分"
	timePattern := regexp.MustCompile(`(\d{1,2})時(\d{2})分～(\d{1,2})時(\d{2})分`)
	matches := timePattern.FindAllStringSubmatch(desc, -1)
	
	for _, match := range matches {
		if len(match) >= 5 {
			openHour := match[1]
			openMin := match[2]
			closeHour := match[3]
			closeMin := match[4]
			
			timeSlot := models.TimeSlot{
				OpenTime:  fmt.Sprintf("%02s:%s", openHour, openMin),
				CloseTime: fmt.Sprintf("%02s:%s", closeHour, closeMin),
			}
			timeSlots = append(timeSlots, timeSlot)
		}
	}
	
	// If no time ranges found, return empty slice
	return timeSlots
}

// extractWebsiteFromHTML extracts website URL from HTML content (placeholder)
func extractWebsiteFromHTML(html string) string {
	// Simple website URL extraction
	patterns := []string{
		`"website_url":"([^"]+)"`,
		`href="(https?://[^"]+)"`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) >= 2 {
			return matches[1]
		}
	}
	
	return ""
}
