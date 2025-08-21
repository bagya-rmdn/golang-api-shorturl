package utils

import (
"net/url"
"strings"
)

func NormalizeURL(raw string) (string, error) {
s := strings.TrimSpace(raw)
if s == "" {
return "", nil
}
u, err := url.Parse(s)
if err != nil {
return "", err
}
if u.Scheme == "" {
// default to https
u.Scheme = "https"
}
// Lowercase host for idempotency
u.Host = strings.ToLower(u.Host)
return u.String(), nil
}