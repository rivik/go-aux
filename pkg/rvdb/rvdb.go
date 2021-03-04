package rvdb

import (
	"fmt"
	"net/url"

	"github.com/rivik/go-aux/pkg/rvgo"
)

type DSNConfig struct {
	RawDSN       string                 `json:"dsn"`
	User         string                 `json:"user"`
	Password     string                 `json:"password"`
	PasswordFile string                 `json:"password_file"`
	RawDBParams  map[string]interface{} `json:"db_params"`

	DSNPrefix *url.URL
	DBParams  url.Values
}

func (d *DSNConfig) Prepare() error {
	u, err := url.Parse(d.RawDSN)
	if err != nil {
		return fmt.Errorf("can't parse DSN %s: %w", d.RawDSN, err)
	}

	params, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return fmt.Errorf("can't parse DSN query %s: %w", u.RawQuery, err)
	}
	u.RawQuery = ""

	for k, v := range d.RawDBParams {
		params.Set(k, fmt.Sprintf("%v", v))
	}

	d.DSNPrefix = u
	d.DBParams = params

	if d.PasswordFile != "" && d.Password == "" {
		pass, err := rvgo.TrimmedStringFromFile(d.PasswordFile)
		if err != nil {
			return fmt.Errorf("can't read passfile %s: %w", d.PasswordFile, err)
		}
		d.Password = pass
	}

	return nil
}

func (d DSNConfig) DSN() *url.URL {
	if d.DSNPrefix == nil || d.DBParams == nil {
		return nil
	}

	u := d.DSNPrefix
	u.RawQuery = d.DBParams.Encode()
	return u
}
