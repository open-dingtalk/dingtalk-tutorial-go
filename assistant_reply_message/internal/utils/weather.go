package utils

func GetWeather() map[string]any {
	result := make(map[string]any)
	result["text"] = "晴天"
	result["temperature"] = 22
	result["humidity"] = 65
	result["wind_direction"] = "东南风"
	return result
}
