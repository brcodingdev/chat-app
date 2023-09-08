package request

// ChatRoomMessagesRequest ...
type ChatRoomMessagesRequest struct {
	RoomId uint `json:"RoomId"`
}

// ChatRoomCreateRequest ...
type ChatRoomCreateRequest struct {
	Name string `json:"Name"`
}
