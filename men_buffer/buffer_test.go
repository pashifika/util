// Package men_buffer
/*
 * Version: 1.0.0
 * Copyright (c) 2021. Pashifika
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package men_buffer

import (
	"fmt"
	"io"
	"testing"
	"unicode/utf8"
)

func TestFakeIO_Write(t *testing.T) {
	type args struct {
		data []byte
		size int64
		loop int
	}
	tests := []struct {
		name    string
		args    args
		result  string
		wantN   int
		wantErr bool
	}{
		{name: "small buff", args: args{
			data: []byte("123"),
			size: 2,
			loop: 3,
		}, result: "123123123", wantN: 3, wantErr: false},
		{name: "default buff", args: args{
			data: []byte("123"),
			size: 0,
			loop: 3,
		}, result: "123123123", wantN: 3, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fio := NewBufferIO(tt.args.size)
			for i := 0; i < tt.args.loop; i++ {
				gotN, err := fio.Write(tt.args.data)
				if (err != nil) != tt.wantErr {
					t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotN != tt.wantN {
					t.Errorf("Write() gotN = %v, want %v", gotN, tt.wantN)
				}
			}
			str := fio.String()
			if fio.String() != tt.result {
				t.Errorf("buffer.string = %v, want %v", str, tt.result)
			}
			fmt.Printf("[%s] result: %s\n", tt.name, str)
		})
	}
}

func TestFakeIO_WriteAt(t *testing.T) {
	type args struct {
		data []byte
		size int64
		len  int
		loop int
	}
	tests := []struct {
		name    string
		args    args
		result  string
		wantN   int
		wantErr bool
	}{
		{name: "small buff", args: args{
			data: []byte("123"),
			size: 2,
			loop: 3,
			len:  3,
		}, result: "123123123", wantN: 3, wantErr: false},
		{name: "default buff", args: args{
			data: []byte("123"),
			size: 0,
			loop: 3,
			len:  3,
		}, result: "123123123", wantN: 3, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fio := NewBufferIO(tt.args.size)
			for i := 0; i < tt.args.loop; i++ {
				gotN, err := fio.WriteAt(tt.args.data, int64(i*tt.args.len))
				if (err != nil) != tt.wantErr {
					t.Errorf("WriteAt() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotN != tt.wantN {
					t.Errorf("WriteAt() gotN = %v, want %v", gotN, tt.wantN)
				}
			}
			str := fio.String()
			if fio.String() != tt.result {
				t.Errorf("buffer.string = %v, want %v", str, tt.result)
			}
			fmt.Printf("[%s] result: %s\n", tt.name, str)
		})
	}
}

func TestFakeIO_WriteRune(t *testing.T) {
	type args struct {
		data string
		size int64
	}
	tests := []struct {
		name    string
		args    args
		result  string
		wantN   int
		wantErr bool
	}{
		{name: "small buff", args: args{
			data: "あいうえお",
			size: 2,
		}, result: "あいうえお", wantN: 3, wantErr: false},
		{name: "default buff", args: args{
			data: "あいうえお",
			size: 0,
		}, result: "あいうえお", wantN: 3, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fio := NewBufferIO(tt.args.size)
			for _, r := range []rune(tt.args.data) {
				gotN, err := fio.WriteRune(r)
				if (err != nil) != tt.wantErr {
					t.Errorf("WriteRune() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotN != tt.wantN {
					t.Errorf("WriteRune() gotN = %v, want %v", gotN, tt.wantN)
				}
			}
			str := fio.String()
			if fio.String() != tt.result {
				t.Errorf("buffer.string = %v, want %v", str, tt.result)
			}
			fmt.Printf("[%s] result: %s\n", tt.name, str)
		})
	}
}

func TestFakeIO_Read(t *testing.T) {
	type args struct {
		data string
		size int64
		max  int
	}
	tests := []struct {
		name    string
		args    args
		result  string
		wantN   int
		wantErr bool
	}{
		{name: "small buff", args: args{
			data: "123456",
			size: 6,
			max:  3,
		}, result: "123", wantN: 1, wantErr: false},
		{name: "default buff", args: args{
			data: "123456",
			size: 0,
			max:  5,
		}, result: "12345", wantN: 1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fio := NewBufferIO(tt.args.size)
			buf := make([]byte, tt.args.max)
			datas := []byte(tt.args.data)
			for i := 0; i < tt.args.max; i++ {
				err := fio.WriteByte(datas[i])
				if (err != nil) != tt.wantErr {
					t.Errorf("WriteByte() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				gotN, err := fio.Read(buf[i:])
				if (err != nil) != tt.wantErr {
					t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotN != tt.wantN {
					t.Errorf("Read() gotN = %v, want %v", gotN, tt.wantN)
				}
			}
			str := string(buf)
			if str != tt.result {
				t.Errorf("buffer.string = %v, want %v", str, tt.result)
			}
			fmt.Printf("[%s] result: %s\n", tt.name, str)
		})
	}
}

func TestFakeIO_ReadAt(t *testing.T) {
	type args struct {
		data string
		size int64
		max  int
	}
	tests := []struct {
		name    string
		args    args
		result  string
		wantN   int
		wantErr bool
	}{
		{name: "small buff", args: args{
			data: "123456",
			size: 6,
			max:  3,
		}, result: "123", wantN: 1, wantErr: false},
		{name: "default buff", args: args{
			data: "123456",
			size: 0,
			max:  5,
		}, result: "12345", wantN: 1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fio := NewBufferIO(tt.args.size)
			buf := make([]byte, tt.args.max)
			_, err := fio.WriteString(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := 0; i < tt.args.max; i++ {
				gotN, err := fio.ReadAt(buf[i:i+1], int64(i))
				if (err != nil) != tt.wantErr {
					t.Errorf("ReadAt() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotN != tt.wantN {
					t.Errorf("ReadAt() gotN = %v, want %v", gotN, tt.wantN)
				}
			}
			str := string(buf)
			if str != tt.result {
				t.Errorf("buffer.string = %v, want %v", str, tt.result)
			}
			fmt.Printf("[%s] result: %s\n", tt.name, str)
		})
	}
}

func TestFakeIO_ReadRune(t *testing.T) {
	type args struct {
		data string
		size int64
		max  int
	}
	tests := []struct {
		name    string
		args    args
		result  string
		wantN   int
		wantErr bool
	}{
		{name: "small buff", args: args{
			data: "あいうえお",
			size: 6,
			max:  3,
		}, result: "あいう", wantN: 3, wantErr: false},
		{name: "default buff", args: args{
			data: "あいうえお",
			size: 0,
			max:  5,
		}, result: "あいうえお", wantN: 3, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fio := NewBufferIO(tt.args.size)
			buf := make([]rune, tt.args.max)
			_, err := fio.WriteString(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := 0; i < tt.args.max; i++ {
				r, gotN, err := fio.ReadRune()
				if (err != nil) != tt.wantErr {
					t.Errorf("ReadRune() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotN != tt.wantN {
					t.Errorf("ReadRune() gotN = %v, want %v", gotN, tt.wantN)
				}
				buf[i] = r
			}
			str := string(buf)
			if str != tt.result {
				t.Errorf("buffer.string = %v, want %v", str, tt.result)
			}
			fmt.Printf("[%s] result: %s\n", tt.name, str)
		})
	}
}

func TestFakeIO_ReadString(t *testing.T) {
	type args struct {
		data string
		size int64
		len  int
	}
	tests := []struct {
		name    string
		args    args
		result  string
		wantN   int
		wantErr bool
	}{
		{name: "small buff", args: args{
			data: "あいうえおLine1\nあいうえおLine2\nあいうえおLine3",
			size: 6,
			len:  3,
		}, result: "あああ", wantN: 3, wantErr: false},
		{name: "default buff", args: args{
			data: "あいうえおLine1\nあいうえおLine2\nあいうえおLine3\nあいうえおLine4\nあいうえおLine5\n\n\n",
			size: 0,
			len:  5,
		}, result: "あああああ", wantN: 3, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fio := NewBufferIO(tt.args.size)
			buf := make([]rune, tt.args.len)
			_, err := fio.WriteString(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			i := 0
			for {
				line, err := fio.ReadString('\n')
				if err == io.EOF {
					break
				} else if err != nil {
					t.Errorf("ReadString() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				runes := []rune(line)
				if line == "\n" {
					continue
				}
				gotN := utf8.RuneLen(runes[0])
				if gotN != tt.wantN {
					t.Errorf("Read() gotN = %v, want %v", gotN, tt.wantN)
					return
				}
				buf[i] = runes[0]
				i++
			}
			str := string(buf)
			if str != tt.result {
				t.Errorf("buffer.string = %v, want %v", str, tt.result)
			}
			fmt.Printf("[%s] result: %s\n", tt.name, str)
		})
	}
}
