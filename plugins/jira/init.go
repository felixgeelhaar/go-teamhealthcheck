package jira

import "github.com/felixgeelhaar/heartbeat/sdk"

func init() {
	sdk.Register(&JiraPlugin{})
}
