package main

import (
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestGetAcquirerPermissions(t *testing.T) {
	tests := []struct{
		name string
		acquirers []string
		expected []string
	}{
		{"No acquirers", nil, nil},
		{"One acquirer", []string{"NI"}, []string{"NI"}},
		{"Multiple acquirers", []string{"NI", "Standard Bank", "Automated Test Acquirer"}, []string{"NI", "Standard Bank", "Automated Test Acquirer"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := url.Values{}
			for _, acquirer := range tt.acquirers {
				data.Add("Acquirers[]", acquirer)
			}

			req, err := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

			got := getAcquirerPermissions(req)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("getAcquirerPermissions() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetResultCodeDescription(t *testing.T) {
	tests := []struct{
		name string
		resultCode string
		expected string
	}{
		{"Unrecognised code", "Not a real code", "Not a real code - Declined"},
		{"00 - Approved", "00", "00 - Approved"},
		{"01 - Call issuer", "01", "01 - Call issuer"},
		{"02 - Call issuer", "02", "02 - Call issuer"},
		{"03 - Invalid Merchant", "03", "03 - Invalid Merchant"},
		//No code 04
		{"05 - Do Not Honor", "05", "05 - Do Not Honor"},
		//No codes 06-11
		{"12 - Invalid Txn", "12", "12 - Invalid Txn"},
		{"13 - Invalid Amount", "13", "13 - Invalid Amount"},
		{"14 - Invalid Card", "14", "14 - Invalid Card"},
		{"15 - <No string returned>", "15", ""},
		//No codes 16-18
		{"19 - Retry the Txn", "19", "19 - Retry the Txn"},
		//No codes 20-24
		{"25 - Declined", "25", "25 - Declined"},
		//No codes 26-29
		{"30 - Format Error", "30", "30 - Format Error"},
		{"31 - Unsupported Txn", "31", "31 - Unsupported Txn"},
		//No codes 32-40
		{"41 - Please Call - LC", "41", "41 - Please Call - LC"},
		//No code 42
		{"43 - Please Call - CC", "43", "43 - Please Call - CC"},
		//No codes 44-50
		{"51 - Declined", "51", "51 - Declined"},
		//No codes 52-53
		{"54 - Expired Card", "54", "54 - Expired Card"},
		{"55 - Incorrect Pin", "55", "55 - Incorrect Pin"},
		//No codes 56-57
		{"58 - Txn not allowed", "58", "58 - Txn not allowed"},
		//No codes 59-64
		{"65 - Perform Contact Txn", "65", "65 - Perform Contact Txn"},
		//No codes 66-73
		{"74 - Txn Unavailable", "74", "74 - Txn Unavailable"},
		//No codes 75-88
		{"89 - Invalid terminal", "89", "89 - Invalid terminal"},
		//No code 90
		{"91 - Auth timed out", "91", "91 - Auth timed out"},
		//No codes 91-93
		{"94 - Duplicate TXN", "94", "94 - Duplicate TXN"},
		{"95 - Txn Cancelled", "95", "95 - Txn Cancelled"},
		{"96 - Declined", "96", "96 - Declined"},
		{"97 - Signature Mismatch", "97", "97 - Signature Mismatch"},
		{"98 - Card Removed", "98", "98 - Card Removed"},
		{"99 - Comms Error", "99", "99 - Comms Error"},
		{"<empty string> - Unknown Reason", "", "Unknown Reason"},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetResultCodeDescription(tt.resultCode)
			if got != tt.expected {
				t.Errorf("GetResultCodeDescription() = %v, want %v", got, tt.expected)
			}
		})
	}
}
