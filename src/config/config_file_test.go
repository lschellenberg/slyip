package config

import "testing"

func TestReadConfigFile(t *testing.T) {
	c := Config{
		DB: SqlDBInfo{
			Host:     "localhost",
			User:     "postgres",
			Password: "postgres",
			Database: "postgres",
			Port:     "5432",
		},
		JWT: JWTTokenConfig{
			TokenExpirationInSec:        600,
			RefreshTokenExpirationInSec: 999999999999,
			CertificatePrivate:          "app.rsa",
			CertificatePublic:           "app.rsa.pub",
			Issuer:                      "https://ip.yours.net",
		},
		Audiences: []Audience{
			{
				URL:     "https://api.yours.net",
				Clients: []string{"https://sample.yours.net"},
			},
		},
		API: API{
			Port:      "8080",
			SwaggerOn: true,
		},
		Chains: []Chain{
			{
				RPCUrl: "https://rpc",
				ID:     "1",
			},
		},
	}

	saveConfigFile(c, "test.json")
}
