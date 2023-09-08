package event

// StockRequest ...
type StockRequest struct {
	RoomId uint   `json:"RoomId"`
	Code   string `json:"Code"`
}
