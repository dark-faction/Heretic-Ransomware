package startup

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type EncryptedFileInfo struct {
	gorm.Model
	Path      string `json:"path,omitempty"`
	FileName  string `json:"file_name,omitempty"`
	RenamedTo string `json:"renamed_to,omitempty"`
}

func InitDb() {
	db, err := gorm.Open(sqlite.Open("heretic.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	_ = db.AutoMigrate(&EncryptedFileInfo{})

}
