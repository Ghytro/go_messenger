package jsonhelpers

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &ns.String)
	ns.Valid = (err == nil)
	return err
}

type NullInt struct {
	sql.NullInt64
}

func (ni NullInt) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

func (ni *NullInt) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &ni.Int64)
	ni.Valid = (err == nil)
	return err
}
