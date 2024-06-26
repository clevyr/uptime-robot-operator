package uptimerobot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	uptimerobotv1 "github.com/clevyr/uptime-robot-operator/api/v1"
	"github.com/clevyr/uptime-robot-operator/internal/uptimerobot/urtypes"
)

func NewClient(apiKey string) Client {
	api := "https://api.uptimerobot.com/v2"
	if env := os.Getenv("UPTIME_ROBOT_API"); env != "" {
		api = strings.TrimSuffix(env, "/")
	}

	return Client{url: api, apiKey: apiKey}
}

type Client struct {
	url    string
	apiKey string
}

func (c Client) NewValues() url.Values {
	v := make(url.Values)
	v.Set("api_key", c.apiKey)
	v.Set("format", "json")
	return v
}

func (c Client) NewRequest(ctx context.Context, endpoint string, form url.Values) (*http.Request, error) {
	u := c.url + "/" + endpoint

	req, err := http.NewRequestWithContext(ctx, "POST", u, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (c Client) Do(ctx context.Context, endpoint string, form url.Values) (*http.Response, error) {
	req, err := c.NewRequest(ctx, endpoint, form)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 400 {
		return nil, fmt.Errorf("%w: %s", ErrStatus, res.Status)
	}

	return res, nil
}

func (c Client) MonitorValues(monitor uptimerobotv1.MonitorValues, form url.Values, contacts uptimerobotv1.MonitorContacts) url.Values {
	form.Set("friendly_name", monitor.Name)
	form.Set("url", monitor.URL)
	form.Set("type", strconv.Itoa(int(monitor.Type)))
	form.Set("interval", strconv.Itoa(int(monitor.Interval.Seconds())))
	form.Set("timeout", strconv.Itoa(int(monitor.Timeout.Seconds())))
	form.Set("alert_contacts", contacts.String())
	form.Set("http_method", strconv.Itoa(int(monitor.Method)))
	switch monitor.Method {
	case urtypes.MethodHEAD, urtypes.MethodGET:
	default:
		if monitor.POST != nil {
			form.Set("post_type", strconv.Itoa(int(monitor.POST.Type)))
			form.Set("post_content_type", strconv.Itoa(int(monitor.POST.ContentType)))
			form.Set("post_value", monitor.POST.Value)
		}
	}
	var username, password string
	if monitor.Auth != nil {
		form.Set("http_auth_type", strconv.Itoa(int(monitor.Auth.Type)))
		username = monitor.Auth.Username
		password = monitor.Auth.Password
	}
	form.Set("http_username", username)
	form.Set("http_password", password)
	switch monitor.Type {
	case urtypes.TypeKeyword:
		if monitor.Keyword != nil {
			form.Set("keyword_type", strconv.Itoa(int(monitor.Keyword.Type)))
			caseType := "1"
			if *monitor.Keyword.CaseSensitive {
				caseType = "0"
			}
			form.Set("keyword_case_type", caseType)
			form.Set("keyword_value", monitor.Keyword.Value)
		}
	case urtypes.TypePort:
		if monitor.Port != nil {
			form.Set("sub_type", strconv.Itoa(int(monitor.Port.Type)))
			if monitor.Port.Type == urtypes.PortCustom {
				form.Set("port", strconv.Itoa(int(monitor.Port.Type)))
			}
		}
	}
	return form
}

type Response struct {
	Status  urtypes.Status  `json:"stat"`
	Monitor ResponseMonitor `json:"monitor"`
}

type ResponseMonitor struct {
	ID json.Number `json:"id"`
}

var (
	ErrStatus   = errors.New("error code from Uptime Robot API")
	ErrResponse = errors.New("received fail from Uptime Robot API")
)

func (c Client) CreateMonitor(ctx context.Context, monitor uptimerobotv1.MonitorValues, contacts uptimerobotv1.MonitorContacts) (string, error) {
	form := c.MonitorValues(monitor, c.NewValues(), contacts)

	res, err := c.Do(ctx, "newMonitor", form)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	parsed := &Response{}
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return "", err
	}

	if parsed.Status != urtypes.StatusOK {
		if id, err := c.FindMonitorID(ctx, FindByURL(monitor)); err == nil {
			// Monitor already exists
			return id, nil
		}
		return "", ErrResponse
	}
	return parsed.Monitor.ID.String(), nil
}

func (c Client) DeleteMonitor(ctx context.Context, id string) error {
	form := c.NewValues()
	form.Set("id", id)

	res, err := c.Do(ctx, "deleteMonitor", form)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	parsed := &Response{}
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return err
	}

	if parsed.Status != urtypes.StatusOK {
		if _, err := c.FindMonitorID(ctx, FindByID(id)); err != nil && errors.Is(err, ErrMonitorNotFound) {
			// Monitor already deleted
			return nil
		}
		return ErrResponse
	}
	return nil
}

func (c Client) EditMonitor(ctx context.Context, id string, monitor uptimerobotv1.MonitorValues, contacts uptimerobotv1.MonitorContacts) (string, error) {
	form := c.MonitorValues(monitor, c.NewValues(), contacts)
	form.Set("id", id)
	form.Set("status", strconv.Itoa(int(monitor.Status)))

	res, err := c.Do(ctx, "editMonitor", form)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	parsed := &Response{}
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return "", err
	}

	if parsed.Status != urtypes.StatusOK {
		if _, err := c.FindMonitorID(ctx, FindByID(id)); err != nil && errors.Is(err, ErrMonitorNotFound) {
			// Recreate deleted monitor
			return c.CreateMonitor(ctx, monitor, contacts)
		}
		return parsed.Monitor.ID.String(), ErrResponse
	}
	return parsed.Monitor.ID.String(), nil
}

type FindMonitorResponse struct {
	Status   urtypes.Status    `json:"stat"`
	Monitors []ResponseMonitor `json:"monitors"`
}

var ErrMonitorNotFound = errors.New("monitor not found")

func (c Client) FindMonitorID(ctx context.Context, opts ...FindOpt) (string, error) {
	form := c.NewValues()
	for _, opt := range opts {
		opt(form)
	}

	res, err := c.Do(ctx, "getMonitors", form)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	parsed := &FindMonitorResponse{}
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return "", err
	}

	if parsed.Status != urtypes.StatusOK {
		return "", ErrResponse
	}

	for _, monitor := range parsed.Monitors {
		return monitor.ID.String(), nil
	}
	return "", ErrMonitorNotFound
}

type FindContactResponse struct {
	Status   urtypes.Status    `json:"stat"`
	Contacts []ResponseContact `json:"alert_contacts"`
}

type ResponseContact struct {
	ID           string `json:"id"`
	FriendlyName string `json:"friendly_name"`
}

var ErrContactNotFound = errors.New("contact not found")

func (c Client) FindContactID(ctx context.Context, friendlyName string) (string, error) {
	form := c.NewValues()
	res, err := c.Do(ctx, "getAlertContacts", form)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	parsed := &FindContactResponse{}
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return "", err
	}

	if parsed.Status != urtypes.StatusOK {
		return "", ErrResponse
	}

	for _, contact := range parsed.Contacts {
		if friendlyName == contact.FriendlyName {
			return contact.ID, nil
		}
	}
	return "", ErrContactNotFound
}

type AccountDetailsResponse struct {
	Status  urtypes.Status `json:"stat"`
	Account Account        `json:"account"`
}

type Account struct {
	Email string `json:"email"`
}

func (c Client) GetAccountDetails(ctx context.Context) (string, error) {
	form := c.NewValues()
	res, err := c.Do(ctx, "getAccountDetails", form)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	parsed := &AccountDetailsResponse{}
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return "", err
	}

	if parsed.Status != urtypes.StatusOK {
		return "", ErrResponse
	}
	return parsed.Account.Email, err
}
