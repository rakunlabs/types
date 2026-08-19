package types

import (
	"encoding/json"
	"testing"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		args    [][]byte
		wantErr bool
	}{
		{
			name: "Test Time UnmarshalJSON",
			args: [][]byte{
				[]byte(`2020-01-01`),
				[]byte(`"2020-01-01 00:00:00"`),
				[]byte(`2020-01-01T00:00:00Z`),
				[]byte(`2006-01-02T15:04:05Z`),
				[]byte(`2020-01-01T00:00:00.000Z`),
				[]byte(`2025-01-31T09:41:17Z`),
				[]byte(`2025-01-31T10:41:55+01:00`),
				[]byte(`2025-01-31T09:43:00.3Z`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tr Time
			for _, arg := range tt.args {
				if err := tr.UnmarshalJSON(arg); (err != nil) != tt.wantErr {
					t.Errorf("Time.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestTime_UnmarshalJSON_Struct(t *testing.T) {
	value := []byte(`{}`)

	type TimeStruct struct {
		Time Time `json:"time"`
	}

	var ts TimeStruct

	if err := json.Unmarshal(value, &ts); err != nil {
		t.Errorf("Time.UnmarshalJSON() error = %v", err)
	}

	value = []byte(`{"time": null}`)
	type TimeStruct2 struct {
		Time Null[Time] `json:"time"`
	}

	var ts2 TimeStruct2

	if err := json.Unmarshal(value, &ts2); err != nil {
		t.Errorf("Time.UnmarshalJSON() error = %v", err)
	}

	var ts3 TimeStruct
	if err := json.Unmarshal(value, &ts3); err != nil {
		t.Errorf("Time.UnmarshalJSON() error = %v", err)
	}
}
