package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	"github.com/nats-io/nats.go"
)

// Config is a structure for configuration
type Config struct {
	Database   Database
	Service    Service
	Mail       Mail
	JWT        JWT
	NatsServer NatsServer
	Hubspot    Hubspot
}

// Database is a structure for relational database configuration
type Database struct {
	Name  string
	Host  string
	Port  string
	User  string
	Pass  string
	Ssl   string
	Ideal string
}

// Mail is a structure for mail service configuration
type Mail struct {
	Host string
	Port string
	From string
	User string
	Pass string
}

// Hubspot is a structure for store the Hubspot information
type Hubspot struct {
	WebsiteID string
	PortalID  string
	GitID     string
}
type JWT struct {
	PrivateKey    string
	JWTExpireTime time.Duration
}

//NatsServer : for nats connection parameters
type NatsServer struct {
	URL      string
	Token    string
	Username string
	Password string
	Subject  string
	Queue    string
}

// Service is a structure for service specific related configuration
type Service struct {
	Port           string
	Environment    string
	Build          string
	RepoURL        string
	RepoPrivateKey string
	SupportEmails  string
	Workspace      string
}

var (
	// Conf is a global variable for configuration
	Conf Config
	// TomlFile is a global variable for toml file path
	TomlFile string
	// DB Database client
	DB *gorm.DB
	//NC for nats connection
	NC *nats.Conn
)

const (
	PersistStoragePath string = "./scripts/websiteScan/"
	GitPath            string = "./scripts/github-scan/"
	MailService        string = "false"
)

// ConfigurationWithEnv is a method to initialize configuration with environment variables
func ConfigurationWithEnv() {

	Conf.JWT.JWTExpireTime = time.Minute * 30
	Conf.JWT.PrivateKey = os.Getenv("PRIVATE_KEY")
	if Conf.JWT.PrivateKey == "" {
		Conf.JWT.PrivateKey = ""
	}

	///mail configuration
	// mail service configuration
	Conf.Mail.Host = os.Getenv("SECURITY_MAIL_SMTP_HOST")
	Conf.Mail.Port = os.Getenv("SECURITY_MAIL_SMTP_PORT")
	Conf.Mail.From = os.Getenv("SECURITY_MAIL_FROM")
	Conf.Mail.User = os.Getenv("SECURITY_MAIL_USERID")
	Conf.Mail.Pass = os.Getenv("SECURITY_MAIL_PASS")

	// cockroach database configuration
	Conf.Database.Host = os.Getenv("SECURITY_DB_HOST")
	Conf.Database.Port = os.Getenv("SECURITY_DB_PORT")
	Conf.Database.User = os.Getenv("SECURITY_DB_USER")
	Conf.Database.Port = os.Getenv("SECURITY_DB_PASS")
	Conf.Database.Name = os.Getenv("SECURITY_DB_NAME")
	Conf.Database.Ideal = os.Getenv("SECURITY_DB_IDEAL_CONNECTIONS")
	Conf.Database.Ssl = "disable"
	// if service port is not defined set default port
	if os.Getenv("SECURITY_PORT") != "" {
		Conf.Service.Port = os.Getenv("SECURITY_PORT")
	} else {
		Conf.Service.Port = "8000"
	}
	Conf.Service.Environment = os.Getenv("ENVIRONMENT")
	Conf.Service.Build = os.Getenv("BUILD_IMAGE")
	Conf.Service.RepoURL = os.Getenv("REPO_URL_OF_SCRIPT")
	Conf.Service.RepoPrivateKey = os.Getenv("PRIVATE_KEY_OF_REPO")

	Conf.Service.SupportEmails = os.Getenv("SECURITY_SUPPORT_EMAILS")

	//base url
	if Conf.Service.Workspace == "" {
		Conf.Service.Workspace = "https://continuous-security.secops.neuralcompany.team/api/workspace"
	}
}

// ConfigurationWithToml is a method to initialize configuration with toml file
func ConfigurationWithToml(filePath string) error {
	// set varible as file path if configuration is done using toml
	TomlFile = filePath
	// parse toml file and save data config structure
	_, err := toml.DecodeFile(filePath, &Conf)
	if err != nil {
		log.Println(err)
		return err
	}
	if Conf.JWT.PrivateKey == "" {
		Conf.JWT.PrivateKey = ""
	}

	if Conf.Service.Port == "" {
		Conf.Service.Port = "8000"
	}
	Conf.Service.Build = os.Getenv("BUILD_IMAGE")
	Conf.Database.Ssl = "disable"

	if Conf.Service.Workspace == "" {
		Conf.Service.Workspace = "https://continuous-security.secops.neuralcompany.team/api/workspace"
	}
	
	return nil
}

// SetConfig is a method to re-intialise configuration at runtime
func SetConfig() {
	if TomlFile == "" {
		ConfigurationWithEnv()
	} else {
		ConfigurationWithToml(TomlFile)
	}
}

// DBConfig is a method that return postgres database connection string
func DBConfig() string {
	//again reset the config if any changes in toml file or environment variables
	SetConfig()
	// creating postgres connection string
	str := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		Conf.Database.Host,
		Conf.Database.Port,
		Conf.Database.User,
		Conf.Database.Pass,
		Conf.Database.Name,
		Conf.Database.Ssl)
	return str
}
