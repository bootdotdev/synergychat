package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Message struct {
	AuthorUsername string
	Text           string
}

type DB interface {
	getMessages() ([]Message, error)
	createMessage(authorUsername, text string) error
	init() error
}

type Disk struct {
	apiDBFilePath string
	mu            *sync.Mutex
}

func (d *Disk) getMessages() ([]Message, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	dat, err := os.ReadFile(d.apiDBFilePath)
	if err != nil {
		return nil, err
	}
	messages := []Message{}
	err = json.Unmarshal(dat, &messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (d *Disk) createMessage(authorUsername, text string) error {
	messages, err := d.getMessages()
	if err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	messages = append(messages, Message{
		AuthorUsername: authorUsername,
		Text:           text,
	})
	dat, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(d.apiDBFilePath, dat, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (d *Disk) init() error {
	d.mu = &sync.Mutex{}
	// make all the directories in the path,
	// not including the file itself
	dir, _ := filepath.Split(d.apiDBFilePath)
	if dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	_, err := os.Stat(d.apiDBFilePath)
	if err != nil && os.IsNotExist(err) {
		return os.WriteFile(d.apiDBFilePath, []byte("[]"), 0644)
	}
	return err
}

type Memory struct {
	mu       *sync.Mutex
	messages []Message
}

func (m *Memory) getMessages() ([]Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return deepCopy(m.messages), nil
}

func (m *Memory) createMessage(authorUsername, text string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = append(m.messages, Message{
		AuthorUsername: authorUsername,
		Text:           text,
	})
	return nil
}

func (m *Memory) init() error {
	m.mu = &sync.Mutex{}
	m.messages = []Message{}
	return nil
}

func deepCopy(messages []Message) []Message {
	result := make([]Message, len(messages))
	copy(result, messages)
	return result
}
