package model

type Person struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" binding:"required" db:"username"`
	Password string `json:"password" binding:"required" db:"password"`
	Role     string `json:"role" db:"role"`
}

type PersonWithChatRoomRole struct {
	Person   *Person
	RoomRole string `db:"role"`
}
