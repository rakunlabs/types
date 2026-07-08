package types

import "testing"

func TestMap_Scan(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		m       *Map[any]
		args    args
		wantErr bool
	}{
		{
			name:    "nil",
			m:       &Map[any]{},
			args:    args{value: nil},
			wantErr: false,
		},
		{
			name:    "null",
			m:       &Map[any]{},
			args:    args{value: []byte("null")},
			wantErr: false,
		},
		{
			name:    "empty",
			m:       &Map[any]{},
			args:    args{value: []byte("{}")},
			wantErr: false,
		},
		{
			name:    "invalid",
			m:       &Map[any]{},
			args:    args{value: []byte("invalid")},
			wantErr: true,
		},
		{
			name:    "unsupported",
			m:       &Map[any]{},
			args:    args{value: 42},
			wantErr: true,
		},
		{
			name:    "nested",
			m:       &Map[any]{},
			args:    args{value: []byte(`{"key":{"nested":"value"}}`)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Map.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
