package retro

import "github.com/felixgeelhaar/heartbeat/sdk"

func init() {
	sdk.Register(&RetroPlugin{})
}
