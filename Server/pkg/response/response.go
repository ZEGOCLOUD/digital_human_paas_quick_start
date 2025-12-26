package response

// CommonResponse 通用响应格式
type CommonResponse struct {
	Code      int         `json:"Code"`
	Message   string      `json:"Message"`
	Data      interface{} `json:"Data"`
	RequestId string      `json:"RequestId,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Code:    0,
		Message: message,
		Data:    data,
	}
}

