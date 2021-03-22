package services

import (
	"os"
	"sync"

	"github.com/heyjoakim/devops-21/helpers"
	"github.com/heyjoakim/devops-21/models"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// DBContext defines the application
type DBContext struct {
	db *gorm.DB
}

var (
	dsn         string
	environment string
	dbContext   DBContext
	lock        = &sync.Mutex{}
)

func configureEnv() {
	envFilePath := helpers.GetFullPath("../.env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Error("Error loading .env file - using system variables.")

	}

	environment = os.Getenv("ENVIRONMENT")
	dsn = os.Getenv("DB_CONNECTION")
}

func (d *DBContext) initialize() {
	configureEnv()
	db, err := d.connectDB()
	if err != nil {
		log.Error(err)
	}
	d.db = db
}

func (d *DBContext) connectDB() (*gorm.DB, error) {
	switch environment {
	case "develop":
		log.Info("Using local SQLite db")
		return gorm.Open(sqlite.Open("./tmp/minitwit.db"), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	case "production":
		log.Info("Using remote postgres db")
		return gorm.Open(postgres.New(
			postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true, // disables implicit prepared statement usage
			}),
			&gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
			})
	case "testing":
		log.Info("Using in memory SQLite db")

		return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	}
	log.Panic("Environment is not specified in the .env file")
	return nil, nil
}

// initDB creates the database tables.
func (d *DBContext) initDB() {
	err := d.db.AutoMigrate(&models.User{}, &models.Follower{}, &models.Message{}, &models.Config{})
	if err != nil {
		log.Fatal("Migration error:", err)
	}
}

// GetDBInstance returns DBContext with specific environment db
func GetDBInstance() DBContext {
	if dbContext.db == nil {
		lock.Lock()
		defer lock.Unlock()
		if dbContext.db == nil {
			log.Info("Creating Single Instance Now")
			dbContext.initialize()
			dbContext.initDB() // AutoMigrate
		}
	}
	return dbContext
}
