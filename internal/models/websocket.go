package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type WebSocketMessageType string

const (
	WSMessageJoinRetrospective   WebSocketMessageType = "join_retrospective"
	WSMessageLeaveRetrospective  WebSocketMessageType = "leave_retrospective"
	WSMessageNewItem             WebSocketMessageType = "new_item"
	WSMessageUpdateItem          WebSocketMessageType = "update_item"
	WSMessageDeleteItem          WebSocketMessageType = "delete_item"
	WSMessageVoteItem            WebSocketMessageType = "vote_item"
	WSMessageUnvoteItem          WebSocketMessageType = "unvote_item"
	WSMessageUpdateRetrospective WebSocketMessageType = "update_retrospective"
	WSMessageNewActionItem       WebSocketMessageType = "new_action_item"
	WSMessageUpdateActionItem    WebSocketMessageType = "update_action_item"
	WSMessageUserJoined          WebSocketMessageType = "user_joined"
	WSMessageUserLeft            WebSocketMessageType = "user_left"
	WSMessageError               WebSocketMessageType = "error"
)

type WebSocketMessage struct {
	Type      WebSocketMessageType `json:"type"`
	Data      json.RawMessage      `json:"data,omitempty"`
	Timestamp int64                `json:"timestamp"`
	UserID    *uuid.UUID           `json:"user_id,omitempty"`
}

type JoinRetrospectiveData struct {
	RetrospectiveID uuid.UUID `json:"retrospective_id"`
	UserID          uuid.UUID `json:"user_id"`
}

type NewItemData struct {
	Item RetrospectiveItem `json:"item"`
}

type UpdateItemData struct {
	Item RetrospectiveItem `json:"item"`
}

type DeleteItemData struct {
	ItemID uuid.UUID `json:"item_id"`
}

type VoteItemData struct {
	ItemID uuid.UUID `json:"item_id"`
	UserID uuid.UUID `json:"user_id"`
}

type UserJoinedData struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
}

type UserLeftData struct {
	UserID uuid.UUID `json:"user_id"`
}

type ErrorData struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type NewActionItemData struct {
	ActionItem ActionItem `json:"action_item"`
}
