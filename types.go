package main

import (
	"encoding/json"
	"time"
)

type config struct {
	AcmeDNSURL          string            `json:"acme_dns_url"`
	PropagationDuration duration          `json:"propagation_duration"`
	Domains             map[string]domain `json:"domains"`
}

type domain struct {
	Username            string   `json:"username"`
	Password            string   `json:"password"`
	FullDomain          string   `json:"fulldomain"`
	Subdomain           string   `json:"subdomain"`
	AllowFrom           []string `json:"allowfrom"`
	AcmeDNSURL          string   `json:"acme_dns_url"`
	PropagationDuration duration `json:"propagation_duration"`
}

type updateBody struct {
	Subdomain string `json:"subdomain"`
	Txt       string `json:"txt"`
}

type duration time.Duration

func (d duration) MarshalJSON() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

func (d *duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*d = duration(dur)

	return nil
}
