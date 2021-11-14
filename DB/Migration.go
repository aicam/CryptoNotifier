package DB

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type WebData struct {
	gorm.Model
	Status      string `json:"status"`
	Error       string `json:"error"`
	AirdropName string `json:"airdrop_name"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	UniqueID    string `json:"unique_id"`
}

type UsersData struct {
	gorm.Model
	Username   string    `json:"username"`
	LastOnline time.Time `json:"last_online"`
}

func DbSqlMigration(url string) *gorm.DB {
	db, err := gorm.Open("mysql", url)
	if err != nil {
		log.Println(err)
	}
	db.AutoMigrate(&WebData{})
	db.AutoMigrate(&UsersData{})
	return db
}
