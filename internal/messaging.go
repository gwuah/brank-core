package internal

import "net/http"

func PublishMessageIntoKafka(req MessageRequest, e EventStore) BrankResponse {
	go func() {
		e.Publish(GenerateTopic("validate_login"), []byte(req.Message))
	}()

	return BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: BrankMeta{
			Data:    map[string]interface{}{},
			Message: "message received",
		},
	}

}
