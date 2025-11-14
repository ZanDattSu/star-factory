package path

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetProjectRoot ищет корневую директорию проекта по наличию go.work файла
func GetProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("не удалось получить рабочую директорию: " + err.Error())
	}

	for {
		_, err = os.Stat(filepath.Join(dir, "go.work"))
		if err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("не удалось найти корень проекта (go.work)")
		}

		dir = parent
	}
}

// GetPathRelativeToRoot получает путь относительно корня проекта
func GetPathRelativeToRoot(relativePath string) string {
	return fmt.Sprintf("%s/%s", GetProjectRoot(), relativePath)
}
