package config

import "flag"

// Config for all servers, apps for the astralis project.
type Config interface {
	// Server
	APIPort() string
	// Third-party APIs
	NasaAPIKey() string
}

type config struct {
	// Server
	apiPort string

	// Third-party APIs
	nasaAPIKey string
}

func LoadConfig() Config {
	// Server
	apiPort := flag.String("api_port", ":8080", "Astralis API port. Defaults to 8080")

	// Third-party APIs
	nasaAPIKey := flag.String("nasa_api_key", "", "NASA API Key")

	flag.Parse()
	return &config{
		apiPort: *apiPort,
		nasaAPIKey: *nasaAPIKey,
	}
}

func (c *config) APIPort() string {
	return c.apiPort
}

func (c *config) NasaAPIKey() string {
	return c.nasaAPIKey
}
