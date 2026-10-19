package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	// Создаётся сокет принимающий сообщения
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Не робит (порт 8080 уже занят):", err)
		return
	}
	defer listener.Close()
	fmt.Println("Сервер запущен на 127.0.0.1:8080. Ждёмссс")

// Цикл для приёма сообщений 
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("Клиент прибыл")

		// Горутина типо нужно для того чтобы запустить несколько терминалов иначе всё залагает
		 go func(c net.Conn) {
			defer c.Close() // То есть говорим функции что нужно закрыть соединение передм самым закрытием функции
			scanner := bufio.NewScanner(c)

			for scanner.Scan() {
				msg := scanner.Text()

				// ВЫДАЧА НА СЕРВЕРЕ
				fmt.Println("[Клиент]:", msg)

				// Отправляем ответ обратно клиенту
				fmt.Fprintf(c, "Сервер: %s\n", msg)
			}
			fmt.Println("Клиент отключился.")
		}(conn)
	}
}
