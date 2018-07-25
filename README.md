# Okta

Okta client written in Go

[![GoDoc](https://godoc.org/github.com/arkan/okta?status.svg)](https://godoc.org/github.com/arkan/okta)
[![Go Report Card](https://goreportcard.com/badge/github.com/arkan/okta)](https://goreportcard.com/report/github.com/arkan/okta)


## Getting Started
```
go get github.com/arkan/okta
```

## Create a new API Key on Okta

First you need [to create a new api key on Okta](https://heetch-admin.okta.com/admin/access/api/tokens) to have a token.

## List users
```
c := okta.New(apiToken, "organisation")
users, err := c.User.GetUsers(context.Background())
```

See the [documentation](https://godoc.org/github.com/arkan/okta) for all the available commands.

## Licence
[MIT](./LICENSE)

