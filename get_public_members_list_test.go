package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicholasf/fakepoint"

	"appengine/aetest"
)

func TestCallPublicMembersListEndpoint(t *testing.T) {
	maker := fakepoint.NewFakepointMaker()
	maker.NewGet("https://api.github.com/orgs/crowdint/public_members", 200).
		SetResponse(`
[
  {
    "login": "octocat",
    "id": 1,
    "avatar_url": "https://github.com/images/error/octocat_happy.gif",
    "gravatar_id": "",
    "url": "https://api.github.com/users/octocat",
    "html_url": "https://github.com/octocat",
    "followers_url": "https://api.github.com/users/octocat/followers",
    "following_url": "https://api.github.com/users/octocat/following{/other_user}",
    "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
    "organizations_url": "https://api.github.com/users/octocat/orgs",
    "repos_url": "https://api.github.com/users/octocat/repos",
    "events_url": "https://api.github.com/users/octocat/events{/privacy}",
    "received_events_url": "https://api.github.com/users/octocat/received_events",
    "type": "User",
    "site_admin": false
  }
]`).
		SetHeader("Content-Type", "application/json")

	c, err := aetest.NewContext(&aetest.Options{"testapp", true})

	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer c.Close()

	mc := &MyContext{Env: "test", Context: c, Client: maker.Client()}
	wrapee := Wrap(GetPublicMembersList, mc)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/pull-people", nil)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	wrapee(w, req)

	if w.Code != 200 {
		t.Log(w.Code)
		t.Log(w.Body.String())
		t.Fail()
	}

}

func TestCallPublicMembersListEndpointWithMultiplePages(t *testing.T) {
	maker := fakepoint.NewFakepointMaker()
	maker.NewGet("https://api.github.com/orgs/crowdint/public_members", 200).
		SetResponse(`
[
  {
    "login": "octocat",
    "id": 1,
    "avatar_url": "https://github.com/images/error/octocat_happy.gif",
    "gravatar_id": "",
    "url": "https://api.github.com/users/octocat",
    "html_url": "https://github.com/octocat",
    "followers_url": "https://api.github.com/users/octocat/followers",
    "following_url": "https://api.github.com/users/octocat/following{/other_user}",
    "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
    "organizations_url": "https://api.github.com/users/octocat/orgs",
    "repos_url": "https://api.github.com/users/octocat/repos",
    "events_url": "https://api.github.com/users/octocat/events{/privacy}",
    "received_events_url": "https://api.github.com/users/octocat/received_events",
    "type": "User",
    "site_admin": false
  }
]`).
		SetHeader("Content-Type", "application/json").
		SetHeader("Link", "<https://api.github.com/orgs/crowdint/public_members?page=2>; rel=\"next\"")

	maker.NewGet("https://api.github.com/orgs/crowdint/public_members?page=2", 200).
		SetResponse(`
[
  {
    "login": "octocat 2",
    "id": 2,
    "avatar_url": "https://github.com/images/error/octocat_happy.gif",
    "gravatar_id": "",
    "url": "https://api.github.com/users/octocat2",
    "html_url": "https://github.com/octocat2",
    "followers_url": "https://api.github.com/users/octocat2/followers",
    "following_url": "https://api.github.com/users/octocat2/following{/other_user}",
    "gists_url": "https://api.github.com/users/octocat2/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/octocat2/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/octocat2/subscriptions",
    "organizations_url": "https://api.github.com/users/octocat2/orgs",
    "repos_url": "https://api.github.com/users/octocat2/repos",
    "events_url": "https://api.github.com/users/octocat2/events{/privacy}",
    "received_events_url": "https://api.github.com/users/octocat2/received_events",
    "type": "User",
    "site_admin": false
  }
]`).
		SetHeader("Content-Type", "application/json")

	c, err := aetest.NewContext(&aetest.Options{"testapp", true})

	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer c.Close()

	mc := &MyContext{Env: "test", Context: c, Client: maker.Client()}
	wrapee := Wrap(GetPublicMembersList, mc)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/pull-people", nil)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	wrapee(w, req)

	if w.Code != 200 {
		t.Log(w.Code)
		t.Log(w.Body.String())
		t.Fail()
	}

	var members []Member
	json.Unmarshal(w.Body.Bytes(), &members)

	if len(members) != 2 {
		t.Fail()
	}
}
