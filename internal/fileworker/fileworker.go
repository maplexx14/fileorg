package fileworker

import (
	//"fmt"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// type File struct { for future upd
// 	filename string
// 	filepath string
// 	dateest  string
// }

func moveFileToFolder(filePath, folderName string) error {
	// 1. Создаем путь к новому местоположению файла
	destDir := filepath.Join(filepath.Dir(filePath), folderName)
	destPath := filepath.Join(destDir, filepath.Base(filePath))

	// 0755 - стандартные права для папок.
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return fmt.Errorf("не удалось создать папку %s: %w", destDir, err)
	}

	err = os.Rename(filePath, destPath)
	if err != nil {
		return fmt.Errorf("не удалось переместить файл %s в %s: %w", filePath, destPath, err)
	}

	return nil
}
func checkFileFormat(filename string) string {
	//fmt.Print(filename, "\n")
	path_splited := strings.Split(filename, "/")
	pointed_file := path_splited[len(path_splited)-1]
	//fmt.Print(pointed_file)

	file_format := strings.Split(pointed_file, ".")[1]
	if file_format != "" {
		return file_format
	}
	return ""
}
func FileSorter(path *string) {
	var files []string
	filepath.WalkDir(*path, func(path string, d os.DirEntry, err error) error {
		files = append(files, path)
		return nil
	})

	for i := range files {
		if i == 0 { //wtf
			continue
		}
		file_format := checkFileFormat(files[i])

		switch file_format {
		case "md", "txt", "pdf", "doc", "docx", "xlsx", "xls":
			moveFileToFolder(files[i], "Documents")
		case "mp4", "avi", "mkv", "mov":
			moveFileToFolder(files[i], "Videos")
		case "jpg", "png", "gif", "webp":
			moveFileToFolder(files[i], "Images")
		default:
			moveFileToFolder(files[i], "Others")

		}

	}
}
