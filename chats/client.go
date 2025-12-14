package chats

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	mu sync.Mutex 
}

func (c *Client) SendMessage(msg Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Conn.WriteJSON(msg)
}

// ------- Hub -------
type Hub struct {
	clients map[string]*Client
	mu sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

func (h *Hub) AddClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c.ID] = c
	log.Printf("Пользователь подключен: %s", c.ID)
}

func (h *Hub) RemoveClient(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if c, ok := h.clients[id]; ok {
		c.Conn.Close()
		delete(h.clients, id)
		log.Printf("Пользователь отключен: %s", id)
	}
}

// SendPrivateMessage ищет получателя и пишет ему напрямую
func (h *Hub) SendPrivateMessage(msg Message) {
	// Берем RLock (блокировка на чтение), чтобы найти клиента
	h.mu.RLock()
	recipient, ok := h.clients[msg.To]
	h.mu.RUnlock() 

	if ok {
		err := recipient.SendMessage(msg)
		if err != nil {
			log.Printf("Ошибка отправки пользователю %s: %v", msg.To, err)
			// Если ошибка записи, удаляем клиента
			h.RemoveClient(recipient.ID)
		}
	} else {
		log.Printf("Пользователь %s не найден", msg.To)
	}
}