package hasher

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

// HashJSONObject convert interface to json object and hash
func HashJSONObject(object interface{}) (string, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return "", err
	}

	hashed := md5.Sum(data)

	return hex.EncodeToString(hashed[:]), nil
}
