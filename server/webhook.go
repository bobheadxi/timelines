package server

import (
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"github.com/bobheadxi/res"
	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/store"
)

type webhookHandler struct {
	db    *db.Database
	store *store.Client

	l *zap.SugaredLogger
}

func newWebhookHandler(
	l *zap.SugaredLogger,
	database *db.Database,
	store *store.Client,
) *webhookHandler {
	return &webhookHandler{database, store, l}
}

func (h *webhookHandler) handleGitHub(w http.ResponseWriter, r *http.Request) {
	t := github.WebHookType(r)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res.R(w, r, res.ErrBadRequest("unable to read request", "error", err))
		return
	}
	payload, err := github.ParseWebHook(t, b)
	if err != nil {
		res.R(w, r, res.ErrBadRequest("unable to parse payload", "error", err))
		return
	}

	switch event := payload.(type) {
	case *github.InstallationEvent, *github.InstallationRepositoriesEvent, *github.MarketplacePurchaseEvent:
		// https://developer.github.com/v3/activity/events/types/#installationevent
		h.l.Infof("received %#v", event)
		// TODO handle all installation-related events together
		res.R(w, r, res.MsgOK("event acknowledged but not processed - not implemented",
			"type", t))

	case *github.CreateEvent, *github.DeleteEvent, *github.MilestoneEvent, *github.ReleaseEvent:
		// https://developer.github.com/v3/activity/events/types/#createevent
		// https://developer.github.com/v3/activity/events/types/#milestoneevent
		// https://developer.github.com/v3/activity/events/types/#releaseevent
		h.l.Infof("received %#v", event)
		// TODO handle tag, milestone sync here - call these "milestones" or something
		res.R(w, r, res.MsgOK("event acknowledged but not processed - not implemented",
			"type", t))

	case *github.IssuesEvent, *github.PullRequest:
		h.l.Infof("received %#v", event)
		// TODO manage issue, pull request updates - aka "items"
		res.R(w, r, res.MsgOK("event acknowledged but not processed - not implemented",
			"type", t))

	case *github.PushEvent:
		h.l.Infof("received %#v", event)
		// TODO manage job queues here for repository updates
		res.R(w, r, res.MsgOK("event acknowledged but not processed - not implemented",
			"type", t))

	default:
		h.l.Infof("unknown type %#v", event)
		res.R(w, r, res.MsgOK("event acknowledged but not processed",
			"type", t))
	}
}
