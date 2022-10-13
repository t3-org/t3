package mockmodel

import (
	"github.com/golang/mock/gomock"
)

func NewStore(t gomock.TestReporter) *MockStore {
	ctl := gomock.NewController(t)
	system := NewMockSystemStore(ctl)
	planet := NewMockPlanetStore(ctl)

	s := NewMockStore(ctl)
	s.EXPECT().System().Return(system).AnyTimes()
	s.EXPECT().Planet().Return(planet).AnyTimes()

	// Set your global expectations on the store mock here.

	return s
}
