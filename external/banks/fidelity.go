package external

import (
	"brank/core/utils"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
)

type FidelityVerifyLoginResponse struct {
	Token              string `json:"token"`
	SecondFactorMethod string `json:"secondFactorMethod"`
	Timestamp          string `json:"timestamp"`
	Status             int    `json:"status"`
	Error              string `json:"error"`
	Path               string `json:"path"`
}

func NewFidelity(a utils.Axios) *BankIntegration {
	return &BankIntegration{
		axios:    a,
		endpoint: "https://retailibank.fidelitybank.com.gh/mmrib/auth/init",
	}
}

func (f *BankIntegration) verifyLogin(username, password string) (bool, *FidelityVerifyLoginResponse, error) {
	reqBody := map[string]string{
		"phoneNumber":        username,
		"password":           password,
		"secondFactorMethod": "SMS-OTP",
	}

	res, err := f.axios.Post(context.Background(), f.endpoint, reqBody)

	if err != nil {
		return false, nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return false, nil, err
	}

	if res.StatusCode == 401 || res.StatusCode == 403 {
		return false, nil, nil
	} else if res.StatusCode == 200 {
		var response FidelityVerifyLoginResponse
		err := json.Unmarshal(body, &response)
		if err != nil {
			return false, nil, errors.New("failed to unmarshal fidelity_verify_login reponse")
		}
		return true, &response, nil
	} else {
		return false, nil, nil
	}

}

func (f *BankIntegration) verifyOtp() {

}
