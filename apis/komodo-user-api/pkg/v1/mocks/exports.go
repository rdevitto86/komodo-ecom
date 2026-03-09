package mocks

import (
	_ "embed"
	"encoding/json"
)

//go:embed user_full.json
var userFullJSON []byte

//go:embed user_minimal.json
var userMinimalJSON []byte

//go:embed user_basic.json
var userBasicJSON []byte

//go:embed address.json
var addressJSON []byte

//go:embed address_input.json
var addressInputJSON []byte

//go:embed payment_method.json
var paymentMethodJSON []byte

//go:embed payment_method_input.json
var paymentMethodInputJSON []byte

//go:embed preferences.json
var preferencesJSON []byte

type mocksExport struct {
	UserFull           any
	UserMinimal        any
	UserBasic          any
	Address            any
	AddressInput       any
	PaymentMethod      any
	PaymentMethodInput any
	Preferences        any
}

var Mocks = mocksExport{
	UserFull:           mustParse(userFullJSON),
	UserMinimal:        mustParse(userMinimalJSON),
	UserBasic:          mustParse(userBasicJSON),
	Address:            mustParse(addressJSON),
	AddressInput:       mustParse(addressInputJSON),
	PaymentMethod:      mustParse(paymentMethodJSON),
	PaymentMethodInput: mustParse(paymentMethodInputJSON),
	Preferences:        mustParse(preferencesJSON),
}

func mustParse(data []byte) any {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		panic("mocks: failed to parse JSON: " + err.Error())
	}
	return v
}
