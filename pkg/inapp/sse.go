package inapp

import (
	"fmt"
	"net/http"
	"sync"
)

// SSEClient đại diện cho một kết nối SSE tới client.
type SSEClient struct {
	Writer  http.ResponseWriter
	Flusher http.Flusher
	Done    chan struct{}
}

// SSEManager quản lý tất cả các kết nối SSE theo userID.
type SSEManager struct {
	clients map[string]*SSEClient
	mutex   sync.Mutex
}

// Manager là instance toàn cục của SSEManager.
var Manager = &SSEManager{
	clients: make(map[string]*SSEClient),
}

// RegisterClient đăng ký một kết nối SSE mới cho user.
func (m *SSEManager) RegisterClient(userID string, w http.ResponseWriter) {
	// Kiểm tra xem ResponseWriter có hỗ trợ Flusher hay không.
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Thiết lập header SSE.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Tạo client mới.
	client := &SSEClient{
		Writer:  w,
		Flusher: flusher,
		Done:    make(chan struct{}),
	}

	// Đăng ký client vào manager.
	m.mutex.Lock()
	m.clients[userID] = client
	m.mutex.Unlock()

	// Giữ kết nối mở cho đến khi client ngắt kết nối.
	<-client.Done

	// Khi client ngắt, xóa khỏi danh sách.
	m.mutex.Lock()
	delete(m.clients, userID)
	m.mutex.Unlock()
}

// SendNotification gửi thông báo tới user thông qua SSE.
func (m *SSEManager) SendNotification(userID string, message string) error {
	m.mutex.Lock()
	client, exists := m.clients[userID]
	m.mutex.Unlock()

	if !exists {
		return fmt.Errorf("user %s không có kết nối SSE", userID)
	}

	// Gửi message theo định dạng SSE.
	fmt.Fprintf(client.Writer, "data: %s\n\n", message)
	client.Flusher.Flush()
	return nil
}

// SSEHandler là HTTP handler để client kết nối tới SSE.
func SSEHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy userID từ query parameter.
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "Thiếu userId", http.StatusBadRequest)
		return
	}

	// Đăng ký kết nối SSE cho user.
	Manager.RegisterClient(userID, w)
}
