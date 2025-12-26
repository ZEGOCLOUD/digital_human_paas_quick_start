package handler

import "zego-digital-human-server/pkg/response"

// buildCommonResponseFromAPIData 将第三方接口返回转换为统一响应格式
func buildCommonResponseFromAPIData(apiData map[string]interface{}) response.CommonResponse {
	code := 500
	if v, ok := apiData["Code"].(float64); ok {
		code = int(v)
	}

	message := ""
	if v, ok := apiData["Message"].(string); ok {
		message = v
	}

	data := apiData["Data"]
	if data == nil {
		// 没有 Data 字段时直接返回完整数据，避免客户端空指针
		data = apiData
	}

	return response.CommonResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
