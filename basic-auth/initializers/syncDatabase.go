package initializers

import "github.com/LoveKhatri/basic-auth/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
