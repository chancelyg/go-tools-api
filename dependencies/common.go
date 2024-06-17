package dependencies

import "encoding/json"

func CalculateSize(data map[string]interface{}) int64 {
	jsonData, _ := json.Marshal(data)
	return int64(len(jsonData))
}
