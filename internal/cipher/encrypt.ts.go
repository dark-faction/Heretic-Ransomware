package cipher

import (
	cryptorand "crypto/rand"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/chacha20poly1305"
	"io"
	"log"
	"os"
)

type FilePathInfo struct {
	Path          string `json:"path,omitempty"`
	FileName      string `json:"file_name,omitempty"`
	FileExtension string `json:"file_extension,omitempty"`
	Key           string `json:"key,omitempty"`
}

const (
	SaltSize   = 32         // in bytes
	NonceSize  = 24         // in bytes. taken from aead.NonceSize()
	KeySize    = uint32(32) // KeySize is 32 bytes (256 bits).
	KeyTime    = uint32(5)
	KeyMemory  = uint32(1024 * 64) // KeyMemory in KiB. here, 64 MiB.
	KeyThreads = uint8(4)
	chunkSize  = 1024 * 32 // chunkSize in bytes. here, 32 KiB.
)

func Encrypt(filePathInfo FilePathInfo) {

	salt := make([]byte, SaltSize)
	if n, err := cryptorand.Read(salt); err != nil || n != SaltSize {
		log.Println("Error when generating random salt.")
		panic(err)
	}

	outfile, err := os.OpenFile(fmt.Sprintf("%s.heretic", filePathInfo.Path), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Error when opening/creating output file.")
		zap.L().Info(err.Error())
	}
	defer func(outfile *os.File) {
		_ = outfile.Close()
	}(outfile)

	_, _ = outfile.Write(salt)

	aead, err := chacha20poly1305.NewX([]byte(filePathInfo.Key))
	if err != nil {
		zap.L().Info(err.Error())
	}

	infile, err := os.Open(filePathInfo.Path)
	if err != nil {
		log.Println("Error when opening input file.")
		panic(err)
	}
	defer func(infile *os.File) {
		_ = infile.Close()
	}(infile)

	buf := make([]byte, chunkSize)
	adCounter := 0 // associated data is a counter

	for {
		n, err := infile.Read(buf)

		if n > 0 {
			// Select a random nonce, and leave capacity for the ciphertext.
			nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+n+aead.Overhead())
			if m, err := cryptorand.Read(nonce); err != nil || m != aead.NonceSize() {
				log.Println("Error when generating random nonce :", err)
				log.Println("Generated nonce is of following size. m : ", m)
				panic(err)
			}

			msg := buf[:n]
			// Encrypt the message and append the ciphertext to the nonce.
			encryptedMsg := aead.Seal(nonce, nonce, msg, []byte(string(rune(adCounter))))
			_, _ = outfile.Write(encryptedMsg)
			adCounter += 1
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error when reading input file chunk :", err)
			panic(err)
		}
	}
}
