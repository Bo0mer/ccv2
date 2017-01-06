package ccv2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

//go:generate counterfeiter . Doer

// Filter specifies the target of a query.
type Filter string

const (
	// FilterName specifies that the query should filter on the name field.
	FilterName Filter = "name"
	// FilterOrganizationGUID specifies that the query should filter on the
	// organization_guid field.
	FilterOrganizationGUID = "organization_guid"
	// FilterSpaceGUID specifies that the query should filter on the space_guid
	// field.
	FilterSpaceGUID = "space_guid"
	// FilterTimestamp specifies that the query should filter on the timestamp
	// field.
	FilterTimestamp = "timestamp"
	// FilterActee specifies that the query should filter on the actee field.
	FilterActee = "actee"
	// FilterActor specifies that the query should filter on the actor field.
	FilterActor = "actor"
	// FilterType specifies that the query should filter on the type field.
	FilterType = "type"
)

// Operator specifies an operator for a query.
type Operator string

const (
	// OperatorEqual specifies that the result should match the value.
	OperatorEqual Operator = ":"
	// OperatorGreater specifies that the result should be greater than the value.
	OperatorGreater = ">"
	// OperatorLess specifies that the result should be less than the value.
	OperatorLess = "<"
)

// Query gives means to filter list of resources.
type Query struct {
	// Filter is the field on which the query will act.
	Filter Filter
	// Op is the operator that will be applied during filtering.
	Op Operator
	// Value is the value that will be used when applying the Op.
	Value string
}

// String returns the string representation of the query.
// It is of the form <Filter><Op><Value>.
// It must be query encoded if to be used as a query param.
func (q Query) String() string {
	return string(q.Filter) + string(q.Op) + q.Value
}

// Doer does HTTP requests and returns the corresponding responses.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// UnexpectedResponseError wraps response that indicates error.
type UnexpectedResponseError struct {
	StatusCode  int    `json:"status_code"`
	Description string `json:"description"`
	ErrorCode   string `json:"error_code"`
}

// Error returns a description of the error.
func (e *UnexpectedResponseError) Error() string {
	return e.Description
}

type paginatedResource struct {
	NextURL   string          `json:"next_url"`
	Resources json.RawMessage `json:"resources"`
}

// Client implements a read-only Cloud Controller client.
type Client struct {
	API        *url.URL
	HTTPClient Doer
}

type requestOpts struct {
	Context context.Context
	Method  string
	Path    string
	Queries []Query
	Body    io.Reader
}

// Info returns the info returned by the Cloud Controller info endpoint.
func (c *Client) Info(ctx context.Context) (i Info, err error) {
	req, err := c.newRequest(requestOpts{
		Context: ctx,
		Method:  http.MethodGet,
		Path:    "/v2/info",
	})
	if err != nil {
		return Info{}, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Info{}, errors.Wrap(err, "doing info request failed")
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != 200 {
		return Info{}, errFromResponse(resp)
	}

	var info Info
	err = json.NewDecoder(resp.Body).Decode(&info)
	return info, err
}

// Organizations list all organizations that conform to the provided queries.
func (c *Client) Organizations(ctx context.Context, queries ...Query) ([]Organization, error) {
	var orgs []Organization
	orgCb := func(resources json.RawMessage) error {
		var res []Organization
		err := json.Unmarshal(resources, &res)
		orgs = append(orgs, res...)
		return err
	}

	opts := requestOpts{
		Context: ctx,
		Method:  http.MethodGet,
		Path:    "/v2/organizations",
		Queries: queries,
	}
	err := c.paginate(opts, orgCb)
	return orgs, err
}

// Spaces list all spaces that conform to the provided queries.
func (c *Client) Spaces(ctx context.Context, queries ...Query) ([]Space, error) {
	var spaces []Space
	spaceCb := func(resources json.RawMessage) error {
		var res []Space
		err := json.Unmarshal(resources, &res)
		spaces = append(spaces, res...)
		return err
	}

	opts := requestOpts{
		Context: ctx,
		Method:  http.MethodGet,
		Path:    "/v2/spaces",
		Queries: queries,
	}
	err := c.paginate(opts, spaceCb)
	return spaces, err
}

// Applications list all applications that conform to the provided queries.
func (c *Client) Applications(ctx context.Context, queries ...Query) ([]Application, error) {
	var apps []Application
	appsCb := func(resources json.RawMessage) error {
		var res []Application
		err := json.Unmarshal(resources, &res)
		apps = append(apps, res...)
		return err
	}

	opts := requestOpts{
		Context: ctx,
		Method:  http.MethodGet,
		Path:    "/v2/apps",
		Queries: queries,
	}
	err := c.paginate(opts, appsCb)
	return apps, err
}

// ApplicationSummary returns summary for a given application.
func (c *Client) ApplicationSummary(ctx context.Context, app Application) (a ApplicationSummary, err error) {
	opts := requestOpts{
		Context: ctx,
		Method:  http.MethodGet,
		Path:    fmt.Sprintf("/v2/apps/%s/summary", app.GUID),
	}
	req, err := c.newRequest(opts)
	if err != nil {
		return ApplicationSummary{}, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return ApplicationSummary{}, errors.Wrap(err, "doing summary request failed")
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return ApplicationSummary{}, errFromResponse(resp)
	}

	var summary ApplicationSummary
	err = json.NewDecoder(resp.Body).Decode(&summary)
	return summary, errors.Wrap(err, "decoding response failed")
}

// Events list all events that conform to the provided queries.
func (c *Client) Events(ctx context.Context, queries ...Query) ([]Event, error) {
	var events []Event
	eventCb := func(resources json.RawMessage) error {
		var res []Event
		err := json.Unmarshal(resources, &res)
		events = append(events, res...)
		return err
	}

	opts := requestOpts{
		Context: ctx,
		Method:  http.MethodGet,
		Path:    "/v2/events",
		Queries: queries,
	}
	err := c.paginate(opts, eventCb)
	return events, err
}

func (c *Client) newRequest(opts requestOpts) (*http.Request, error) {
	url, err := c.API.Parse(opts.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing url for path %q failed", opts.Path)
	}
	q := url.Query()
	for _, query := range opts.Queries {
		q.Add("q", query.String())
	}
	url.RawQuery = q.Encode()
	req, err := http.NewRequest(opts.Method, url.String(), opts.Body)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest failed")
	}
	return req.WithContext(opts.Context), nil
}

func (c *Client) paginate(opts requestOpts, pageCb func(json.RawMessage) error) (err error) {
	for {
		req, err := c.newRequest(opts)
		if err != nil {
			return errors.Wrap(err, "creating page request failed")
		}
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return errors.Wrap(err, "doing page request failed")
		}
		defer func(resp *http.Response) {
			if cerr := resp.Body.Close(); cerr != nil && err == nil {
				err = cerr
			}
		}(resp)

		if resp.StatusCode != http.StatusOK {
			return errFromResponse(resp)
		}
		var page paginatedResource
		if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
			return errors.Wrap(err, "page response decoding failed")
		}
		if err := pageCb(page.Resources); err != nil {
			return errors.Wrap(err, "page content processor failed")
		}
		if page.NextURL == "" {
			break
		}
		opts.Path = page.NextURL
	}
	return nil
}

func errFromResponse(resp *http.Response) error {
	e := &UnexpectedResponseError{
		StatusCode: resp.StatusCode,
	}
	if err := json.NewDecoder(resp.Body).Decode(e); err != nil {
		e.Description = err.Error()
	}
	return e
}
