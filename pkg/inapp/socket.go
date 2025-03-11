package inapp

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// InAppNotification định nghĩa cấu trúc thông báo in-app.
type InAppNotification struct {
	ID      string `json:"id"`
	UserID  string `json:"userId"`
	Message string `json:"message"`
}

// ClientManager quản lý kết nối WebSocket.
type ClientManager struct {
	clients map[string]*websocket.Conn // key: userId
	lock    sync.Mutex
}

var manager = ClientManager{
	clients: make(map[string]*websocket.Conn),
}

// RegisterClient đăng ký kết nối WebSocket cho user.
func RegisterClient(userID string, conn *websocket.Conn) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.clients[userID] = conn
}

// SendInApp gửi thông báo in-app tới client thông qua WebSocket.
func SendInApp(notification InAppNotification) error {
	manager.lock.Lock()
	conn, exists := manager.clients[notification.UserID]
	manager.lock.Unlock()

	if !exists {
		return nil // Nếu client chưa kết nối, bạn có thể lưu DB để gửi khi online.
	}

	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Lỗi gửi in-app notification: %v", err)
		return err
	}
	log.Printf("In-app notification gửi đến user %s thành công.", notification.UserID)
	return nil
}

// Ví dụ HTTP handler để upgrade kết nối WebSocket.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	RegisterClient(userID, conn)
}
