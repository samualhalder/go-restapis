package types

type Student struct {
	Id    int    `json:"id"`
	Name  string `validate:"required" json:"name"`
	Age   int    `validate:"required" json:"age"`
	Email string `validate:"required" json:"email"`
}
