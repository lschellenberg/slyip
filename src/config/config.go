package config

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	DefaultConfigLocation = "yip.json"
)

type Config struct {
	DB        SqlDBInfo      `json:"db"`
	JWT       JWTTokenConfig `json:"jwt"`
	Audiences []Audience     `json:"audiences"`
	Clients   []Client       `json:"clients"`
	API       API            `json:"api"`
	Test      Test           `json:"test"`
	Email     EmailConfig    `json:"email"`
	EthConfig EthConfig      `json:"eth"`
}

type EmailConfig struct {
	SenderName  string        `json:"senderName"`
	SenderEmail string        `json:"senderEmail"`
	MailJet     MailJetConfig `json:"mailjet"`
}

type MailJetConfig struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
	TemplateId int    `json:"templateId"`
}

func ReadConfig() (Config, error) {
	// TODO check for custom locations of config file
	c, err := readConfigFile(DefaultConfigLocation)
	if err != nil {
		return Config{}, err
	}

	err = c.verifyConfig()
	return c, err
}

func (c Config) verifyConfig() error {
	_, err := url.Parse(c.JWT.Issuer)
	if err != nil {
		return fmt.Errorf("malformed issuer url in config file: %s", c.JWT.Issuer)
	}

	for _, aud := range c.Audiences {
		_, err := url.Parse(aud.URL)
		if err != nil {
			return fmt.Errorf("malformed audience url in config file: %s", aud.URL)
		}
		for _, cli := range aud.Clients {
			_, err := url.Parse(cli)
			if err != nil {
				return fmt.Errorf("malformed client url in config file: %s", cli)
			}
		}
	}

	return nil
}

func (c Config) String() string {
	marshaled, err := json.MarshalIndent(c, "", "   ")
	if err != nil {
		return "error printing json"
	}
	return string(marshaled)
}

type API struct {
	Port      string `json:"port"`
	SwaggerOn bool   `json:"swagger_on"`
	Admin     Admin  `json:"admin"`
}

type Admin struct {
	Username       string `json:"username"`
	PasswordHashed string `json:"password_hashed"`
}

type JWTTokenConfig struct {
	TokenExpirationInSec        int64  `json:"token_expiration_in_sec"`
	RefreshTokenExpirationInSec int64  `json:"refresh_token_expiration_in_sec"`
	CertificatePrivate          string `json:"certificate_private"`
	CertificatePublic           string `json:"certificate_public"`
	Issuer                      string `json:"issuer"`
}

type Audience struct {
	ID      string   `json:"id"`
	URL     string   `json:"url"`
	Clients []string `json:"clients"`
	Scopes  []string `json:"scopes"`
}

type Client struct {
	ID        string   `json:"id"`
	Domain    string   `json:"domain"`
	Label     string   `json:"label"`
	Audiences []string `json:"audiences"`
}

func (c Config) AudiencesByClient(clientId string) []string {
	audiences := make([]string, 0)
	for _, a := range c.Audiences {
		for _, cl := range a.Clients {
			if cl == clientId {
				audiences = append(audiences, a.URL)
			}
		}
	}

	return audiences
}

type SqlDBInfo struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Port     string `json:"port"`
}

type Test struct {
	On bool `json:"on"`
}

func (pi SqlDBInfo) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", pi.Host, pi.Port, pi.User, pi.Password, pi.Database)
}

func (c Config) VerifyAudiencesExist(audience []string) bool {
	exist := false
	for _, a := range audience {
		exist = false
		for _, v := range c.Audiences {
			if v.URL == a {
				exist = true
			}
		}
		// no matter - if one audience doesn't exist the whole thing is false
		if !exist {
			return false
		}
	}
	return true
}
