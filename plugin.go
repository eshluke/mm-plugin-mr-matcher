package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"sync"
	"time"

	_const "com.github.eshluke.helloworld/const"
	"com.github.eshluke.helloworld/gitlab"
	"com.github.eshluke.helloworld/helper"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type configuration struct {
}

// Clone shallow copies the configuration. Your implementation may require a deep copy if
// your configuration has reference types.
func (c *configuration) Clone() *configuration {
	var clone = *c
	return &clone
}

type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

func (p *Plugin) getClient() *model.Client4 {
	client := model.NewAPIv4Client(_const.MMDOMAIN)
	client.SetToken(_const.MMTOKEN)
	return client
}

func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		// Ignore assignment if the configuration struct is empty. Go will optimize the
		// allocation for same to point at the same memory address, breaking the check
		// above.
		if reflect.ValueOf(*configuration).NumField() == 0 {
			return
		}

		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}

func (p *Plugin) OnActivate() error {

	//if err := p.OnConfigurationChange(); err != nil {
	//	return err
	//}

	return nil
}

func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	// Load the public configuration fields from the Mattermost server configuration.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return fmt.Errorf("failed to load plugin configuration: %w", err)
	}

	p.setConfiguration(configuration)

	return nil
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
	if mr.ShouldBeMatched() {
		user, drawErr := p.drawUserInChannel()
		if drawErr != nil {
			http.Error(w, fmt.Sprintf("drawUserInChannel() failed: %s", drawErr.Error()), http.StatusInternalServerError)
		}
		postErr := p.createPost(&mr, user)
		if postErr != nil {
			http.Error(w, fmt.Sprintf("createPost() failed: %s", postErr.Error()), http.StatusInternalServerError)
		}
	}
}

func (p *Plugin) drawUserInChannel() (*model.User, error) {
	//???????????? ?????? ????????? ????????????
	members, res := p.getClient().GetUsersInChannel(_const.MMCHANNELID, 0, 100, "")
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("GetUsersInChannel failed: %w", res.Error)
	}
	//????????? ???????????? ?????? ?????????
	rand.Seed(time.Now().UnixNano())
	return members[rand.Intn(len(members))], nil
}

func (p *Plugin) createPost(mr *gitlab.MergeRequest, reviewer *model.User) error {
	//????????? ????????? ???????????? ????????? ?????????
	message := fmt.Sprintf("@%s [Review Request] A new Merge Request has been created in %s by %s.",
		reviewer.Username, mr.Project.Name, mr.User.Username)

	// ??? ????????? ??????
	_, res := p.getClient().CreatePost(&model.Post{
		ChannelId: _const.MMCHANNELID,
		Message:   message,
	})
	if res.StatusCode >= 400 {
		return fmt.Errorf("create post failed: %s", res.Error.Message)
	}
	return nil
}

func main() {
	plugin.ClientMain(&Plugin{})
}
