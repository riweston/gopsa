package main

import (
	"log"

	"github.com/zalando/go-keyring"
)

func storeConfig(appconfig AppConfig) {
	service := appconfig.endpoint
	user := appconfig.username
	password := appconfig.password

	// set password
	err := keyring.Set(service, user, password)
	if err != nil {
		log.Fatal(err)
	}

	// get password
	secret, err := keyring.Get(service, user)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(secret)
}
