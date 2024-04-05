package uptimerobot

import (
	"net/url"
	"strings"
)

type FindOpt func(form url.Values)

func FindBySearch(val string) FindOpt {
	return func(form url.Values) {
		form.Set("search", val)
	}
}

func FindByURL(monitor Monitor) FindOpt {
	return FindBySearch(monitor.URL)
}

func FindByID(id ...string) FindOpt {
	return func(form url.Values) {
		form.Set("monitors", strings.Join(id, "-"))
	}
}
