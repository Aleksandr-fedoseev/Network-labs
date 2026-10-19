package main

import (
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"
)

func handleFileProtocol(rawURL string) ([]byte, error) {
	// 1. Разбираем строку URI
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	// 2. Проверяем что схема именно file://
	if parsed.Scheme != "file" {	
		return nil, fmt.Errorf("ожидается схема file://, получено: %s", parsed.Scheme)
	}

	filePath := parsed.Path

	// 3. Исправление пути для Windows из за того что начинается со слэша
	if runtime.GOOS == "windows" && strings.HasPrefix(filePath, "/") && len(filePath) > 2 && filePath[2] == ':' {
		filePath = filePath[1:]
	} // Проверяем что в начале стоит слэш, 3 символ это : и отрезаем первый слэш

	// 4. Читаем байты файла
	return os.ReadFile(filePath)
}

func main() {
	// Создаем тестовый файл с данными
	testFileName := "test_lab.txt"
	os.WriteFile(testFileName, []byte("Привет! Это данные, прочитанные по протоколу file://"), 0644)

	// Собираем валидный URI так как в сетях не используются \\
	pwd, _ := os.Getwd()
	testURL := fmt.Sprintf("file:///%s/%s", strings.ReplaceAll(pwd, "\\", "/"), testFileName)
	fmt.Printf("Тестируем URI: %s\n", testURL)

	// Читаем файл по протоколу file://
	data, err := handleFileProtocol(testURL)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("Успешно прочитано содержимое файла:\n\"%s\"\n", string(data))
}
