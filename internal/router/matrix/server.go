package matrix

import (
	"context"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"space.org/space/internal/app"
)

type Server struct {
	cli    *mautrix.Client
	app    app.App
	doneCh chan error
}

func NewServer(cli *mautrix.Client, app app.App) *Server {
	return &Server{cli: cli, app: app, doneCh: make(chan error)}
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

func (s *Server) Shutdown(ctx context.Context) error {
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
		body := evt.Content.AsMessage().Body
		hlog.Info("received new message",
			hlog.String("sender", evt.Sender.String()),
			hlog.String("type", evt.Type.String()),
			hlog.String("id", evt.ID.String()),
			hlog.String("body", body),
		)

		relatedTo := evt.Content.AsMessage().GetRelatesTo()
		if relatedTo.GetThreadParent() != "" && body == "!party" {
			content := event.MessageEventContent{
				MsgType:   event.MsgText,
				Body:      ":))))",
				RelatesTo: relatedTo,
			}
			_, err := s.cli.SendMessageEvent(evt.RoomID, event.EventMessage, &content)

			if err != nil {
				hlog.Error("can not send message to matrix channel",
					hlog.String("room_id", string(evt.RoomID)),
				)
			}
		}
	})
}

var _ hexa.Runnable = &Server{}
var _ hexa.Shutdownable = &Server{}
