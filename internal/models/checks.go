package models

import "encoding/json"

func IsJSON(str string) bool {
	return json.Valid([]byte(str))
}