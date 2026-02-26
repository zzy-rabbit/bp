package model

import (
	"encoding/json"
	"github.com/zzy-rabbit/xtools/xerror"
)

type Network struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Header struct {
	Authorization string `json:"authorization"`
}

type HttpResponse struct {
	xerror.IError
	Data any `json:"data"`
}

func (r *HttpResponse) MarshalJSON() ([]byte, error) {
	if r.IError == nil {
		r.IError = xerror.ErrSuccess
	}
	type response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}
	resp := response{
		Code:    r.Code(),
		Message: r.Message(),
		Data:    r.Data,
	}
	if resp.Data == nil {
		resp.Data = json.RawMessage("{}")
	}
	return json.Marshal(resp)
}

func (r *HttpResponse) UnmarshalJSON(data []byte) error {
	type response struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	var resp response
	err := json.Unmarshal(data, &resp)
	if xerror.Error(err) {
		return err
	}
	r.IError = xerror.New(resp.Code, resp.Message)
	r.Data = resp.Data
	return nil
}
