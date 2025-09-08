package checker

import (
	"encoding/csv"
	"fmt"
	"iranAccess/internal/cache"
	"iranAccess/internal/models"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// DomainChecker handles domain checking operations
type DomainChecker struct {
	httpClient *http.Client
	csvPath    string
}

// NewDomainChecker creates a new domain checker instance
func NewDomainChecker(csvPath string) *DomainChecker {
	return &DomainChecker{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		csvPath: csvPath,
	}
}

// CheckDomain checks if a single domain is accessible
func (dc *DomainChecker) CheckDomain(domain string) models.DomainCheckResult {
	url := "https://" + domain
	resp, err := dc.httpClient.Get(url)

	result := models.DomainCheckResult{Domain: domain}

	if err != nil {
		result.Error = err.Error()
		result.IsUp = false
		return result
	}

	defer resp.Body.Close()
	result.StatusCode = resp.StatusCode
	result.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 400

	return result
}

// LoadDomainsFromCSV reads domains from CSV file
func (dc *DomainChecker) LoadDomainsFromCSV() ([]string, error) {
	file, err := os.Open(dc.csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) <= 1 {
		return nil, fmt.Errorf("no domains found in CSV file")
	}

	domains := make([]string, 0, len(records)-1)
	for _, record := range records[1:] { // Skip header
		if len(record) > 0 && record[0] != "" {
			domains = append(domains, record[0])
		}
	}

	return domains, nil
}

// CheckAllDomains checks all domains concurrently
func (dc *DomainChecker) CheckAllDomains(domains []string) []models.DomainCheckResult {
	resultChannel := make(chan models.DomainCheckResult, len(domains))
	var wg sync.WaitGroup

	for _, domain := range domains {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			result := dc.CheckDomain(d)
			resultChannel <- result
		}(domain)
	}

	wg.Wait()
	close(resultChannel)

	results := make([]models.DomainCheckResult, 0, len(domains))
	for result := range resultChannel {
		results = append(results, result)
	}

	return results
}

// StartDomainChecker starts the domain checking worker
func (dc *DomainChecker) StartDomainChecker() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	// Run immediately
	dc.runCheck()

	// Then run every 10 minutes
	for range ticker.C {
		dc.runCheck()
	}
}

func (dc *DomainChecker) runCheck() {
	log.Println("Starting domain check...")

	domains, err := dc.LoadDomainsFromCSV()
	if err != nil {
		log.Printf("Error loading domains: %v", err)
		return
	}

	if len(domains) == 0 {
		log.Println("No domains to check")
		return
	}

	results := dc.CheckAllDomains(domains)
	cache.DomainCache.SetResults(results)

	log.Printf("Domain checking completed. %d domains checked and cached", len(results))
}
