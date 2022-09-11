package response

const (
	Success           Code = "AU0000"
	ServerError       Code = "AU0001"
	BadRequest        Code = "AU0002"
	InvalidRequest    Code = "AU0004"
	Failed            Code = "AU0073"
	Pending           Code = "AU0050"
	InvalidInputParam Code = "AU0032"
	DuplicateUser     Code = "AU0033"
	NotFound          Code = "AU0034"

	Unauthorized   Code = "AU0502"
	Forbidden      Code = "AU0503"
	GatewayTimeout Code = "AU0048"
)

type Code string

var codeMap = map[Code]string{
	Success:           "success",
	Failed:            "failed",
	Pending:           "pending",
	BadRequest:        "bad or invalid request",
	Unauthorized:      "Unauthorized Token",
	GatewayTimeout:    "Gateway Timeout",
	ServerError:       "Internal Server Error",
	InvalidInputParam: "Other invalid argument",
	DuplicateUser:     "duplicate user",
	NotFound:          "Not found",
}

func (c Code) GetStatus() string {
	switch c {
	case Success:
		return "SUCCESS"

	default:
		return "FAILED"
	}
}

func (c Code) GetMessage() string {
	return codeMap[c]
}

func (c Code) GetVersion() string {
	return "1"
}
