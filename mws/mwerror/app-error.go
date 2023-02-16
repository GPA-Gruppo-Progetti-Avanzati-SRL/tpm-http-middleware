package mwerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AppError interface {
	GetStatusCode() int
	Error() string
	Marshal(ct string) ([]byte, error)
	Sanitized() AppError
}

type AppErrorImpl struct {
	StatusCode  int    `yaml:"-" mapstructure:"-" json:"-"`
	ErrCode     string `json:"error-code,omitempty" yaml:"error-code,omitempty" mapstructure:"error-code,omitempty"`
	Ambit       string `json:"ambit,omitempty" yaml:"ambit,omitempty" mapstructure:"ambit,omitempty"`
	Step        string `yaml:"step,omitempty" mapstructure:"step,omitempty" json:"step,omitempty"`
	Text        string `json:"text,omitempty" yaml:"text,omitempty" mapstructure:"text,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`
	Message     string `yaml:"message,omitempty" mapstructure:"message,omitempty" json:"message,omitempty"`
	Ts          string `yaml:"timestamp,omitempty" mapstructure:"timestamp,omitempty" json:"timestamp,omitempty"`
}

func (ae AppErrorImpl) Error() string {
	var sv strings.Builder
	const sep = " - "
	if ae.StatusCode != 0 {
		sv.WriteString(fmt.Sprintf("status-code: %d"+sep, ae.StatusCode))
	}

	if ae.ErrCode != "" {
		sv.WriteString(fmt.Sprintf("error-code: %s"+sep, ae.ErrCode))
	}

	if ae.Ambit != "" {
		sv.WriteString(fmt.Sprintf("ambit: %s"+sep, ae.Ambit))
	}

	if ae.Step != "" {
		sv.WriteString(fmt.Sprintf("step: %s"+sep, ae.Step))
	}

	if ae.Text != "" {
		sv.WriteString(fmt.Sprintf("text: %s"+sep, ae.Text))
	}

	if ae.Description != "" {
		sv.WriteString(fmt.Sprintf("description: %s"+sep, ae.Description))
	}

	if ae.Message != "" {
		sv.WriteString(fmt.Sprintf("message: %s"+sep, ae.Message))
	}

	if ae.Ts != "" {
		sv.WriteString(fmt.Sprintf("timestamp: %s"+sep, ae.Ts))
	}

	return strings.TrimSuffix(sv.String(), sep)
}

func (ae AppErrorImpl) GetStatusCode() int {
	return ae.StatusCode
}

func (ae AppErrorImpl) GetMessage() string {
	return ae.Text
}

func (ae AppErrorImpl) Marshal(ct string) ([]byte, error) {

	if ct == "application/json" {
		b, err := json.Marshal(ae)
		return b, err
	}

	return nil, errors.New("app error cannot marshal to " + ct)
}

func (ae AppErrorImpl) Sanitized() AppError {

	nae := &AppErrorImpl{
		StatusCode: ae.StatusCode,
		Text:       ae.Text,
	}

	return nae
}

type AppErrorOption func(ae *AppErrorImpl)

func AppErrorWithStatusCode(sc int) AppErrorOption {
	return func(ae *AppErrorImpl) {
		ae.StatusCode = sc
	}
}

func AppErrorWithErrorCode(ec string) AppErrorOption {
	return func(ae *AppErrorImpl) {
		ae.ErrCode = ec
	}
}

func AppErrorWithAmbit(a string) AppErrorOption {
	return func(ae *AppErrorImpl) {
		ae.Ambit = a
	}
}

func AppErrorWithStep(s string) AppErrorOption {
	return func(ae *AppErrorImpl) {
		ae.Step = s
	}
}

func AppErrorWithText(t string) AppErrorOption {
	return func(ae *AppErrorImpl) {
		ae.Text = t
	}
}

func AppErrorWithMessage(m string) AppErrorOption {
	return func(ae *AppErrorImpl) {
		ae.Message = m
	}
}

func AppErrorWithDescription(d string) AppErrorOption {
	return func(ae *AppErrorImpl) {
		ae.Description = d
	}
}

func NewAppError(opts ...AppErrorOption) *AppErrorImpl {
	ae := &AppErrorImpl{StatusCode: http.StatusInternalServerError, Text: "internal server error", Ts: time.Now().Format(time.RFC3339Nano)}
	for _, o := range opts {
		o(ae)
	}

	return ae
}
