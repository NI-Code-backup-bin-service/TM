package main

import "crypto/rand"

func GenerateCSRFToken() []byte {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		logging.Error(err)
	}
	return token
}
