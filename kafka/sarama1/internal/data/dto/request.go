package dto

type MessageDemo struct {
	Msg string `json:"msg"`
}

type QueMessage struct {
	Msg  string
	Time int64
}
