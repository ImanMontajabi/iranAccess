package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DomainCheckResult struct {
	Domain     string `json:"domain"`
	StatusCode int    `json:"statusCode"`
	IsUp       bool   `json:"isUp"`
	Error      error  `json:"-"`
}

func main() {
	app := fiber.New()
	app.Get("/check", checkDomainsHandler)
	fmt.Println("Server is up on port 3000")
	log.Fatal(app.Listen(":3000"))
}

func checkDomainsHandler(c *fiber.Ctx) error {
	file, err := os.Open("./domains.csv")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "CSV file not found or failed to open!",
		})
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Couldn't read the CSV records",
		})
	}
	if len(records) == 0 {
		return c.JSON([]DomainCheckResult{})
	}

	resultChannel := make(chan DomainCheckResult)
	wg := sync.WaitGroup{}
	for _, record := range records {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			checkDomain(d, resultChannel)
		}(record[0])
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	var results []DomainCheckResult
	for result := range resultChannel {
		results = append(results, result)
	}

	return c.JSON(results)
}

func checkDomain(domain string, resultChannel chan<- DomainCheckResult) {
	url := "https://" + domain
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	result := DomainCheckResult{Domain: domain}

	if err != nil {
		result.Error = err
		result.IsUp = false
	} else {
		defer resp.Body.Close()
		result.StatusCode = resp.StatusCode
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			result.IsUp = true
		} else {
			result.IsUp = false
		}
	}
	resultChannel <- result
}
