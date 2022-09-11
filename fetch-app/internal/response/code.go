package response

const (
	Success           Code = "FE0000"
	ServerError       Code = "FE0001"
	BadRequest        Code = "FE0002"
	InvalidRequest    Code = "FE0004"
	Failed            Code = "FE0073"
	Pending           Code = "FE0050"
	InvalidInputParam Code = "FE0032"
	DuplicateUser     Code = "FE0033"
	NotFound          Code = "FE0034"

	Unauthorized   Code = "FE0502"
	Forbidden      Code = "FE0503"
	GatewayTimeout Code = "FE0048"
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
