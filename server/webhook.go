package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/bobheadxi/res"
	"github.com/bobheadxi/timelines/db"
	"github.com/bobheadxi/timelines/host"
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
	case *github.InstallationRepositoriesEvent:
		// https://developer.github.com/v3/activity/events/types/#installationevent
		install := event.GetInstallation()
		h.l.Infof("received installation %#v", install)
		if err := h.handleInstall(r.Context(),
			host.InstallationFromGitHub(install),
			host.ReposFromGitHub(event.RepositoriesAdded),
			host.ReposFromGitHub(event.RepositoriesRemoved),
		); err != nil {
			res.R(w, r, res.ErrBadRequest("unexpected error encountered",
				"error", err))
			return
		}
		res.R(w, r, res.MsgOK("event acknowledged but not processed - not implemented",
			"type", t))

	case *github.InstallationEvent, *github.MarketplacePurchaseEvent:
		// Installation updates
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

func (h *webhookHandler) handleInstall(
	ctx context.Context,
	install host.Installation,
	added, removed []host.Repo,
) error {
	errSet := map[string]error{}
	for _, repo := range added {
		if err := h.db.Repos().NewRepository(
			ctx,
			install.GetID(),
			repo,
		); err != nil {
			if !db.IsNotFound(err) {
				errSet[repo.GetOwner()+"/"+repo.GetName()] = err
			} else {
				h.l.Errorw("unexpected not-found error",
					"repo", repo,
					"install", install,
					"error", err)
			}
		}
		h.store.RepoJobs().Queue(&store.RepoJob{
			ID:             uuid.New(),
			Owner:          repo.GetOwner(),
			Repo:           repo.GetName(),
			InstallationID: install.GetID(), // TODO: not needed?
		})
	}
	for _, repo := range removed {
		dbr, err := h.db.Repos().GetRepository(
			ctx,
			repo.GetHost(),
			repo.GetOwner(),
			repo.GetName())
		if err != nil {
			if !db.IsNotFound(err) {
				errSet[repo.GetOwner()+"/"+repo.GetName()] = err
			} else {
				h.l.Errorw("unexpected not-found error",
					"repo", repo,
					"install", install,
					"error", err)
			}
			continue
		}

		if err := h.db.Repos().DeleteRepository(ctx, dbr.ID); err != nil {
			errSet[repo.GetOwner()+"/"+repo.GetName()] = err
		}
	}
	if len(errSet) > 0 {
		h.l.Errorw("error(s) occured on installation handling",
			"errors", errSet)
		return fmt.Errorf("errors: %+v", errSet)
	}
	return nil
}
