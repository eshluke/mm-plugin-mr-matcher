package main

import (
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/plugin/rpcplugin"
)

type Plugin struct{}

func (p *Plugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

// This example demonstrates a plugin that handles HTTP requests which respond by greeting the
// world.
func main() {
	rpcplugin.Main(&Plugin{})
}
