package domain

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type Cursor struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func EncodeCursor(c Cursor) (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to encode cursor: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func DecodeCursor(encoded string) (Cursor, error) {
	var c Cursor
	bytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return Cursor{}, fmt.Errorf("failed to decode base64 cursor: %w", err)
	}
	if err := json.Unmarshal(bytes, &c); err != nil {
		return Cursor{}, fmt.Errorf("failed to unmarshal cursor json: %w", err)
	}
	return c, nil
}
