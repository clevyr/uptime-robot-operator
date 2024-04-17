package uptimerobot

import (
	"net/url"
	"strings"

	uptimerobotv1 "github.com/clevyr/uptime-robot-operator/api/v1"
)

type FindOpt func(form url.Values)

func FindBySearch(val string) FindOpt {
	return func(form url.Values) {
		form.Set("search", val)
	}
}

func FindByURL(monitor uptimerobotv1.MonitorValues) FindOpt {
	return FindBySearch(monitor.URL)
}

func FindByID(id ...string) FindOpt {
	return func(form url.Values) {
		form.Set("monitors", strings.Join(id, "-"))
	}
}
