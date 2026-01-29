package fuzzy

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type TLMethod struct {
	Name string `json:"name"`
}

type TLData struct {
	Methods []TLMethod `json:"methods"`
}

// LoadTLMethods loads method names from a JSON file path.
func LoadTLMethods(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open TL methods file: %w", err)
	}
	defer f.Close()

	return LoadTLMethodsFromReader(f)
}

// LoadTLMethodsFromReader loads method names from any io.Reader.
// This allows loading from files, HTTP responses, or embedded data.
func LoadTLMethodsFromReader(r io.Reader) ([]string, error) {
	var data TLData
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode TL methods JSON: %w", err)
	}

	out := make([]string, len(data.Methods))
	for i, m := range data.Methods {
		out[i] = m.Name
	}

	return out, nil
}
