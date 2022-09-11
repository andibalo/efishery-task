package response

type Wrapper struct {
	ResponseCode Code        `json:"code"`
	Status       string      `json:"status"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
}

func NewResponse(code Code, data interface{}) *Wrapper {
	return &Wrapper{
		ResponseCode: code,
		Status:       code.GetStatus(),
		Message:      code.GetMessage(),
		Data:         data,
	}
}

func (w *Wrapper) SetResponseMessage(customMessage string) {
	w.Message = customMessage
}
