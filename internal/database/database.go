package database

import (
	"database/sql"
	"fmt"
	"telebot/internal/models"

	_ "github.com/golang-migrate/migrate/v4/source/file" // required for go-migrate via files
	_ "github.com/lib/pq"                                // required for PostgreSQL connection
)

// PostgresWasteStorage incapsulates PostgreSQL storage
type TelebotLanguageStorage struct {
	db *sql.DB
}

// Migrate ups version of DB model
// func (t *TelebotLanguageStorage) Migrate() {
// 	driver, err := postgres.WithInstance(p.db, &postgres.Config{})
// 	if err != nil {
// 		log.Fatal("[MIGRATE] Unable to get driver due to: " + err.Error())
// 	}
// 	m, err := migrate.NewWithDatabaseInstance(
// 		"file:///app/migrations",
// 		"postgres", driver)
// 	if err != nil {
// 		log.Fatal("[MIGRATE] Unable to get migrate instance due to: " + err.Error())
// 	}
// 	err = m.Up()
// 	switch err {
// 	case migrate.ErrNoChange:
// 		return
// 	default:
// 		log.Fatal("[MIGRATE] Unable to apply DB migrations due to: " + err.Error())
// 	}
// }

func NewTelebotLanguageStorage(config *models.Config) *TelebotLanguageStorage {
	dbURL := config.DbURL
	if dbURL == "" {
		// dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
		dbURL = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DbHost, config.DbPort, config.DbUser, config.DbPassword, config.DbName)
	}
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		fmt.Printf(err.Error())
	}
	if err = db.Ping(); err != nil {
		fmt.Printf(err.Error())
	}

	storage := TelebotLanguageStorage{db: db}
	return &storage
}

func (t *TelebotLanguageStorage) GetLangMessage(langCode int) (*map[int]string, error) {
	var queryLang string
	switch langCode {
	case RU:
		queryLang = "SELECT \"ID\", \"RU\" FROM telebot_language"
	case EN:
		queryLang = "SELECT \"ID\", \"ENG\" FROM telebot_language"
	default:
		queryLang = "SELECT \"ID\", \"RU\" FROM telebot_language"
	}

	rows, err := t.db.Query(queryLang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messageMap := make(map[int]string, 0)
	for rows.Next() {
		var telebotMessage TelebotMessage
		if err := rows.Scan(&telebotMessage.ID, &telebotMessage.Message); err != nil {
			return nil, err
		}

		messageMap[int(telebotMessage.ID.Int32)] = telebotMessage.Message.String
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &messageMap, nil
}
