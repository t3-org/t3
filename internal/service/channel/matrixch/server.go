package matrixch

import (
	"context"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

type Server struct {
	// TODO: move this fn to the registry package after
	//  changing the sr package to include BootFn definition in itself.
	bootFn func() error
	cli    *mautrix.Client
	router *Router
	doneCh chan error
}

func NewServer(cli *mautrix.Client, r *Router, bootFn func() error) *Server {
	return &Server{bootFn: bootFn, cli: cli, router: r, doneCh: make(chan error)}
}

func (s *Server) Boot() error {
	return tracer.Trace(s.bootFn())
}

func (s *Server) Run() (done <-chan error, err error) {
	cli := s.cli
	s.registerInviteHandler()
	if err := s.syncToLastState(); err != nil {
		return nil, tracer.Trace(err)
	}
	s.registerMsgHandler()

	go func() {
		cli.SyncWithContext(context.Background())
	}()

	return s.doneCh, nil
}

func (s *Server) doShutdown(_ context.Context) error {
	s.cli.StopSync()
	close(s.doneCh)
	return nil
}

func (s *Server) syncToLastState() error {
	cli := s.cli
	resp, err := cli.SyncRequest(
		30000,
		cli.Store.LoadNextBatch(cli.UserID),
		"",
		true,
		cli.SyncPresence,
		context.Background(),
	)
	if err != nil {
		return tracer.Trace(err)
	}

	cli.Store.SaveNextBatch(cli.UserID, resp.NextBatch)
	return nil
}

func (s *Server) registerInviteHandler() {
	syncer := s.cli.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.StateMember, func(source mautrix.EventSource, evt *event.Event) {
		if evt.GetStateKey() == s.cli.UserID.String() && evt.Content.AsMember().Membership == event.MembershipInvite {
			_, err := s.cli.JoinRoomByID(evt.RoomID)
			l := hlog.With(
				hlog.String("room_id", evt.RoomID.String()),
				hlog.String("inviter", evt.Sender.String()),
			)
			if err != nil {
				l.Error("can not join room", hlog.Err(err))
			} else {
				l.Message("joined room")
			}
		}
	})
}

func (s *Server) registerMsgHandler() {
	syncer := s.cli.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		if err := s.router.Route(context.Background(), source, evt); err != nil {
			hlog.Error("can not handle matrix message",
				hlog.Err(err),
				hlog.String("room_id", string(evt.RoomID)),
			)
		}
	})
}

var _ hexa.Bootable = &Server{}
var _ hexa.Runnable = &Server{}
