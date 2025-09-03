package models

// DomainCheckResult represents the result of a domain check
type DomainCheckResult struct {
	Domain     string `json:"domain"`
	StatusCode int    `json:"statusCode"`
	IsUp       bool   `json:"isUp"`
	Error      string `json:"error,omitempty"`
}
