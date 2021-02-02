package external

import (
	"brank/core/utils"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
)

type FidelityAccount struct {
	AccountNumber string `json:"accountNumber"`
	Currency      string `json:"currency"`
	Id            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
}

type FidelityUser struct {
	UserUid        string            `json:"userUid"`
	UserType       string            `json:"userType"`
	Status         string            `json:"status"`
	Name           string            `json:"name"`
	MobileNumber   string            `json:"mobileNumber"`
	ActiveChannels []int             `json:"activeChannels"`
	Accounts       []FidelityAccount `json:"accounts"`
}

type FidelityBalance struct {
	Id      int    `json:"id"`
	Balance string `json:"balance"`
}

type FidelityHTTPResponse struct {
	// login endpoint
	Token              string `json:"token"`
	SecondFactorMethod string `json:"secondFactorMethod"`

	// api request error
	Error     string `json:"error"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Timestamp string `json:"timestamp"`

	// verify otp endpoint
	RemeberMeToken string       `json:"rememberMeToken"`
	User           FidelityUser `json:"user"`

	// balance endpoint
	Balances              []FidelityBalance `json:"balances"`
	TrackingId            string            `json:"trackingId"`
	TransactionStatusCode string            `json:"transactionStatusCode"`
}

func NewFidelity(a utils.Axios) *BankIntegration {
	return &BankIntegration{
		axios:                      a,
		loginEndpoint:              "https://retailibank.fidelitybank.com.gh/mmrib/auth/init",
		verifyOtpEndpoint:          "https://retailibank.fidelitybank.com.gh/mmrib/auth/continue",
		balanceEndpoint:            "https://retailibank.fidelitybank.com.gh/mmrib/account/balance",
		recentTransactionsEndpoint: "https://retailibank.fidelitybank.com.gh/mmrib/account/statement/0",
		statementEndpoint:          "https://retailibank.fidelitybank.com.gh/mmrib/file/statement/0?start=2021-02-01&end=2021-02-02",
	}
}

func (f *BankIntegration) SetBearerToken(token string) {
	f.axios.SetBearerToken(token)
}

func (f *BankIntegration) verifyLogin(username, password string) (bool, *FidelityHTTPResponse, error) {
	reqBody := map[string]string{
		"phoneNumber":        username,
		"password":           password,
		"secondFactorMethod": "SMS-OTP",
	}

	res, err := f.axios.Post(context.Background(), f.loginEndpoint, reqBody)

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
		var response FidelityHTTPResponse
		err := json.Unmarshal(body, &response)
		if err != nil {
			return false, nil, errors.New("failed to unmarshal fidelity reponse")
		}
		return true, &response, nil
	}

	return false, nil, nil

}

func (f *BankIntegration) verifyOtp(otp string) (bool, *FidelityHTTPResponse, error) {
	reqBody := map[string]string{
		"otp": otp,
	}

	res, err := f.axios.Post(context.Background(), f.verifyOtpEndpoint, reqBody)

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
		var response FidelityHTTPResponse
		err := json.Unmarshal(body, &response)
		if err != nil {
			return false, nil, errors.New("failed to unmarshal fidelity reponse")
		}
		return true, &response, nil
	}

	return false, nil, nil
}

func (f *BankIntegration) getBalance() (bool, *FidelityHTTPResponse, error) {
	res, err := f.axios.Get(context.Background(), f.balanceEndpoint)

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
		var response FidelityHTTPResponse
		err := json.Unmarshal(body, &response)
		if err != nil {
			return false, nil, errors.New("failed to unmarshal fidelity reponse")
		}
		return true, &response, nil
	}

	return false, nil, nil
}
