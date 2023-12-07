package matrixch

import (
	"fmt"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
	"t3.org/t3/internal/config"
	"t3.org/t3/internal/service/channel"
	"t3.org/t3/pkg/md"
)

type HomeOpts struct {
	// Set keyPrefix with "." (config.InternalLabelKeyPrefix) to set that label as an internal label.
	KeyPrefix     string
	OkEmoji       string
	CommandPrefix string

	Client   *mautrix.Client
	KVStore  channel.KVStore
	MarkDown *md.Markdown
	UI       config.UI
}

func (o *HomeOpts) Key(roomID id.RoomID) string {
	return fmt.Sprintf("%s:%s", o.KeyPrefix, roomID)
}
