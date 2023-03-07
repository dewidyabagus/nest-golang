package main

import (
	"fmt"
	"net/http"

	resty "github.com/go-resty/resty/v2"
)

type notify struct {
	rest *resty.Client
}

func NewNotify(cfg NotifyConfig) *notify {
	r := resty.New()
	r.SetBaseURL(cfg.Host)
	r.SetHeader("X-API-KEY", cfg.ApiKey)

	return &notify{r}
}

func (n *notify) EmailNotify(notif string, info PayloadNotify) error {
	fullPath := fmt.Sprintf("%s/email/%s", n.rest.BaseURL, notif)

	res, err := n.rest.R().SetBody(info).Post(fullPath)
	if err != nil {
		return err

	} else if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("code=%d message=%s", res.StatusCode(), res.String())
	}

	return nil
}
