package fileworker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMoveFileToFolder(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "fileworker_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Создаем тестовый файл
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		folderName string
		setup      func() // подготовка перед тестом
		wantErr    bool
	}{
		{
			name:       "Перемещение в новую папку",
			folderName: "Documents",
			setup:      func() {},
			wantErr:    false,
		},
		{
			name:       "Перемещение в существующую папку",
			folderName: "ExistingFolder",
			setup: func() {
				existingDir := filepath.Join(tempDir, "ExistingFolder")
				os.MkdirAll(existingDir, 0755)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новый файл для каждого теста
			testFile := filepath.Join(tempDir, "test_"+tt.name+".txt")
			err := os.WriteFile(testFile, []byte("test content"), 0644)
			if err != nil {
				t.Fatal(err)
			}

			tt.setup()

			err = moveFileToFolder(testFile, tt.folderName)
			if (err != nil) != tt.wantErr {
				t.Errorf("moveFileToFolder() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Проверяем, что файл перемещен
			destPath := filepath.Join(filepath.Dir(testFile), tt.folderName, filepath.Base(testFile))
			if _, err := os.Stat(destPath); os.IsNotExist(err) {
				t.Errorf("Файл не был перемещен в %s", destPath)
			}
		})
	}
}

func TestCheckFileFormat(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "Обычное расширение",
			filename: "/home/user/document.pdf",
			expected: "pdf",
		},
		{
			name:     "Файл с несколькими точками",
			filename: "/home/user/archive.tar.gz",
			expected: "gz",
		},
		{
			name:     "Файл без расширения",
			filename: "/home/user/README",
			expected: "",
		},
		{
			name:     "Путь с пробелами",
			filename: "/home/user/my file.txt",
			expected: "txt",
		},
		{
			name:     "Только имя файла",
			filename: "image.jpg",
			expected: "jpg",
		},
		{
			name:     "Скрытый файл",
			filename: "/home/user/.bashrc",
			expected: "bashrc",
		},
		{
			name:     "Пустая строка",
			filename: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkFileFormat(tt.filename)
			if result != tt.expected {
				t.Errorf("checkFileFormat(%s) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestFileSorter(t *testing.T) {
	// Создаем временную директорию
	tempDir, err := os.MkdirTemp("", "filesorter_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Создаем тестовые файлы разных типов
	testFiles := []struct {
		name     string
		content  string
		fileType string
	}{
		{"document.txt", "text content", "Documents"},
		{"presentation.pdf", "pdf content", "Documents"},
		{"video.mp4", "video content", "Videos"},
		{"image.jpg", "image content", "Images"},
		{"unknown.xyz", "unknown content", "Others"},
		{"no_extension", "no extension", "Others"},
	}

	for _, tf := range testFiles {
		filePath := filepath.Join(tempDir, tf.name)
		err := os.WriteFile(filePath, []byte(tf.content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Запускаем сортировку
	FileSorter(&tempDir)

	// Проверяем результаты
	checks := []struct {
		fileName     string
		expectedDir  string
		shouldExists bool
	}{
		{"document.txt", "Documents", true},
		{"presentation.pdf", "Documents", true},
		{"video.mp4", "Videos", true},
		{"image.jpg", "Images", true},
		{"unknown.xyz", "Others", true},
		{"no_extension", "Others", true},
	}

	for _, check := range checks {
		expectedPath := filepath.Join(tempDir, check.expectedDir, check.fileName)
		_, err := os.Stat(expectedPath)
		if check.shouldExists && os.IsNotExist(err) {
			t.Errorf("Файл %s должен быть в папке %s, но не найден", check.fileName, check.expectedDir)
		}
		if !check.shouldExists && !os.IsNotExist(err) {
			t.Errorf("Файл %s не должен существовать в %s", check.fileName, expectedPath)
		}
	}
}

func TestFileSorterWithExistingFolders(t *testing.T) {
	// Создаем временную директорию
	tempDir, err := os.MkdirTemp("", "filesorter_existing_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Предварительно создаем папки
	folders := []string{"Documents", "Videos", "Images", "Others"}
	for _, folder := range folders {
		err := os.MkdirAll(filepath.Join(tempDir, folder), 0755)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Создаем тестовые файлы
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Запускаем сортировку (не должно быть ошибок)
	FileSorter(&tempDir)

	// Проверяем, что файл перемещен
	expectedPath := filepath.Join(tempDir, "Documents", "test.txt")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Файл не был перемещен в существующую папку Documents")
	}
}

// Тест на конкуррентность
func TestFileSorterConcurrency(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "concurrency_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Создаем много файлов
	numFiles := 100
	for i := 0; i < numFiles; i++ {
		fileName := filepath.Join(tempDir, "file_%d.txt")
		err := os.WriteFile(fileName, []byte("content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Запускаем сортировку
	FileSorter(&tempDir)

	// Проверяем, что все файлы перемещены
	destDir := filepath.Join(tempDir, "Documents")
	files, err := os.ReadDir(destDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != numFiles {
		t.Errorf("Ожидалось %d файлов в папке Documents, получено %d", numFiles, len(files))
	}
}

// Бенчмарк тесты
func BenchmarkCheckFileFormat(b *testing.B) {
	filename := "/very/long/path/to/some/document/file.pdf"

	for i := 0; i < b.N; i++ {
		checkFileFormat(filename)
	}
}

func BenchmarkMoveFileToFolder(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "benchmark_test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testFile := filepath.Join(tempDir, "test.txt")
		os.WriteFile(testFile, []byte("content"), 0644)
		moveFileToFolder(testFile, "Documents")
	}
}
