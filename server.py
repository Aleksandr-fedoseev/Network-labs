import socket
# Создание сокета
server = socket.socket()
# Привязка к адресу локалхост
server.bind(('127.0.0.1', 8080))
# Режим ожидания 
server.listen()
print("Сервер включен и ждет подключения на порту 8080...")
conn, addr = server.accept()
print(f"Кто-то подключился с адреса: {addr}")
request = conn.recv(1024).decode('utf-8')
print("=== Запрос от браузера ===")
print(request)
body = "<h1>200 OK: Мой первый сервер на Python работает!</h1>"
response = (
    "HTTP/1.1 200 OK\r\n"
    "Content-Type: text/html; charset=utf-8\r\n"
    f"Content-Length: {len(body.encode('utf-8'))}\r\n"
    "Connection: close\r\n\r\n"
    f"{body}"
)
conn.sendall(response.endcode('utf-8'))
conn.close
print("Ответ успешно отправлен, соединие закрыто")