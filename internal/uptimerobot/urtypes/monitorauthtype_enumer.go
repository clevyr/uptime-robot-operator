// Code generated by "enumer -type MonitorAuthType -json -trimprefix Auth"; DO NOT EDIT.

package urtypes

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _MonitorAuthTypeName = "BasicDigest"

var _MonitorAuthTypeIndex = [...]uint8{0, 5, 11}

const _MonitorAuthTypeLowerName = "basicdigest"

func (i MonitorAuthType) String() string {
	i -= 1
	if i >= MonitorAuthType(len(_MonitorAuthTypeIndex)-1) {
		return fmt.Sprintf("MonitorAuthType(%d)", i+1)
	}
	return _MonitorAuthTypeName[_MonitorAuthTypeIndex[i]:_MonitorAuthTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _MonitorAuthTypeNoOp() {
	var x [1]struct{}
	_ = x[AuthBasic-(1)]
	_ = x[AuthDigest-(2)]
}

var _MonitorAuthTypeValues = []MonitorAuthType{AuthBasic, AuthDigest}

var _MonitorAuthTypeNameToValueMap = map[string]MonitorAuthType{
	_MonitorAuthTypeName[0:5]:       AuthBasic,
	_MonitorAuthTypeLowerName[0:5]:  AuthBasic,
	_MonitorAuthTypeName[5:11]:      AuthDigest,
	_MonitorAuthTypeLowerName[5:11]: AuthDigest,
}

var _MonitorAuthTypeNames = []string{
	_MonitorAuthTypeName[0:5],
	_MonitorAuthTypeName[5:11],
}

// MonitorAuthTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func MonitorAuthTypeString(s string) (MonitorAuthType, error) {
	if val, ok := _MonitorAuthTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _MonitorAuthTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to MonitorAuthType values", s)
}

// MonitorAuthTypeValues returns all values of the enum
func MonitorAuthTypeValues() []MonitorAuthType {
	return _MonitorAuthTypeValues
}

// MonitorAuthTypeStrings returns a slice of all String values of the enum
func MonitorAuthTypeStrings() []string {
	strs := make([]string, len(_MonitorAuthTypeNames))
	copy(strs, _MonitorAuthTypeNames)
	return strs
}

// IsAMonitorAuthType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i MonitorAuthType) IsAMonitorAuthType() bool {
	for _, v := range _MonitorAuthTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for MonitorAuthType
func (i MonitorAuthType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for MonitorAuthType
func (i *MonitorAuthType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("MonitorAuthType should be a string, got %s", data)
	}

	var err error
	*i, err = MonitorAuthTypeString(s)
	return err
}
