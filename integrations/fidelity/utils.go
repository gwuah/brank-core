package fidelity

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func CreateTransaction(values [][]string) Transaction {
	main := values[0]
	trans := Transaction{
		Date:        main[0],
		Description: []string{main[1]},
		ValueDate:   main[2],
		Debit:       main[3],
		Credit:      main[4],
		Balance:     main[5],
	}

	for _, value := range values[1:] {
		trans.Description = append(trans.Description, value[0])
	}

	return trans
}

func ConvertToFloat(v string) float64 {
	res := []string{}
	for _, c := range v {
		char := fmt.Sprintf("%c", c)
		if char == "," {
			continue
		}
		res = append(res, char)
	}

	value, err := strconv.ParseFloat(strings.Join(res, ""), 64)

	if err != nil {
		fmt.Println("err occured whiles converting string -> int. err:", err)
		return 0
	}

	return value
}

func ConvertToInt(v string) int64 {
	res := []string{}
	for _, c := range v {
		char := fmt.Sprintf("%c", c)
		if char == "," || char == "." {
			continue
		}
		res = append(res, char)
	}

	value, err := strconv.ParseInt(strings.Join(res, ""), 10, 64)

	if err != nil {
		fmt.Println("err occured whiles converting string -> int. err:", err)
		return 0
	}

	return value
}

func IsInt(v string) bool {
	if _, err := strconv.ParseInt(v, 10, 64); err == nil {
		return true
	}
	return false
}

func logRow(value []string) {
	fmt.Println("row -> ", "['"+strings.Join(value, `','`)+`']`)
}

func logSubRow(value []string) {
	fmt.Println("subrow ---------> ", "['"+strings.Join(value, `','`)+`']`)
}

func isValidDate(date string) (bool, error) {
	matched, err := regexp.MatchString(`\d{2}-\d{2}-\d{4}`, date)
	return matched, err
}
