package retro

import "github.com/felixgeelhaar/go-teamhealthcheck/sdk"

func init() {
	sdk.Register(&RetroPlugin{})
}
