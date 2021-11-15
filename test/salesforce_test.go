package gopsa

import (
	"testing"
)

func TestSetConfig(t *testing.T) {
	const (
		endpoint string = "https://cloudreach.my.salesforce.com/"
		username string = "user@cloudreach.com"
		password string = "password"
		token    string = "token"
	)
	t.Run("SetConfig", func(t *testing.T) {
		got := setConfig(endpoint, username, password, token)
		want := AppConfig{
			endpoint: "https://cloudreach.my.salesforce.com/",
			username: "user@cloudreach.com",
			password: "password",
			token:    "token",
		}
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}
