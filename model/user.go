package model

type User struct {
	Id    int64
	Name  string
	Email string
	Age   int32
}

type ListUserParams struct {
	Id    *int64
	Name  *string
	Email *string
}
