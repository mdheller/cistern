package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nbedos/citop/cache"
	"github.com/nbedos/citop/utils"
)

type AppVeyorClient struct {
	client      *http.Client
	rateLimiter <-chan time.Time
	token       string
	accountID   string
}

var AppVeyorURL = url.URL{
	Scheme:  "https",
	Host:    "ci.appveyor.com",
	Path:    "/api",
	RawPath: "/api",
}

func NewAppVeyorClient(accountID string, token string, rateLimit time.Duration) AppVeyorClient {
	return AppVeyorClient{
		client:      &http.Client{Timeout: 10 * time.Second},
		rateLimiter: time.Tick(rateLimit),
		token:       token,
		accountID:   accountID,
	}
}

func (c AppVeyorClient) AccountID() string {
	return c.accountID
}

func (c AppVeyorClient) Log(ctx context.Context, repository cache.Repository, jobID int) (string, bool, error) {
	return "", true, nil
}

func (c AppVeyorClient) BuildFromURL(ctx context.Context, u string) (cache.Build, error) {
	owner, repo, id, err := parseAppVeyorURL(u)
	if err != nil {
		return cache.Build{}, err
	}

	return c.fetchBuild(ctx, owner, repo, id)
}

func (c AppVeyorClient) getJSON(ctx context.Context, u url.URL, i interface{}) error {
	body, err := c.get(ctx, u)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body.Bytes(), i); err != nil {
		return err
	}

	return nil
}

func (c AppVeyorClient) get(ctx context.Context, u url.URL) (*bytes.Buffer, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req = req.WithContext(ctx)

	select {
	case <-c.rateLimiter:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // FIXME return error

	body := new(bytes.Buffer)
	if _, err := body.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = HTTPError{
			Method:  req.Method,
			URL:     req.URL.String(),
			Status:  resp.StatusCode,
			Message: body.String(),
		}
	}

	return body, err
}

func (c AppVeyorClient) fetchBuild(ctx context.Context, owner string, repoName string, id int) (cache.Build, error) {
	// We only have the build ID and need a build object. We have to query two endpoints:
	// 		1. /projects/owner/repoName/history with startBuildId = id gives us a build object with
	//      an empty job list but with a version number
	//      2. /projects/owner/repoName/build/<version> using the version number from the last call
	//      gives us a build with a complete job list
	history := AppVeyorURL
	historyFormat := "/projects/%s/%s/history"
	history.Path += fmt.Sprintf(historyFormat, owner, repoName)
	history.RawPath += fmt.Sprintf(historyFormat, url.PathEscape(owner), url.PathEscape(repoName))
	params := history.Query()
	params.Add("recordsNumber", "1")
	params.Add("startBuildId", strconv.Itoa(id+1))
	history.RawQuery = params.Encode()

	var b struct {
		Project struct {
			ID    int    `json:"projectId"`
			Owner string `json:"accountName"`
			Name  string `json:"name"`
		}
		Builds []appVeyorBuild `json:"builds"`
	}
	if err := c.getJSON(ctx, history, &b); err != nil {
		return cache.Build{}, err
	}

	if len(b.Builds) != 1 {
		return cache.Build{}, fmt.Errorf("found no build with id %d", id)
	}
	if b.Builds[0].ID != id {
		return cache.Build{}, fmt.Errorf("expected build #%d but got %d", id, b.Builds[0].ID)
	}
	version := b.Builds[0].Version

	repository := cache.Repository{
		AccountID: c.accountID,
		ID:        b.Project.ID,
		URL:       "",
		Owner:     b.Project.Owner,
		Name:      b.Project.Name,
	}

	endpoint := AppVeyorURL
	pathFormat := "/projects/%s/%s/build/%s"
	endpoint.Path += fmt.Sprintf(pathFormat, owner, repoName, version)
	endpoint.RawPath += fmt.Sprintf(pathFormat, url.PathEscape(owner),
		url.PathEscape(repoName), url.PathEscape(version))
	var bVersion struct {
		Build appVeyorBuild `json:"build"`
	}
	if err := c.getJSON(ctx, endpoint, &bVersion); err != nil {
		return cache.Build{}, err
	}

	return bVersion.Build.toCacheBuild(c.accountID, &repository)
}

// Extract owner, repository and build ID from web URL of build
func parseAppVeyorURL(u string) (string, string, int, error) {
	v, err := url.Parse(u)
	if err != nil {
		return "", "", 0, err
	}

	if !strings.HasSuffix(v.Hostname(), "appveyor.com") {
		return "", "", 0, cache.ErrUnknownURL
	}

	// URL format: https://ci.appveyor.com/project/nbedos/citop/builds/29070120
	cs := strings.Split(v.EscapedPath(), "/")
	if len(cs) < 6 || cs[1] != "project" || cs[4] != "builds" {
		return "", "", 0, cache.ErrUnknownURL
	}

	owner, repo := cs[2], cs[3]
	id, err := strconv.Atoi(cs[5])
	if err != nil {
		return "", "", 0, err
	}
	return owner, repo, id, nil
}

func fromAppVeyorState(s string) cache.State {
	switch strings.ToLower(s) {
	case "queued", "received":
		return cache.Pending
	case "success":
		return cache.Passed
	case "failed":
		return cache.Failed
	default:
		panic(fmt.Errorf("unknown state: %q", s))
		return cache.Unknown
	}
}

type appVeyorBuild struct {
	ID          int           `json:"buildId"`
	Jobs        []appVeyorJob `json:"jobs"`
	Number      int           `json:"buildNumber"`
	Version     string        `json:"version"`
	Message     string        `json:"message"`
	Branch      string        `json:"branch"`
	IsTag       bool          `json:"isTag"`
	Sha         string        `json:"commitId"`
	Author      string        `json:"authorUsername"`
	CommittedAt string        `json:"committed"`
	Status      string        `json:"status"`
	CreatedAt   string        `json:"created"`
	StartedAt   string        `json:"started"`
	FinishedAt  string        `json:"finished"`
	UpdatedAt   string        `json:"updated"`
}

func (b appVeyorBuild) toCacheBuild(accountID string, repo *cache.Repository) (cache.Build, error) {
	build := cache.Build{
		Repository: repo,
		ID:         strconv.Itoa(b.ID),
		Commit: cache.Commit{
			Sha:     b.Sha,
			Message: b.Message,
		},
		Ref:             b.Branch,
		IsTag:           b.IsTag,
		RepoBuildNumber: strconv.Itoa(b.Number),
		State:           fromAppVeyorState(b.Status),
		Stages:          make(map[int]*cache.Stage),
		Jobs:            make(map[int]*cache.Job),
	}
	var err error
	build.Commit.Date, err = utils.NullTimeFromString(b.CommittedAt)
	if err != nil {
		return build, err
	}
	build.CreatedAt, err = utils.NullTimeFromString(b.CreatedAt)
	if err != nil {
		return build, err
	}
	build.StartedAt, err = utils.NullTimeFromString(b.StartedAt)
	if err != nil {
		return build, err
	}
	build.FinishedAt, err = utils.NullTimeFromString(b.FinishedAt)
	if err != nil {
		return build, err
	}
	if build.UpdatedAt, err = time.Parse(time.RFC3339, b.UpdatedAt); err != nil {
		return build, err
	}

	if build.FinishedAt.Valid && build.StartedAt.Valid {
		build.Duration = utils.NullDuration{
			Valid:    true,
			Duration: build.FinishedAt.Time.Sub(build.StartedAt.Time),
		}
	}

	build.WebURL = fmt.Sprintf("https://ci.appveyor.com/project/%s/%s/builds/%d",
		url.PathEscape(repo.Owner), url.PathEscape(repo.Name), b.ID)

	for i, job := range b.Jobs {
		j, err := job.toCacheJob(i + 1)
		if err != nil {
			return build, err
		}
		build.Jobs[j.ID] = &j
	}

	return build, nil
}

type appVeyorJob struct {
	Name         string `json:"name"`
	AllowFailure bool   `json:"allowFailure"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created"`
	StartedAt    string `json:"started"`
	FinishedAt   string `json:"finished"`
}

func (j appVeyorJob) toCacheJob(id int) (cache.Job, error) {
	if id <= 0 {
		return cache.Job{}, errors.New("job id must be > 0")
	}
	job := cache.Job{
		ID:           id,
		State:        fromAppVeyorState(j.Status),
		Name:         j.Name,
		WebURL:       "",
		AllowFailure: j.AllowFailure,
	}

	var err error
	job.CreatedAt, err = utils.NullTimeFromString(j.CreatedAt)
	if err != nil {
		return job, err
	}
	job.StartedAt, err = utils.NullTimeFromString(j.StartedAt)
	if err != nil {
		return job, err
	}
	job.FinishedAt, err = utils.NullTimeFromString(j.FinishedAt)
	if err != nil {
		return job, err
	}
	if job.FinishedAt.Valid && job.StartedAt.Valid {
		job.Duration = utils.NullDuration{
			Valid:    true,
			Duration: job.FinishedAt.Time.Sub(job.StartedAt.Time),
		}
	}

	return job, nil
}
