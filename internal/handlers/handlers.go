package handlers

import (
	"iranAccess/internal/cache"

	"github.com/gofiber/fiber/v2"
)

// CheckDomainsHandler handles the domain check API endpoint
func CheckDomainsHandler(c *fiber.Ctx) error {
	results := cache.DomainCache.GetResults()
	return c.JSON(results)
}
