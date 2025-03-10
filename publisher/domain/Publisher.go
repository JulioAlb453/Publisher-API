package domain

type Event struct {
    ID      string `json:"ID"`
    Payload string `json:"payload"`
    Title   string `json:"Title"`
}