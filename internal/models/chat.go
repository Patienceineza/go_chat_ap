package models

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
    "time"
)

type User struct {
    ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
    Username  string    `gorm:"unique;not null" json:"username"`
    Password  string    `gorm:"not null" json:"password"`
    Role      string    `gorm:"not null;default:'user'" json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Message struct {
    ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
    SenderID   uuid.UUID `gorm:"type:uuid;not null" json:"sender_id"`
    ReceiverID uuid.UUID `gorm:"type:uuid;not null" json:"receiver_id"`
    Content    string    `gorm:"type:text;not null" json:"content"`
    ReplyToID  *uuid.UUID `gorm:"type:uuid;" json:"reply_to_id"`
    Timestamp  int64     `gorm:"autoCreateTime" json:"timestamp"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
    
    Sender   User `gorm:"foreignKey:SenderID;references:ID"`
    Receiver User `gorm:"foreignKey:ReceiverID;references:ID"`
}
type ChatRoom struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type RoomMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	RoomID    uuid.UUID `json:"room_id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (ChatRoom) TableName() string {
	return "chat_rooms"
}

func (RoomMessage) TableName() string {
	return "room_messages"
}
