package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// scanPort проверяет конкретный порт на конкретном IP
func scanPort(host string, port int, timeout time.Duration, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return 
	}
	conn.Close()

	results <- fmt.Sprintf("%s:%d", host, port)
}

func main() {
	subnet := "10.26.93."
	startHost := 1
	endHost := 254

	// Список популярных портов для проверки каждого хоста
	ports := []int{21, 22, 80, 135, 139, 443, 445, 3389, 8080}
	timeout := 400 * time.Millisecond

	var wg sync.WaitGroup
	results := make(chan string, 1000)

	// Ограничиваем число одновременных соединений
	semaphore := make(chan struct{}, 200)

	fmt.Printf("Сканирование сети %s%d-%d по портам %v...\n", subnet, startHost, endHost, ports)

	// Двойной цикл: перебираем IP-адреса и порты
	for h := startHost; h <= endHost; h++ {
		host := fmt.Sprintf("%s%d", subnet, h)

		for _, port := range ports {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(ip string, p int) {
				defer func() { <-semaphore }()
				scanPort(ip, p, timeout, &wg, results)
			}(host, port)
		}
	}

	// Закрываем канал после завершения всех горутин
	go func() {
		wg.Wait()
		close(results)
	}()

	openServices := []string{}
	// Читаем все результаты из канала
	for res := range results {
		openServices = append(openServices, res)
	}

	if len(openServices) == 0 {
		fmt.Println("Активных устройств и открытых портов не найдено")
		return
	}

	fmt.Println("Найдены активные адреса и открытые порты:")
	for _, s := range openServices {
		fmt.Printf("%s/tcp открыт\n", s)
	}
}
