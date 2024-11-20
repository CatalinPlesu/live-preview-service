package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ChannelID   uuid.UUID  `json:"channel_id"`   
	ParentID    *uuid.UUID `json:"parent_id,omitempty"` 
	UserID      uuid.UUID  `json:"user_id"`      
	MessageText string     `json:"message_text"`  
	CreatedAt   *time.Time `json:"created_at,omitempty"` 
}


type MessageMin struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"` 
	UserID      uuid.UUID  `json:"user_id"`      
	MessageText string     `json:"message_text"`  
}
