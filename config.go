package gobuild

import (
	"encoding/json"
	"os"
)

func GetTargetFromJson(path string) (*Target, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result Target
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
