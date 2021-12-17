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
	"bytes"
	"errors"
	"io"
	"sync"
	"unicode/utf8"

	"github.com/pashifika/util/conv"
)

// smallBufferSize is an initial allocation minimal capacity.
const smallBufferSize = 64

// FakeIO io.Write + io.Reader
type FakeIO struct {
	buf      []byte     // buff data
	m        sync.Mutex // multi control
	off      int64      // current writing index
	i        int64      // current reading index
	prevRune int        // index of previous rune; or < 0
}

// NewBufferIO creates a FakeIO with an internal buffer
// provided by buf.
func NewBufferIO(size int64) *FakeIO {
	bufSize := size
	if bufSize <= 0 {
		bufSize = smallBufferSize
	}
	return &FakeIO{buf: make([]byte, 0, bufSize), off: 0, i: 0}
}

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as Truncate(0).
func (fio *FakeIO) Reset() {
	fio.buf = fio.buf[:0]
	fio.off = 0
	fio.i = 0
	fio.prevRune = -1
}

// Truncate discards all but the first n unread bytes from the buffer
// but continues to use the same allocated storage.
// It panics if n is negative or greater than the length of the buffer.
func (fio *FakeIO) Truncate(n int64) {
	if n == 0 {
		fio.Reset()
		return
	}
	fio.m.Lock()
	defer fio.m.Unlock()
	if n < 0 || n > int64(len(fio.buf)) {
		panic("bytes.Buffer: truncation out of range")
	}
	fio.buf = fio.buf[:fio.off+n]
}

// Size returns the original length of the underlying byte slice.
// Size is the number of bytes available for reading via ReadAt.
// The returned value is always the same and is not affected by calls
// to any other method.
func (fio *FakeIO) Size() int64 {
	fio.m.Lock()
	defer fio.m.Unlock()
	return int64(len(fio.buf))
}

// Len returns the number of bytes of the unread portion of the buffer;
// b.Len() == len(b.Bytes()).
func (fio *FakeIO) Len() int {
	fio.m.Lock()
	defer fio.m.Unlock()
	return len(fio.buf)
}

// Cap returns the capacity of the buffer's underlying byte slice, that is, the
// total space allocated for the buffer's data.
func (fio *FakeIO) Cap() int {
	fio.m.Lock()
	defer fio.m.Unlock()
	return cap(fio.buf)
}

// Read implements the io.Reader interface.
func (fio *FakeIO) Read(p []byte) (n int, err error) {
	fio.m.Lock()
	defer fio.m.Unlock()
	if fio.i >= int64(len(fio.buf)) {
		return 0, io.EOF
	}
	fio.prevRune = -1
	n = copy(p, fio.buf[fio.i:])
	fio.i += int64(n)
	return
}

// ReadAt implements the io.ReaderAt interface.
func (fio *FakeIO) ReadAt(b []byte, off int64) (n int, err error) {
	// cannot modify state - see io.ReaderAt
	if off < 0 {
		return 0, errors.New("BufferIO.Reader.ReadAt: negative offset")
	}
	fio.m.Lock()
	defer fio.m.Unlock()
	if off >= int64(len(fio.buf)) {
		return 0, io.EOF
	}
	n = copy(b, fio.buf[off:])
	if n < len(b) {
		err = io.EOF
	}
	return
}

// ReadByte implements the io.ByteReader interface.
func (fio *FakeIO) ReadByte() (byte, error) {
	fio.m.Lock()
	defer fio.m.Unlock()
	fio.prevRune = -1
	if fio.i >= int64(len(fio.buf)) {
		return 0, io.EOF
	}
	b := fio.buf[fio.i]
	fio.i++
	return b, nil
}

// UnreadByte complements ReadByte in implementing the io.ByteScanner interface.
func (fio *FakeIO) UnreadByte() error {
	fio.m.Lock()
	defer fio.m.Unlock()
	if fio.i <= 0 {
		return errors.New("BufferIO.Reader.UnreadByte: at beginning of slice")
	}
	fio.prevRune = -1
	fio.i--
	return nil
}

// ReadRune implements the io.RuneReader interface.
func (fio *FakeIO) ReadRune() (ch rune, size int, err error) {
	fio.m.Lock()
	defer fio.m.Unlock()
	if fio.i >= int64(len(fio.buf)) {
		fio.prevRune = -1
		return 0, 0, io.EOF
	}
	fio.prevRune = int(fio.i)
	if c := fio.buf[fio.i]; c < utf8.RuneSelf {
		fio.i++
		return rune(c), 1, nil
	}
	ch, size = utf8.DecodeRune(fio.buf[fio.i:])
	fio.i += int64(size)
	return
}

// UnreadRune complements ReadRune in implementing the io.RuneScanner interface.
func (fio *FakeIO) UnreadRune() error {
	fio.m.Lock()
	defer fio.m.Unlock()
	if fio.i <= 0 {
		return errors.New("BufferIO.Reader.UnreadRune: at beginning of slice")
	}
	if fio.prevRune < 0 {
		return errors.New("BufferIO.Reader.UnreadRune: previous operation was not ReadRune")
	}
	fio.i = int64(fio.prevRune)
	fio.prevRune = -1
	return nil
}

// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
func (fio *FakeIO) Write(p []byte) (n int, err error) {
	pLen := len(p)
	fio.m.Lock()
	defer fio.m.Unlock()
	bufLen := len(fio.buf)
	expLen := bufLen + pLen
	if bufLen < expLen {
		if cap(fio.buf) < expLen {
			nextCap := smallBufferSize
			if nextCap < pLen {
				nextCap = pLen
			}
			newBuf := make([]byte, bufLen, bufLen+nextCap)
			copy(newBuf, fio.buf)
			fio.buf = newBuf
		}
		fio.buf = fio.buf[:fio.off+int64(pLen)]
	}
	copy(fio.buf[fio.off:], p)
	fio.off += int64(pLen)
	return pLen, nil
}

// WriteAt writes a slice of bytes to a buffer starting at the position provided
// The number of bytes written will be returned, or error. Can overwrite previous
// written slices if the write ats overlap.
func (fio *FakeIO) WriteAt(p []byte, pos int64) (n int, err error) {
	pLen := len(p)
	expLen := pos + int64(pLen)
	fio.m.Lock()
	defer fio.m.Unlock()
	if int64(len(fio.buf)) < expLen {
		if int64(cap(fio.buf)) < expLen {
			newBuf := make([]byte, expLen, expLen)
			copy(newBuf, fio.buf)
			fio.buf = newBuf
		}
		fio.buf = fio.buf[:expLen]
	}
	copy(fio.buf[pos:], p)
	return pLen, nil
}

// WriteByte appends the byte c to the buffer, growing the buffer as needed.
// The returned error is always nil, but is included to match bufio.Writer's
// WriteByte. If the buffer becomes too large, WriteByte will panic with
// ErrTooLarge.
func (fio *FakeIO) WriteByte(c byte) error {
	_, err := fio.Write([]byte{c})
	return err
}

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with ErrTooLarge.
func (fio *FakeIO) WriteString(s string) (n int, err error) {
	return fio.Write(conv.StringToBytes(s))
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to the
// buffer, returning its length and an error, which is always nil but is
// included to match bufio.Writer's WriteRune. The buffer is grown as needed;
// if it becomes too large, WriteRune will panic with ErrTooLarge.
func (fio *FakeIO) WriteRune(r rune) (n int, err error) {
	var buf []byte
	// Compare as uint32 to correctly handle negative runes.
	if uint32(r) < utf8.RuneSelf {
		buf = []byte{byte(r)}
	} else {
		buf = make([]byte, utf8.RuneLen(r))
		_ = utf8.EncodeRune(buf, r)
	}
	return fio.Write(buf)
}

// readSlice is like ReadBytes but returns a reference to internal buffer data.
func (fio *FakeIO) readSlice(delim byte) (line []byte, err error) {
	fio.m.Lock()
	defer fio.m.Unlock()

	i := int64(bytes.IndexByte(fio.buf[fio.i:], delim))
	end := fio.i + i + 1
	if i < 0 {
		end = int64(len(fio.buf))
		if len(fio.buf[fio.i:]) == 0 || fio.buf[fio.i:][0] == delim {
			err = io.EOF
		}
	}
	line = fio.buf[fio.i:end]
	fio.i = end
	return line, err
}

// ReadBytes reads until the first occurrence of delim in the input,
// returning a slice containing the data up to and including the delimiter.
// If ReadBytes encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadBytes returns err != nil if and only if the returned data does not end in
// delim.
func (fio *FakeIO) ReadBytes(delim byte) (line []byte, err error) {
	return fio.readSlice(delim)
}

// ReadString reads until the first occurrence of delim in the input,
// returning a string containing the data up to and including the delimiter.
// If ReadString encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadString returns err != nil if and only if the returned data does not end
// in delim.
func (fio *FakeIO) ReadString(delim byte) (line string, err error) {
	slice, err := fio.readSlice(delim)
	return conv.BytesToString(slice), err
}

// Seek implements the io.Seeker interface.
func (fio *FakeIO) Seek(offset int64, whence int) (int64, error) {
	fio.m.Lock()
	defer fio.m.Unlock()
	fio.prevRune = -1
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = fio.i + offset
	case io.SeekEnd:
		abs = int64(len(fio.buf)) + offset
	default:
		return 0, errors.New("BufferIO.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("BufferIO.Reader.Seek: negative position")
	}
	fio.i = abs
	return abs, nil
}

// Bytes returns a slice of bytes written to the buffer.
func (fio *FakeIO) Bytes() []byte {
	fio.m.Lock()
	defer fio.m.Unlock()
	return fio.buf
}

// String returns the contents of the unread portion of the buffer
// as a string. If the Buffer is a nil pointer, it returns "<nil>".
//
// To build strings more efficiently, see the strings.Builder type.
func (fio *FakeIO) String() string {
	if fio == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	fio.m.Lock()
	defer fio.m.Unlock()
	return conv.BytesToString(fio.buf)
}