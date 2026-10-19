package main

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"
)

const myIP = "10.37.2.138"    // Мой айпи
const SIO_RCVALL = 0x98000001 // Для винды чтобы пакеты принимать

func main() {
	// "ip4:0" — IPv4, протокол 0 = все пакеты; IPAddr — IP карты
	conn, err := net.ListenIP("ip4:0", &net.IPAddr{IP: net.ParseIP(myIP)})
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}
	defer conn.Close()

	// Включаем приём всех пакетов
	rawConn, _ := conn.SyscallConn() //  доступ к сокету 
	on := uint32(1) // Запуск 
	var ret uint32
	rawConn.Control(func(fd uintptr) { // По сути контрол вызывает функцию и передаёт номер сокета
		err = syscall.WSAIoctl(syscall.Handle(fd), SIO_RCVALL,
			(*byte)(unsafe.Pointer(&on)), 4, nil, 0, &ret, nil, 0) // включается приём пакетов
	})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	buf := make([]byte, 65535)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		var proto string
		switch buf[9] { // 9 номер протокола в RFC
		case 6:
			proto = "TCP"
		case 17:
			proto = "UDP"
		default:
			continue 
		}

		h := int(buf[0]&0x0F) * 4               // длина IP-заголовка
		src := int(buf[h])<<8 | int(buf[h+1])   // порт отправителя
		dst := int(buf[h+2])<<8 | int(buf[h+3]) // порт получателя

		srcIP := net.IP(buf[12:16]) // IP отправителя (байты 12-15)
		dstIP := net.IP(buf[16:20]) // IP получателя (байты 16-19)

		fmt.Printf("%s  %s:%d -> %s:%d  длина: %d байт\n", proto, srcIP, src, dstIP, dst, n)
	}
}
