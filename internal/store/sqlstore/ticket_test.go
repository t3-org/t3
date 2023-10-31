package sqlstore

//func TestTicketStore_Get(t *testing.T) {
//	defer testbox.Global().TeardownIfPanic()
//
//	s := store()
//	ctx := context.Background()
//
//	var p model.Ticket
//	require.NoError(t, p.Create(&input.CreateTicket{
//		Name: "name-abc",
//		Code: "code-abc",
//	}))
//
//	require.NoError(t, s.Ticket().Create(ctx, &p))
//	defer func() { require.NoError(t, s.Ticket().Delete(ctx, &p)) }()
//
//	var p2 model.Ticket
//	require.NoError(t, p2.Create(&input.CreateTicket{
//		Name: "name-abc2",
//		Code: "code-abc2",
//	}))
//
//	require.NoError(t, s.Ticket().Create(ctx, &p2))
//	defer func() { require.NoError(t, s.Ticket().Delete(ctx, &p2)) }()
//
//	res, err := s.Ticket().Get(ctx, p.ID)
//	require.NoError(t, err)
//	require.Equal(t, &p, res)
//
//	res2, err := s.Ticket().Get(ctx, p2.ID)
//	require.NoError(t, err)
//	require.Equal(t, &p2, res2)
//}
//
//func TestTicketStore_Create(t *testing.T) {
//	defer testbox.Global().TeardownIfPanic()
//
//	s := store()
//	ctx := context.Background()
//
//	var p model.Ticket
//	require.NoError(t, p.Create(&input.CreateTicket{
//		Name: "name-abc",
//		Code: "code-abc",
//	}))
//
//	require.NoError(t, s.Ticket().Create(ctx, &p))
//	defer func() { require.NoError(t, s.Ticket().Delete(ctx, &p)) }()
//
//	count, err := s.Ticket().Count(ctx, "")
//	require.NoError(t, err)
//	require.Equal(t, 1, count)
//
//	var p2 model.Ticket
//	require.NoError(t, p2.Create(&input.CreateTicket{
//		Name: "name-abc2",
//		Code: "code-abc2",
//	}))
//
//	require.NoError(t, s.Ticket().Create(ctx, &p2))
//	defer func() { require.NoError(t, s.Ticket().Delete(ctx, &p2)) }()
//
//	count, err = s.Ticket().Count(ctx, "")
//	require.NoError(t, err)
//	require.Equal(t, 2, count)
//}
