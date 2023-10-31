package model

// TicketTag is a tag on a ticket.
// @model
type TicketTag struct {
	TicketID int64  `json:"ticket_id"`
	Term     string `json:"term"`
}

func Tags(tags ...*TicketTag) []string {
	res := make([]string, len(tags))
	for i, v := range tags {
		res[i] = v.Term
	}
	return res
}
