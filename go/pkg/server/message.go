package server

import "github.com/gorilla/websocket"

type Message struct {
    Type    uint
    ID      uint
    Message string
}

type ChatMessage struct {
    ChannelName      string `json:"channel_name"`
    ChannelUserCount int    `json:"channel_user_count"`
    From             uint   `json:"from"`
    Msg string `json:"msg"`
}

func (m *Message) FromMessage(message string) *Message {
    return &Message {
        websocket.TextMessage,
        m.ID,
        message,
    }
}

func NewMessage(id uint, message string) *Message {
    return &Message {
        Type:    websocket.TextMessage,
        ID:      id,
        Message: message,
    }
}

func CloseMessage(id uint) *Message {
    return &Message {
        websocket.CloseMessage,
        id,
        "",
    }
}
