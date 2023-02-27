package main

import (
	"errors"
	"fmt"
	"net/http"

	"com.github.eshluke.helloworld/const"
	"com.github.eshluke.helloworld/gitlab"
	"com.github.eshluke.helloworld/helper"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin
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
		p.createPost(&mr)
	}
}

func (p *Plugin) createPost(mr *gitlab.MergeRequest) {
	//채널에서 맴버 리스트 불러오기
	members, resp := p.API.GetUsersInChannel(MMCHANNELID, "", 0, 100)
	if resp.Error != nil {
		http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
		return
	}
	//그중에 랜덤으로 한명 고르기
	rand.Seed(time.Now().UnixNano())
	randomMember := members[rand.Intn(len(members))]

	//골라진 사람을 맨션하는 포스팅 올리기
	message := fmt.Sprintf("@%s [Review Request] A new Merge Request has been created in %s by %s.",
		randomMember.Username, mrEvent.Project.Name, mr.User.Name)

	// 새 게시물 작성
	post := &model.Post{
		ChannelId: channelID,
		Message:   message,
	}
	_, resp := p.API.CreatePost(post)
	if resp.Error != nil {
		http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
		return
	}
}

// This example demonstrates a plugin that handles HTTP requests which respond by greeting the
// world.
func main() {
	plugin.ClientMain(&Plugin{})
}
