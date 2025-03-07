package cmd

type Environment int

type Config struct {
	yipURL    string
	oracleURL string
}

const (
	Local Environment = 0
	Cloud Environment = 1
)

func GetEnvironment(env Environment) Config {
	switch env {
	case Cloud:
		return getCloudEnvironment()
	default:
		return getLocalEnvironment()
	}
}

func getLocalEnvironment() Config {
	return Config{
		yipURL: "http://localhost:8080",
	}
}

func getCloudEnvironment() Config {
	return Config{
		yipURL: "https://auth.singularry.xyz",
	}
}
