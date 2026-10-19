package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// 1. Запускаем DNS-сервер на порту 5354 (UDP)
	serverConn, err := net.ListenPacket("udp", "127.0.0.1:5354")
	if err != nil {
		panic(err)
	}
	defer serverConn.Close()

	fmt.Println("DNS-сервер успешно запущен на 127.0.0.1:5354 (UDP)")

	// Сервер работает в фоновом режиме (горутине) 
	go func() {
		buf := make([]byte, 512) // 512 байт по спецификации RFC 1035  - размер DNS пакета
		for {
			n, clientAddr, err := serverConn.ReadFrom(buf) // n - количество байт и clientAddr - адрес клиента
			if err != nil {
				return
			}

			// Формируем DNS-ответ
			resp := make([]byte, n+16)
			copy(resp, buf[:n])

			// По стандарту DNS имеет 12 ячеек памяти, которые нужны для того чтобы сервер мог получить данные
			resp[2] = 0x81 // Это ответ от сервера 
			resp[3] = 0x80 // Это статус выполнения 
			resp[6] = 0x00
			resp[7] = 0x01 // Сколько ответов лежит в ячейке памяти - это IP адрес

			// Сам ответ из 16 ячеек
			offset := n
			resp[offset] = 0xc0 // Имя домена
			resp[offset+1] = 0x0c
			resp[offset+2] = 0x00 // TYPE = 1 (A-запись, IPv4)
			resp[offset+3] = 0x01
			resp[offset+4] = 0x00 // CLASS = 1 (IN, Internet)
			resp[offset+5] = 0x01
			resp[offset+6] = 0x00 // TTL = 60 секунд (0x0000003C) или кэш хранится 60 секунд
			resp[offset+7] = 0x00
			resp[offset+8] = 0x00
			resp[offset+9] = 0x3c
			resp[offset+10] = 0x00 // RDLENGTH = 4 байта или далее идут 4 байта
			resp[offset+11] = 0x04
			resp[offset+12] = 1 // IP-адрес 1.2.3.4
			resp[offset+13] = 2
			resp[offset+14] = 3
			resp[offset+15] = 4

			// Отправляем ответ обратно
			serverConn.WriteTo(resp, clientAddr)
			fmt.Printf("DNS-сервер: ответил клиенту %s -> IP 1.2.3.4\n", clientAddr.String())
		}
	}()

	// Это время на то чтобы сервер открыл клиентский udp сокет
	time.Sleep(100 * time.Millisecond)
	fmt.Println("DNS-клиент: отправляем запрос для домена test.local...")

	clientConn, err := net.Dial("udp", "127.0.0.1:5354")
	if err != nil {
		panic(err)
	}
	defer clientConn.Close()

	// Бинарный DNS-запрос
	query := []byte{
		// ЧАСТЬ 1: Заголовок запроса (12 байт)
		0x12, 0x34, // ID
		0x01, 0x00, // Флаги (RD=1)
		0x00, 0x01, // QDCOUNT = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// ЧАСТЬ 2: Секция Question (Сам вопрос)
		0x04, 't', 'e', 's', 't', 0x05, 'l', 'o', 'c', 'a', 'l', 0x00, // 4"test" 5"local" 0
		0x00, 0x01, // QTYPE = A
		0x00, 0x01, // QCLASS = IN
	}

	clientConn.Write(query)

	// Читаем ответ
	clientConn.SetReadDeadline(time.Now().Add(2 * time.Second))
	respBuf := make([]byte, 512)
	n, err := clientConn.Read(respBuf)
	if err != nil {
		fmt.Printf("DNS-клиент: ошибка ответа: %v\n", err)
		return
	}

	if n >= 4 {
		ip := respBuf[n-4 : n]
		fmt.Printf("DNS-клиент: успешно получен IP-адрес: %d.%d.%d.%d\n", ip[0], ip[1], ip[2], ip[3])
	}
}
