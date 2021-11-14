package gopsa_test

import (
	"testing"
)

func TestSetConfig(t *testing.T) {
	const (
	)
	t.Run("SetConfig", func(t *testing.T) {
		got := setConfig(endpoint, username, password, token)
		want := ""
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}
