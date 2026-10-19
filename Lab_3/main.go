package main
import (
	"fmt"
	"net"
	"sync"
	"time"
)
func scanPort(host string, port int, timeout time.Duration, wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()

	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return // порт закрыт или недоступен
	}
	conn.Close()
	results <- port
}
func main() {
	host := "127.0.0.1" // локальный адрес
	startPort := 1
	endPort := 1000
	timeout := 500 * time.Millisecond
	var wg sync.WaitGroup
	results := make(chan int, endPort-startPort+1)
	// ограничиваем число одновременных горутин чтобы не было множества одновременных сетевых соединений 
	semaphore := make(chan struct{}, 200)
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(p int) {
			defer func() { <-semaphore }()
			scanPort(host, p, timeout, &wg, results)
		}(port)
	}
	// закрываем канал после завершения всех горутин
	go func() {
		wg.Wait()
		close(results)
	}()
	fmt.Printf("Сканирование %s, порты %d-%d...\n", host, startPort, endPort)
	openPorts := []int{}
	// аналог foreach
	for port := range results {
		// Добавить в массив openPorts значения port
		openPorts = append(openPorts, port)
	}
	if len(openPorts) == 0 {
		fmt.Println("Открытых портов нету")
		return
	}
	fmt.Println("Порты открыты:")
	for _, p := range openPorts {
		fmt.Printf("  %d/tcp открыт\n", p)
	}
}