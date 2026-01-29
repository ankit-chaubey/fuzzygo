package fuzzy

import (
	"encoding/json"
	"os"
)

type TLMethod struct {
	Name string `json:"name"`
}

type TLData struct {
	Methods []TLMethod `json:"methods"`
}

func LoadTLMethods(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data TLData
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	out := make([]string, 0, len(data.Methods))
	for _, m := range data.Methods {
		out = append(out, m.Name)
	}

	return out, nil
}
