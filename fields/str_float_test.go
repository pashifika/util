// Package fields
package fields

import (
	"reflect"
	"testing"
)

func TestStrFloat_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		s       StrFloat
		want    []byte
		wantErr bool
	}{
		{name: "test 01", s: 3.1415926535, want: []byte("\"3.1415927\""), wantErr: false},
		{name: "test 02", s: -3.1415926535, want: []byte("\"-3.1415927\""), wantErr: false},
		{name: "test 03", s: -3.1415, want: []byte("\"-3.1415\""), wantErr: false},
		{name: "test 04", s: 0, want: []byte("\"0\""), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrFloat_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		s       StrFloat
		args    args
		wantErr bool
	}{
		{name: "test 01", s: 3.1415926535, args: args{data: []byte("\"3.1415926535\"")}, wantErr: false},
		{name: "test 02", s: -3.1415926535, args: args{data: []byte("\"-3.1415926535\"")}, wantErr: false},
		{name: "test 03", s: -3.1415, args: args{data: []byte("\"-3.1415\"")}, wantErr: false},
		{name: "test 04", s: 0, args: args{data: []byte("\"0\"")}, wantErr: false},
		{name: "test error", s: 0, args: args{data: []byte("\"3.1s\"")}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStrFloat64_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		s       StrFloat64
		want    []byte
		wantErr bool
	}{
		{name: "test 01", s: 3.1415926535, want: []byte("\"3.1415926535\""), wantErr: false},
		{name: "test 02", s: -3.1415926535, want: []byte("\"-3.1415926535\""), wantErr: false},
		{name: "test 03", s: -3.1415, want: []byte("\"-3.1415\""), wantErr: false},
		{name: "test 04", s: 0, want: []byte("\"0\""), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrFloat64_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		s       StrFloat64
		args    args
		wantErr bool
	}{
		{name: "test 01", s: 3.1415926535, args: args{data: []byte("\"3.1415926535\"")}, wantErr: false},
		{name: "test 02", s: -3.1415926535, args: args{data: []byte("\"-3.1415926535\"")}, wantErr: false},
		{name: "test 03", s: -3.1415, args: args{data: []byte("\"-3.1415\"")}, wantErr: false},
		{name: "test 04", s: 0, args: args{data: []byte("\"0\"")}, wantErr: false},
		{name: "test error", s: 0, args: args{data: []byte("\"3.1s\"")}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
