package fileworker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type File struct {
	filename string
	filepath string
	dateest  string
}

func moveFileToFolder(filePath, folderName string) error {
	// Создаем путь к новому местоположению файла
	destDir := filepath.Join(filepath.Dir(filePath), folderName)
	destPath := filepath.Join(destDir, filepath.Base(filePath))

	// MkdirAll НЕ выдает ошибку, если папка уже существует
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return fmt.Errorf("не удалось создать папку %s: %w", destDir, err)
	}

	// Перемещаем файл
	err = os.Rename(filePath, destPath)
	if err != nil {
		return fmt.Errorf("не удалось переместить файл %s в %s: %w", filePath, destPath, err)
	}

	return nil
}

func checkFileFormat(filename string) string {
	// Безопасное получение расширения файла
	path_splited := strings.Split(filename, "/")
	pointed_file := path_splited[len(path_splited)-1]

	parts := strings.Split(pointed_file, ".")
	if len(parts) < 2 {
		return "" // нет расширения
	}

	return parts[len(parts)-1]
}

func FileSorter(path *string) {
	var files []string
	filepath.WalkDir(*path, func(path string, d os.DirEntry, err error) error {
		if !d.IsDir() { // Добавляем только файлы, пропускаем папки
			files = append(files, path)
		}
		return nil
	})

	var wg sync.WaitGroup

	for i := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			file_format := checkFileFormat(filePath)
			switch file_format {
			case "md", "txt", "pdf", "doc", "docx", "xlsx", "xls":
				moveFileToFolder(filePath, "Documents")
			case "mp4", "avi", "mkv", "mov":
				moveFileToFolder(filePath, "Videos")
			case "jpg", "png", "gif", "webp":
				moveFileToFolder(filePath, "Images")
			default:
				moveFileToFolder(filePath, "Others")
			}
		}(files[i]) //  это  корчое передается  в горутину значение
	}

	wg.Wait()
}
