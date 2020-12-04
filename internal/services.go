package internal

import (
	"net/http"
)

func HandleMessagePost(req MessageRequest, e EventStore) (int, map[string]interface{}, error) {
	go func() {
		e.Publish("validate_login", []byte(req.Message))
	}()

	return http.StatusOK, map[string]interface{}{
		"status": "ok",
	}, nil

}
