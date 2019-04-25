package http

var (
	InternalServerErrorCode    = "internal-server-error"
	InternalServerErrorMessage = "The server encountered an internal error or misconfiguration and was unable to complete your request."
)

type errorResponse struct {
	Status int
	Data   map[string]interface{}
}

// TODO: acquire/release from pool
func newErrorResponse(status int) *errorResponse {
	return &errorResponse{
		Status: status,
		Data:   make(map[string]interface{}),
	}
}

func (response *errorResponse) SetParam(name string, value interface{}) {
	switch name {
	case "statusCode":
		if v, ok := value.(int); ok {
			response.Status = v
		}
	default:
		response.Data[name] = value
	}
}
