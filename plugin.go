package main

import (
	"errors"
	"fmt"
	"net/http"

	"com.github.eshluke.helloworld/gitlab"
	"com.github.eshluke.helloworld/helper"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type Plugin struct {
	api plugin.API
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path; path {
	case "/gitlab/mr":
		p.handleMergeRequest(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) handleMergeRequest(w http.ResponseWriter, r *http.Request) {
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
	if mr.ObjectKind == "merge_request" &&
		(mr.ObjectAttributes.Action == "open" || mr.ObjectAttributes.Action == "reopen") &&
		len(mr.Reviewers) == 0 {
		fmt.Fprintf(w, "TODO: handle merge request event...")
	}
}

func (p *Plugin) createPost(mr *gitlab.MergeRequest) {

}

// This example demonstrates a plugin that handles HTTP requests which respond by greeting the
// world.
func main() {
	plugin.ClientMain(&Plugin{})
}
