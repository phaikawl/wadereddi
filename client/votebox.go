package client

import (
	"github.com/phaikawl/wade/rt"
	"github.com/phaikawl/wade/utils"
	c "github.com/phaikawl/wadereddi/common"
)

// VoteBoxModel is the prototype for the "votebox" custom element
type VoteBoxModel struct {
	rt.BaseProto
	Vote      *c.Score
	VoteUrl   string
	AfterVote func() // function to be called after a vote is done
}

// DoVote performs the vote with the given point changed (-1 = vote down, 1 = vote up)
func (m *VoteBoxModel) DoVote(vote int) {
	if m.VoteUrl == "" {
		panic("VoteUrl has not been set.")
	}

	if vote < -1 || vote > 1 {
		panic("Illegal vote value.")
	}

	lastVote := m.Vote.Voted
	m.Vote.Voted = vote
	if vote == lastVote {
		vote, lastVote = -lastVote, 0
		m.Vote.Voted = 0
	}

	url := utils.UrlQuery(m.VoteUrl, utils.Map{
		"vote":     vote,
		"lastvote": lastVote,
	})

	// performs an http request to the server to vote, and assign the updated score
	// to m.Vote.Score after that
	go func() {
		r, _ := app().Http.GET(url)
		r.ParseJSON(&m.Vote.Score)

		if m.AfterVote != nil {
			m.AfterVote()
		}

		app().Render()
	}()
}
