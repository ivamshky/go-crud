package model

type User struct {
	Id    int64  `json:"id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int32  `json:"age"`
}

type ListUserParams struct {
	Id     *int64
	Name   *string
	Email  *string
	Limit  *int
	Offset *int
}
