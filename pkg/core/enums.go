package core

import (
	"encoding/json"
	"errors"
	"fmt"
)

// CCFailReason represents the reason to fail a credit card.
type CCFailReason uint

const (
	// CCFailReason_Authorise represents an authorisation fail.
	CCFailReason_Authorise = iota + 1
	// CCFailReason_Capture represents a capture fail.
	CCFailReason_Capture
	// CCFailReason_Refund represents a refund fail.
	CCFailReason_Refund
	// CCFailReason_Void represents a void fail.
	CCFailReason_Void
)

// String returns the string representation of CCFailReason.
func (ccfr CCFailReason) String() string {
	return [...]string{"", "authorise fail", "capture fail", "refund fail", "void fail"}[ccfr]
}

var ccFailReasonToEnum = map[string]CCFailReason{
	"authorise fail": CCFailReason_Authorise,
	"capture fail":   CCFailReason_Capture,
	"refund fail":    CCFailReason_Refund,
	"void fail":      CCFailReason_Void,
}

func (ccfr *CCFailReason) Load(reason string) error {
	if reasonEnum, ok := ccFailReasonToEnum[reason]; ok {
		*ccfr = reasonEnum
		return nil
	}
	return fmt.Errorf("unknown reason to fail")
}

// UnmarshalJSON unmashals a quoted json string to the CCFaiLReason Enum
func (ccfr *CCFailReason) UnmarshalJSON(data []byte) error {
	var j string
	err := json.Unmarshal(data, &j)
	if err != nil {
		return err
	}

	result, ok := ccFailReasonToEnum[j]
	if !ok {
		return errors.New("couldn't find matching CCFailReason enum value")
	}

	*ccfr = result
	return nil
}

func (ccfr *CCFailReason) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	err := unmarshal(&j)
	if err != nil {
		return err
	}

	result, ok := ccFailReasonToEnum[j]
	if !ok {
		return errors.New("couldn't find matching CCFailReason enum value")
	}

	*ccfr = result
	return nil
}
