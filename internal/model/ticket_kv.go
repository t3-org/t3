package model

type TicketKV struct {
	TicketID int64  `json:"ticket_id"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}
