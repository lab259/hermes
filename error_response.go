package hermes

import "sync"

var (
	InternalServerErrorCode    = "internal-server-error"
	InternalServerErrorMessage = "We encountered an internal error or misconfiguration and was unable to complete your request."

	NotFoundErrorCode    = "not-found"
	NotFoundErrorMessage = "We could not find the resource you requested."

	MethodNotAllowedErrorCode    = "method-not-allowed"
	MethodNotAllowedErrorMessage = "We believe that the used request method is inappropriate for the resource you requested."
)

var errorResponsePool = &sync.Pool{
	New: func() interface{} {
		return &errorResponse{
			Data: make(map[string]interface{}),
		}
	},
}

type errorResponse struct {
	defaultStatus int
	Status        int
	Data          map[string]interface{}
}

func acquireErrorResponse(status int) *errorResponse {
	r := errorResponsePool.Get().(*errorResponse)
	r.defaultStatus = status
	r.Status = status
	return r
}

func releaseErrorResponse(r *errorResponse) {
	r.reset()
	errorResponsePool.Put(r)
}

func (response *errorResponse) reset() {
	response.defaultStatus = 0
	response.Status = 0
	for key := range response.Data {
		delete(response.Data, key)
	}
}

func (response *errorResponse) SetParam(name string, value interface{}) {
	switch name {
	case "statusCode":
		if v, ok := value.(int); ok && response.Status == response.defaultStatus {
			response.Status = v
		}
	default:
		if _, found := response.Data[name]; !found {
			response.Data[name] = value
		}
	}
}
