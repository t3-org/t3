package mockmodel

import (
	"github.com/golang/mock/gomock"
)

func NewStore(t gomock.TestReporter) *MockStore {
	ctl := gomock.NewController(t)
	system := NewMockSystemStore(ctl)
	ticket := NewMockTicketStore(ctl)
	ticketKV := NewMockTicketKVStore(ctl)

	s := NewMockStore(ctl)
	s.EXPECT().System().Return(system).AnyTimes()
	s.EXPECT().Ticket().Return(ticket).AnyTimes()
	s.EXPECT().TicketKV().Return(ticketKV).AnyTimes()

	// Set your global expectations on the store mock here.

	return s
}
