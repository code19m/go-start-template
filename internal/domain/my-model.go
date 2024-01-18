package domain

type MyModel struct {
	Id   int32
	Name string
	Age  int32
}

type CreateMyModelParams struct {
	Name string
	Age  int32
}
