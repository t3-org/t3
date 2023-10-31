package app

//func TestAppCore_CreateTicket(t *testing.T) {
//	defer testbox.Global().TeardownIfPanic()
//	_, s, a := setup(t)
//	now := time.Now().UnixMilli()
//
//	_, _ = s, a
//	in := input.CreateTicket{
//		Name: "abc",
//		Code: "abc",
//	}
//
//	var res *model.Ticket
//	s.Ticket().(*mockmodel.MockTicketStore).EXPECT().Create(_any, _any).DoAndReturn(
//		func(_ context.Context, m *model.Ticket) error {
//			res = m
//			return nil
//		},
//	)
//
//	p, err := a.CreateTicket(context.Background(), &in)
//	require.NoError(t, err)
//	require.Equal(t, p, res)
//	require.NotEmpty(t, p.ID)
//	require.GreaterOrEqual(t, p.CreatedAt, now)
//	require.GreaterOrEqual(t, p.UpdatedAt, now)
//	require.Equal(t, p.Name, in.Name)
//	require.Equal(t, p.Code, in.Code)
//}
