package internal

import (
	"brank/internal/repository"
	"errors"
	"log"
	"net/http"

	"gorm.io/gorm"
)

func ExchageContractCode(req ExchangeContractCode, r repository.Repo, c *Config, e EventStore) BrankResponse {
	link, err := r.Link.FindByCode(req.Code)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return BrankResponse{
			Error: true,
			Code:  http.StatusUnauthorized,
			Meta: BrankMeta{
				Data:    nil,
				Message: "unauthorized request",
			},
		}
	}

	if err != nil {
		log.Println("link-FindByCode failed. err", err)
		return BrankResponse{
			Error: true,
			Code:  http.StatusBadRequest,
			Meta: BrankMeta{
				Data:    nil,
				Message: "account already exists. Wanna login?",
			},
		}
	}

	return BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: BrankMeta{
			Data: map[string]interface{}{
				"link": link,
			},
			Message: "exchange was successful",
		},
	}
}

func LinkAccount(req LinkAccountRequest, r repository.Repo, c *Config, e EventStore) BrankResponse {
	return BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: BrankMeta{
			Data:    nil,
			Message: "client created successfully",
		},
	}
}
