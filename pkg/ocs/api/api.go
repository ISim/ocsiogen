package api

const (
	headerOCSAppID  = "X-Huwwi-Application-Id"
	headerOCSAppKey = "X-Huwwi-API-Key"
)

const (
	apiResultStatusOK = "OK"
)

const (
	apiCreateCustomer = "createCustomer"
)

var url4actions = map[string]string{
	apiCreateCustomer: "casm/accounts/create",
}
