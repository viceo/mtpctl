package dto

import (
	"encoding/json"
	"os"
)

type CommandResponse struct {
	HasError   bool   `json:"hasError"`
	Error      string `json:"error,omitempty"`
	StackTrace string `json:"stackTrace,omitempty"`
	Data       any    `json:"data,omitempty"`
	Hostname   string `json:"hostname"`
}

func (x *CommandResponse) Print() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	x.Hostname = hostname
	if err := json.NewEncoder(os.Stdout).Encode(x); err != nil {
		panic(err)
	}
}
