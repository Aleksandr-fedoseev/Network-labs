package main

import (
	"fmt"
	"net"
)

func main() {
	// 1. Открываем TCP-сокет на локальном порту 8080
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close() // гарантирует что сетевой сокет закроется и освободит порт при завершении

	fmt.Println("HTTP-сервер запущен на http://127.0.0.1:8080")

	for {
		// 2. Принимаем входящее подключение
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// 3. Считываем HTTP-запрос
		buf := make([]byte, 1024) // Буфер для хранения байтов которые присылает браузер
		n, _ := conn.Read(buf)
		fmt.Printf("\nПолучен HTTP-запрос:\n%s\n", string(buf[:n]))

		// 4. Формируем тело и HTTP-заголовки ответа
		body := "<h1>200 OK: Ответ от Go HTTP-сервера</h1>"
		response := fmt.Sprintf(
			"HTTP/1.1 200 OK\r\n"+
				"Content-Type: text/html; charset=utf-8\r\n"+
				"Content-Length: %d\r\n"+
				"Connection: close\r\n\r\n"+
				"%s",
			len([]byte(body)), body, 
		)

		// 5. Отправляем ответ и закрываем сокет
		conn.Write([]byte(response))
		conn.Close()
	}
}
