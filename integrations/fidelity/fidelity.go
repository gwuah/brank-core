package fidelity

import (
	"brank/core/axios"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ledongthuc/pdf"
)

func New(a *axios.Axios) *Integration {
	return &Integration{
		axios:                      a,
		loginEndpoint:              "https://retailibank.fidelitybank.com.gh/mmrib/auth/init",
		verifyOtpEndpoint:          "https://retailibank.fidelitybank.com.gh/mmrib/auth/continue",
		balanceEndpoint:            "https://retailibank.fidelitybank.com.gh/mmrib/account/balance",
		recentTransactionsEndpoint: "https://retailibank.fidelitybank.com.gh/mmrib/account/statement/%d",
		statementEndpoint:          "https://retailibank.fidelitybank.com.gh/mmrib/file/statement/%d?start=%s&end=%s",
	}
}

func (f *Integration) SetBearerToken(token string) {
	f.axios.SetBearerToken(token)
}

func (f *Integration) VerifyLogin(username, password string) (bool, *HTTPResponse, error) {
	reqBody := map[string]string{
		"phoneNumber":        username,
		"pin":                password,
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
		var response HTTPResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return false, nil, fmt.Errorf("failed to unmarshal fidelity reponse. err: %v", err)
		}
		return true, &response, nil
	}

	return false, nil, nil

}

func (f *Integration) VerifyOtp(otp string) (bool, *HTTPResponse, error) {
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
		var response HTTPResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return false, nil, fmt.Errorf("failed to unmarshal fidelity reponse. err: %v", err)
		}
		return true, &response, nil
	}

	return false, nil, nil
}

func (f *Integration) GetBalance() (bool, *HTTPResponse, error) {
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
		var response HTTPResponse
		err := json.Unmarshal(body, &response)
		if err != nil {
			return false, nil, fmt.Errorf("failed to unmarshal fidelity reponse. err: %v", err)
		}
		return true, &response, nil
	}

	return false, nil, nil
}

func (f *Integration) DownloadStatement(accountId int64, start, end string) (bool, []byte, error) {
	var body []byte
	res, err := f.axios.Get(context.Background(), fmt.Sprintf(f.statementEndpoint, accountId, start, end))

	if err != nil {
		return false, body, err
	}

	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return false, body, err
	}

	if res.StatusCode == 401 || res.StatusCode == 403 {
		return false, body, nil
	}

	return true, body, nil
}

func (f *Integration) ProcessPDF(body []byte) (*TransactionTree, error) {
	var tree TransactionTree

	bytesWithReader := bytes.NewReader(body)
	rows := [][]string{}
	length := 0

	r, err := pdf.NewReader(bytesWithReader, int64(len(body)))
	if err != nil {
		return nil, err
	}

	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		pdfRows, _ := p.GetTextByRow()
		for _, row := range pdfRows {
			var array []string
			for _, word := range row.Content {
				if string(word.S) == "Not for official Use" {
					continue
				}

				if strings.Contains(string(word.S), "Please review your statements and draw your Relationship Manager's") {
					continue
				}

				if strings.Contains(string(word.S), "shall be deemed as correct and shall be binding on you") {
					continue
				}

				if strings.Contains(string(word.S), "Fidelity Bank Ghana") {
					continue
				}

				if word.S == "" {
					continue
				}

				array = append(array, string(word.S))
			}

			if len(array) > 0 {
				rows = append(rows, array)
				length += 1
			}
		}
	}

	commenceConstruction := false

	for i := 0; i < length; i++ {
		iRow := rows[i]

		// if the row length is 0
		// iRow[0] == ""  might introduce bugs but let's keep it there for now
		if len(iRow) == 0 || iRow[0] == "" {
			continue
		}

		if !commenceConstruction {
			// if we have not seen the row that starts with date, ignore and continue
			// else, you can start processsing the pdf rows
			if strings.ToLower(iRow[0]) == "date" {
				commenceConstruction = true
			}
			continue
		}

		// main row, overflow of description...
		processingQueue := [][]string{}
		processingQueue = append(processingQueue, iRow)

		for j := i + 1; j < length; j++ {
			jRow := rows[j]

			// if we start sliding our window and we see a date, it means we have encountered
			// a new row, and as such we exit and go back to "i"
			if ok, _ := isValidDate(jRow[0]); ok {
				i = j - 1
				break
			}

			// if we see the string below, it means we've completed parsing our pdf
			// so we set i and j at max value and then we continue
			// by continueing, we'll exit both loops
			if strings.Contains(jRow[0], "No of Debits") {
				i = length - 1
				j = length - 1
				tree.Summary = jRow
				continue
			}

			// if none of our 2 conditions run, it means we are currently amassing
			// overflow description, so we just continue sliding our window
			processingQueue = append(processingQueue, jRow)
		}

		// once we exit the sliding window loop, we pass all the terms in a particular window
		// and pass it to create transaction that combines them and creates a transaction
		tree.Transactions = append(tree.Transactions, CreateTransaction(processingQueue))

	}

	return &tree, err

}
