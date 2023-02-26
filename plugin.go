package main

import (
	"errors"
	"fmt"
	"net/http"

	"com.github.eshluke.helloworld/gitlab"
	"com.github.eshluke.helloworld/helper"
	"github.com/mattermost/mattermost-server/plugin/rpcplugin"
)

type Plugin struct{}

func (p *Plugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
	var mr gitlab.MergeRequest
	err := helper.DecodeJSONBody(w, r, &mr)
	if err != nil {
		var mr *helper.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	if mr.ObjectKind == "merge_request" {
		fmt.Fprintf(w, "merge_request event received!!!")
	}
}

// This example demonstrates a plugin that handles HTTP requests which respond by greeting the
// world.
func main() {
	rpcplugin.Main(&Plugin{})
}
