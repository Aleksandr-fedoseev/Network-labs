package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Создаю сокет отдающий сообщения
	conn, _ := net.Dial("tcp", "127.0.0.1:8080")
	defer conn.Close()

	fmt.Println("Подключилось вводите сообщение")

	inputScanner := bufio.NewScanner(os.Stdin) // Читает текст с клаваиатуры 
	serverScanner := bufio.NewScanner(conn)	 // Читает текст с приёмника, в данном случае с сокета 

	// Цикл для диалога чата
	for {
		fmt.Print("Вы: ")
		// Если нет входящего сообщения то выходим
		if !inputScanner.Scan() {
			break
		}

		text := inputScanner.Text()

		// Fprintln - типо отправляет сообщение в указанную сеть (куда отправить, текст)
		fmt.Fprintln(conn, text)

		// Получаем ответ от сервера
		if serverScanner.Scan() {
			fmt.Println(serverScanner.Text())
		}
	}
}
