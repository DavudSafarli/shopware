package internal

import (
	"fmt"
	"net/url"
	"sort"
)

type FullMatchRule struct {
	FromRaw string
	// FromCanonical(maybe rename to FromNormalized) is full path but the query string variables are sorted alphabetically
	FromCanonical string
	Target        string
}

func NewFullMatchRule(fromRaw string, target string) (*FullMatchRule, error) {
	urlFrom, err := url.Parse(fromRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse From url: %w", err)
	}

	_, err = url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Target url: %w", err)
	}

	r := &FullMatchRule{
		FromRaw: fromRaw,
		Target:  target,
	}

	r.FromCanonical = canonicalPathQuery(urlFrom)
	return r, nil
}

func canonicalPathQuery(u *url.URL) string {
	q := u.Query()
	for k, v := range q {
		if len(v) > 1 {
			sort.Strings(v)
			q[k] = v
		}
	}

	// keys are sorted alphabetically by #Encode
	encoded := q.Encode()
	path := u.EscapedPath()
	if path == "" {
		path = "/"
	}
	if encoded == "" {
		return path
	}
	return path + "?" + encoded
}
