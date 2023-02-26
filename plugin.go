package main

import (
	"com.github.eshluke.helloworld/helper"
	"fmt"
	"net/http"

	"com.github.eshluke.helloworld/gitlab"
	"github.com/mattermost/mattermost-server/plugin/rpcplugin"
)

type Plugin struct{}

func (p *Plugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
	var mr gitlab.MergeRequest
	err := helper.DecodeJSONBody(w, r, &mr)
	if err == nil {
		fmt.Fprintf(w, "received: %s", mr.ObjectKind)
	}
}

// This example demonstrates a plugin that handles HTTP requests which respond by greeting the
// world.
func main() {
	rpcplugin.Main(&Plugin{})
}
