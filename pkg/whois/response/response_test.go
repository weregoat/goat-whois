package response

import (
	"errors"
	"testing"
)

func TestResponse_IsValid(t *testing.T) {
	validResponses := []Response{
		{"something", "SE", "", []byte("body"), nil, []string{"SE"}},
		{"something", "", "192.168.0.0/16", []byte("body"), nil, []string{}},
		{"something", "SE", "192.168.0.0/16", []byte("body"), nil, []string{"SE"}},
	}
	for _,response := range validResponses {
		if ! response.IsValid() {
			t.Errorf("Valid response %q marked as invalid", response)
		}
	}

	invalidResponses := []Response{
		{"", "SE", "", []byte("body"), nil, []string{}},
		{"something", "", "", []byte("body"), nil, []string{}},
		{Resource: "something", CountryCode: "SE", CIDR: "192.168.0.0/16", Error: nil, CountryCodes: []string{"SE"}},
		{"something", "SE", "", []byte("body"), errors.New("error"), []string{"SE"}},
	}
	for _,response := range invalidResponses {
		if response.IsValid() {
			t.Errorf("invalid response %q marked as valid", response)
		}
	}
}

func TestGetCIDR(t *testing.T) {
	payloads := map[string]string {
		"192.168.13.0/24": "CIDR: 192.168.13.0/24 ",
		"192.168.12.0/24": "inetnum: 192.168.12.0 - 192.168.12.255 ",
		"69.175.0.0/17": "NetRange:       69.175.0.0 - 69.175.127.255\nCIDR:           69.175.0.0/17 ",
		"2a00:1450:4000::/37": "inet6num:       2a00:1450:4000::/37 ",
		"2001:470::/32": "CIDR:           2001:470::/32 ",
		"179.6.0.0/16": "inetnum:     179.6/16 ",
		"202.214.194.128/25": "a. [Network Number]             202.214.194.128/25 ",
	}
	for expected, payload := range payloads {
		cidr,err := GetCIDR([]byte(payload))
		if err != nil {
			t.Error(err)
		}
		if cidr != expected {
			t.Errorf("failed to correctly extract CIDR from %s; expected %s, got %s", payload, expected, cidr)
		}
	}
}

func TestGetCountries(t *testing.T) {
	payloads := map[string]string {
		"PA" : "Registrant Country: PA ",
		"BB" : "Registrant Country: PA\nblablah\nRegistrant Country: BB\n",
		"GG" : "Country: GG\n",
		"PP" : "country: PP ",
	}
	for expected, payload := range payloads {
		isoCodes,err := GetCountries([]byte(payload))
		if err != nil {
			t.Error(err)
		}
		if isoCodes[len(isoCodes)-1] != expected {
			t.Errorf("failed to correctly extract Country from %s; expected %s, got %s", payload, expected, isoCodes)
		}
	}
}
