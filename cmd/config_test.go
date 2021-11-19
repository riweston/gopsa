package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
)

func testSetConfigEndpoint(t *testing.T) {

	url := "test.com"
	setConfigEndpoint(url)

	got := viper.GetString("endpoint")
	want := url

	if got != want {
		t.Errorf("got %q want %q given, %v", got, want, url)
	}
}

func testSetConfigUsername(t *testing.T) {

	username := "user@test.com"
	setConfigUsername(username)

	got := viper.GetString("username")
	want := username

	if got != want {
		t.Errorf("got %q want %q given, %v", got, want, username)
	}
}

func testSetConfigPassword(t *testing.T) {

	url := "test.com"
	username := "user@test.com"
	password := "myduMMypa55w0$d"

	defer keyring.Delete(("gopsa." + url + ".password"), username)
	setConfigPassword(url, username, password)
	key := viper.GetString("keychainPassword")

	got, _ := keyring.Get(key, username)
	want := password

	if got != want {
		t.Errorf("got %q want %q given, %v", got, want, password)
	}
}

func testSetConfigToken(t *testing.T) {

	url := "test.com"
	username := "user@test.com"
	token := "myduMMyt0k$n"

	defer keyring.Delete(("gopsa." + url + ".token"), username)
	setConfigToken(url, username, token)
	key := viper.GetString("keychaintoken")

	got, _ := keyring.Get(key, username)
	want := token

	if got != want {
		t.Errorf("got %q want %q given, %v", got, want, token)
	}
}
