package api

type ErrorMessage struct {
	Error []string `json:"error"`
}

func NewErrorMessage(err string) ErrorMessage {
	return ErrorMessage{Error: []string{err}}
}
