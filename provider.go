package main

import (
	b64 "encoding/base64"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	HostEnv     string = "D42_HOST"
	UsernameEnv string = "D42_USER"
	PasswordEnv string = "D42_PASS"
)

var (
	data       string = UsernameEnv + ":" + UsernameEnv
	AuthString string = "Basic " + b64.StdEncoding.EncodeToString([]byte(data))
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"device": resourceDevice(),
		},
	}
}
