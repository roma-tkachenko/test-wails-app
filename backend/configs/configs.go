package configs

const (
	// Application settings
	AppName    = "MyApp"
	AppVersion = "1.0.0"

	// Predefined URLS
	BaseURL       = "https://animestars.org/"
	СontrollerURL = "https://animestars.org/engine/ajax/controller.php"
	BoostClubURL  = "https://animestars.org/clubs/19/boost/"

	// Reward cards
	ClaimRewardCardInterval = 3 // value in minutes

	// Boos cards
	RefreshCardInterval    = 10  // value in milliseconds
	RetryBoostCardInterval = 100 // value in milliseconds

	// Database settings
	DatabaseHost     = "localhost"
	DatabasePort     = 5432
	DatabaseUser     = "postgres"
	DatabasePassword = "password"
	DatabaseName     = "myapp_db"

	// Logging settings
	LogLevel   = "DEBUG"
	LogFile    = "logs/app.log"
	MaxLogSize = 10 // Максимальний розмір лог-файлу в MB

	// Security settings
	EncryptionKey = "replace_this_with_a_secure_key"
	JWTSecret     = "replace_this_with_a_secure_jwt_key"
)

var UserHash string
