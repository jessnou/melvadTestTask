package models

type RedisRequest struct {
	Key   string `json:"key" binding:"required"`
	Value int    `json:"value" binding:"required"`
}

type SignRequest struct {
	Text string `json:"text" binding:"required"`
	Key  string `json:"key" binding:"required"`
}

type UserRequest struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required"`
}

type UserResponse struct {
	ID int `json:"id"`
}
