package model

import "sync"

type Cache struct {
	mu          sync.Mutex
	mapMessages map[string][]Message
}

func NewCache() *Cache {
	return &Cache{
		mapMessages: make(map[string][]Message),
	}
}

// AddMessages добавляет Message в кэш, где ключом является имя файла(FileID)
func (c *Cache) AddMessages(msg Message) {
	c.mu.Lock()
	c.mapMessages[msg.FileID] = append(c.mapMessages[msg.FileID], msg)
	c.mu.Unlock()
}

// GetMessages извлекает все сообщения из кэша и очищает кэш.
func (c *Cache) GetMessages() map[string][]Message {
	c.mu.Lock()
	defer c.mu.Unlock()
	messages := c.mapMessages
	c.mapMessages = make(map[string][]Message)
	return messages
}
