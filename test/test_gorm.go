package main

import (
	"fmt"
	"ginchat/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&models.UserBasic{})

	// Create
	user := &models.UserBasic{}
	user.Name = "Tony"
	db.Create(user)

	// Read
	fmt.Println(db.First(user, 1)) // Read first user with integer primary key
	// db.First(&user, "name = ?", "Tony") // Read user with name "Tony"
	//db.First(&user, "code = ?", "D42")

	// Update - update user password
	// db.Model(user).Update("PassWord", "1234")
	db.Model(user).Update("PassWord", "1234")
	// Update - update multiple fields
	//db.Model(&user).Updates(UserBasic{Age: 30, Name: "Tom"}) // Only update non-zero fields
	//db.Model(&user).Updates(map[string]interface{}{"Age": 30, "Name": "Tom"})

	// Delete - delete user
	//db.Delete(&user, 1)
}
