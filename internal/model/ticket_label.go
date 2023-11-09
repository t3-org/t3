package model

// TicketLabel is a label on a ticket.
// @model
type TicketLabel struct {
	TicketID string `json:"ticket_id"`
	Key      string `json:"key"`
	Val      string `json:"val"`
}

func LabelsFromMap(ticketID string, m map[string]string) []*TicketLabel {
	l := make([]*TicketLabel, len(m))
	var i int
	for k, v := range m {
		l[i] = &TicketLabel{
			TicketID: ticketID,
			Key:      k,
			Val:      v,
		}
		i++
	}
	return l
}

func LabelsMap(labels ...*TicketLabel) map[string]string {
	res := make(map[string]string)
	for _, v := range labels {
		res[v.Key] = v.Val
	}
	return res
}
