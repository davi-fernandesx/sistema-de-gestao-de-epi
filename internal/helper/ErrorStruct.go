package helper


type HTTPError struct {
    Message string `json:"message" `
    Code    int    `json:"code,omitempty"`
}