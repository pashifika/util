// Package mem
/*
 * Version: 1.0.0
 * Copyright (c) 2022. Pashifika
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
package mem

// Simple byte buffer for marshaling data.

import (
	"bytes"
	"errors"
	"io"
	"unicode/utf8"

	"github.com/pashifika/util/conv"
)

// smallBufferSize is an initial allocation minimal capacity.
const smallBufferSize = 64

// A FakeIO is a variable-sized buffer of bytes with Read and Write methods.
// The zero value for FakeIO is an empty buffer ready to use.
type FakeIO struct {
	buf      []byte // contents are the bytes buf[off : len(buf)]
	off      int64  // read at &buf[off], write at &buf[len(buf)]
	lastRead readOp // last read operation, so that Unread* can work correctly.
}

// The readOp constants describe the last action performed on
// the buffer, so that UnreadRune and UnreadByte can check for
// invalid usage. opReadRuneX constants are chosen such that
// converted to int they correspond to the rune size that was read.
type readOp int8

// Don't use iota for these, as the values need to correspond with the
// names and comments, which is easier to see when being explicit.
const (
	opRead      readOp = -1 // Any other read operation.
	opInvalid   readOp = 0  // Non-read operation.
	opReadRune1 readOp = 1  // Read rune of size 1.
)

// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.
var ErrTooLarge = errors.New("bytes.FakeIO: too large")
var errNegativeRead = errors.New("bytes.FakeIO: reader returned negative count from Read")

const maxInt = int(^uint(0) >> 1)

// Bytes returns a slice of length b.Len() holding the unread portion of the buffer.
// The slice is valid for use only until the next buffer modification (that is,
// only until the next call to a method like Read, Write, Reset, or Truncate).
// The slice aliases the buffer content at least until the next buffer modification,
// so immediate changes to the slice will affect the result of future reads.
func (fio *FakeIO) Bytes() []byte { return fio.buf[fio.off:] }

// String returns the contents of the unread portion of the buffer
// as a string. If the FakeIO is a nil pointer, it returns "<nil>".
//
// To build strings more efficiently, see the strings.Builder type.
func (fio *FakeIO) String() string {
	if fio == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	return conv.BytesToString(fio.buf[fio.off:])
}

// empty reports whether the unread portion of the buffer is empty.
func (fio *FakeIO) empty() bool { return len(fio.buf) <= int(fio.off) }

// Len returns the number of bytes of the unread portion of the buffer;
// b.Len() == len(b.Bytes()).
func (fio *FakeIO) Len() int { return len(fio.buf) - int(fio.off) }

// Cap returns the capacity of the buffer's underlying byte slice, that is, the
// total space allocated for the buffer's data.
func (fio *FakeIO) Cap() int { return cap(fio.buf) }

func (fio *FakeIO) Size() int64 { return int64(len(fio.buf)) }

// Truncate discards all but the first n unread bytes from the buffer
// but continues to use the same allocated storage.
// It panics if n is negative or greater than the length of the buffer.
func (fio *FakeIO) Truncate(n int) {
	if n == 0 {
		fio.Reset()
		return
	}
	fio.lastRead = opInvalid
	if n < 0 || n > fio.Len() {
		panic("bytes.FakeIO: truncation out of range")
	}
	fio.buf = fio.buf[:int(fio.off)+n]
}

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as Truncate(0).
func (fio *FakeIO) Reset() {
	fio.buf = fio.buf[:0]
	fio.off = 0
	fio.lastRead = opInvalid
}

// ResetTo resets the Reader to be reading from b.
func (fio *FakeIO) ResetTo(b []byte) { *fio = FakeIO{buf: b, off: 0, lastRead: opRead} }

// tryGrowByReslice is a inlineable version of grow for the fast-case where the
// internal buffer only needs to be resliced.
// It returns the index where bytes should be written and whether it succeeded.
func (fio *FakeIO) tryGrowByReslice(n int) (int, bool) {
	if l := len(fio.buf); n <= cap(fio.buf)-l {
		fio.buf = fio.buf[:l+n]
		return l, true
	}
	return 0, false
}

// grow grows the buffer to guarantee space for n more bytes.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func (fio *FakeIO) grow(n int) int {
	m := fio.Len()
	// If buffer is empty, reset to recover space.
	if m == 0 && fio.off != 0 {
		fio.Reset()
	}
	// Try to grow by means of a reslice.
	if i, ok := fio.tryGrowByReslice(n); ok {
		return i
	}
	if fio.buf == nil && n <= smallBufferSize {
		fio.buf = make([]byte, n, smallBufferSize)
		return 0
	}
	c := cap(fio.buf)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(fio.buf, fio.buf[fio.off:])
	} else if c > maxInt-c-n {
		panic(ErrTooLarge)
	} else {
		// Not enough space anywhere, we need to allocate.
		buf := makeSlice(2*c + n)
		copy(buf, fio.buf[fio.off:])
		fio.buf = buf
	}
	// Restore b.off and len(b.buf).
	fio.off = 0
	fio.buf = fio.buf[:m+n]
	return m
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the
// buffer without another allocation.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with ErrTooLarge.
func (fio *FakeIO) Grow(n int) {
	if n < 0 {
		panic("bytes.FakeIO.Grow: negative count")
	}
	m := fio.grow(n)
	fio.buf = fio.buf[:m]
}

// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
func (fio *FakeIO) Write(p []byte) (n int, err error) {
	fio.lastRead = opInvalid
	m, ok := fio.tryGrowByReslice(len(p))
	if !ok {
		m = fio.grow(len(p))
	}
	return copy(fio.buf[m:], p), nil
}

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with ErrTooLarge.
func (fio *FakeIO) WriteString(s string) (n int, err error) {
	fio.lastRead = opInvalid
	m, ok := fio.tryGrowByReslice(len(s))
	if !ok {
		m = fio.grow(len(s))
	}
	return copy(fio.buf[m:], conv.StringToBytes(s)), nil
}

// MinRead is the minimum slice size passed to a Read call by
// FakeIO.ReadFrom. As long as the FakeIO has at least MinRead bytes beyond
// what is required to hold the contents of r, ReadFrom will not grow the
// underlying buffer.
const MinRead = 512

// ReadFrom reads data from r until EOF and appends it to the buffer, growing
// the buffer as needed. The return value n is the number of bytes read. Any
// error except io.EOF encountered during the read is also returned. If the
// buffer becomes too large, ReadFrom will panic with ErrTooLarge.
func (fio *FakeIO) ReadFrom(r io.Reader) (n int64, err error) {
	fio.lastRead = opInvalid
	for {
		i := fio.grow(MinRead)
		fio.buf = fio.buf[:i]
		m, e := r.Read(fio.buf[i:cap(fio.buf)])
		if m < 0 {
			panic(errNegativeRead)
		}

		fio.buf = fio.buf[:i+m]
		n += int64(m)
		if e == io.EOF {
			return n, nil // e is EOF, so return nil explicitly
		}
		if e != nil {
			return n, e
		}
	}
}

// makeSlice allocates a slice of size n. If the allocation fails, it panics
// with ErrTooLarge.
func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	return make([]byte, n)
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (fio *FakeIO) WriteTo(w io.Writer) (n int64, err error) {
	fio.lastRead = opInvalid
	if nBytes := fio.Len(); nBytes > 0 {
		m, e := w.Write(fio.buf[fio.off:])
		if m > nBytes {
			panic("bytes.FakeIO.WriteTo: invalid Write count")
		}
		fio.off += int64(m)
		n = int64(m)
		if e != nil {
			return n, e
		}
		// all bytes should have been written, by definition of
		// Write method in io.Writer
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	// FakeIO is now empty; reset.
	fio.Reset()
	return n, nil
}

// WriteByte appends the byte c to the buffer, growing the buffer as needed.
// The returned error is always nil, but is included to match bufio.Writer's
// WriteByte. If the buffer becomes too large, WriteByte will panic with
// ErrTooLarge.
func (fio *FakeIO) WriteByte(c byte) error {
	fio.lastRead = opInvalid
	m, ok := fio.tryGrowByReslice(1)
	if !ok {
		m = fio.grow(1)
	}
	fio.buf[m] = c
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to the
// buffer, returning its length and an error, which is always nil but is
// included to match bufio.Writer's WriteRune. The buffer is grown as needed;
// if it becomes too large, WriteRune will panic with ErrTooLarge.
func (fio *FakeIO) WriteRune(r rune) (n int, err error) {
	// Compare as uint32 to correctly handle negative runes.
	if uint32(r) < utf8.RuneSelf {
		return 1, fio.WriteByte(byte(r))
	}
	fio.lastRead = opInvalid
	m, ok := fio.tryGrowByReslice(utf8.UTFMax)
	if !ok {
		m = fio.grow(utf8.UTFMax)
	}
	n = utf8.EncodeRune(fio.buf[m:m+utf8.UTFMax], r)
	fio.buf = fio.buf[:m+n]
	return n, nil
}

// Read reads the next len(p) bytes from the buffer or until the buffer
// is drained. The return value n is the number of bytes read. If the
// buffer has no data to return, err is io.EOF (unless len(p) is zero);
// otherwise it is nil.
func (fio *FakeIO) Read(p []byte) (n int, err error) {
	if fio.off >= int64(len(fio.buf)) {
		return 0, io.EOF
	}
	fio.lastRead = opInvalid
	if fio.empty() {
		// FakeIO is empty, reset to recover space.
		fio.Reset()
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n = copy(p, fio.buf[fio.off:])
	fio.off += int64(n)
	if n > 0 {
		fio.lastRead = opRead
	}
	return n, nil
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (fio *FakeIO) Next(n int) []byte {
	fio.lastRead = opInvalid
	m := fio.Len()
	if n > m {
		n = m
	}
	data := fio.buf[fio.off : fio.off+int64(n)]
	fio.off += int64(n)
	if n > 0 {
		fio.lastRead = opRead
	}
	return data
}

// ReadByte reads and returns the next byte from the buffer.
// If no byte is available, it returns error io.EOF.
func (fio *FakeIO) ReadByte() (byte, error) {
	if fio.empty() {
		// FakeIO is empty, reset to recover space.
		fio.Reset()
		return 0, io.EOF
	}
	c := fio.buf[fio.off]
	fio.off++
	fio.lastRead = opRead
	return c, nil
}

// ReadRune reads and returns the next UTF-8-encoded
// Unicode code point from the buffer.
// If no bytes are available, the error returned is io.EOF.
// If the bytes are an erroneous UTF-8 encoding, it
// consumes one byte and returns U+FFFD, 1.
func (fio *FakeIO) ReadRune() (r rune, size int, err error) {
	if fio.empty() {
		// FakeIO is empty, reset to recover space.
		fio.Reset()
		return 0, 0, io.EOF
	}
	c := fio.buf[fio.off]
	if c < utf8.RuneSelf {
		fio.off++
		fio.lastRead = opReadRune1
		return rune(c), 1, nil
	}
	r, n := utf8.DecodeRune(fio.buf[fio.off:])
	fio.off += int64(n)
	fio.lastRead = readOp(n)
	return r, n, nil
}

// UnreadRune unreads the last rune returned by ReadRune.
// If the most recent read or write operation on the buffer was
// not a successful ReadRune, UnreadRune returns an error.  (In this regard
// it is stricter than UnreadByte, which will unread the last byte
// from any read operation.)
func (fio *FakeIO) UnreadRune() error {
	if fio.lastRead <= opInvalid {
		return errors.New("bytes.FakeIO: UnreadRune: previous operation was not a successful ReadRune")
	}
	if fio.off >= int64(fio.lastRead) {
		fio.off -= int64(fio.lastRead)
	}
	fio.lastRead = opInvalid
	return nil
}

var errUnreadByte = errors.New("bytes.FakeIO: UnreadByte: previous operation was not a successful read")

// UnreadByte unreads the last byte returned by the most recent successful
// read operation that read at least one byte. If a write has happened since
// the last read, if the last read returned an error, or if the read read zero
// bytes, UnreadByte returns an error.
func (fio *FakeIO) UnreadByte() error {
	if fio.lastRead == opInvalid {
		return errUnreadByte
	}
	fio.lastRead = opInvalid
	if fio.off > 0 {
		fio.off--
	}
	return nil
}

// ReadBytes reads until the first occurrence of delim in the input,
// returning a slice containing the data up to and including the delimiter.
// If ReadBytes encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadBytes returns err != nil if and only if the returned data does not end in
// delim.
func (fio *FakeIO) ReadBytes(delim byte) (line []byte, err error) {
	slice, err := fio.readSlice(delim)
	// return a copy of slice. The buffer's backing array may
	// be overwritten by later calls.
	line = append(line, slice...)
	return line, err
}

// readSlice is like ReadBytes but returns a reference to internal buffer data.
func (fio *FakeIO) readSlice(delim byte) (line []byte, err error) {
	i := bytes.IndexByte(fio.buf[fio.off:], delim)
	end := fio.off + int64(i) + 1
	if i < 0 {
		end = int64(len(fio.buf))
		err = io.EOF
	}
	line = fio.buf[fio.off:end]
	fio.off = end
	fio.lastRead = opRead
	return line, err
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

// ReadAt implements the io.ReaderAt interface.
func (fio *FakeIO) ReadAt(b []byte, off int64) (n int, err error) {
	// cannot modify state - see io.ReaderAt
	if off < 0 {
		return 0, errors.New("FakeIO.Multi.ReadAt: negative offset")
	}
	if off >= int64(len(fio.buf)) {
		return 0, io.EOF
	}
	n = copy(b, fio.buf[off:])
	if n < len(b) {
		err = io.EOF
	}
	return
}

// Seek implements the io.Seeker interface.
func (fio *FakeIO) Seek(offset int64, whence int) (int64, error) {
	fio.lastRead = opRead
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = fio.off + offset
	case io.SeekEnd:
		abs = int64(len(fio.buf)) + offset
	default:
		return 0, errors.New("bytes.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("bytes.Reader.Seek: negative position")
	}
	fio.off = abs
	return abs, nil
}

// WriteAt writes a slice of bytes to a buffer starting at the position provided
// The number of bytes written will be returned, or error. Can overwrite previous
// written slices if the write ats overlap.
func (fio *FakeIO) WriteAt(p []byte, pos int64) (n int, err error) {
	pLen := len(p)
	expLen := pos + int64(pLen)
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

// NewFakeIO creates and initializes a new FakeIO using buf as its
// initial contents. The new FakeIO takes ownership of buf, and the
// caller should not use buf after this call. NewFakeIO is intended to
// prepare a FakeIO to read existing data. It can also be used to set
// the initial size of the internal buffer for writing. To do that,
// buf should have the desired capacity but a length of zero.
//
// In most cases, new(FakeIO) (or just declaring a FakeIO variable) is
// sufficient to initialize a FakeIO.
func NewFakeIO(buf []byte) *FakeIO { return &FakeIO{buf: buf} }

// NewFakeIOString creates and initializes a new FakeIO using string s as its
// initial contents. It is intended to prepare a buffer to read an existing
// string.
//
// In most cases, new(FakeIO) (or just declaring a FakeIO variable) is
// sufficient to initialize a FakeIO.
func NewFakeIOString(s string) *FakeIO {
	return &FakeIO{buf: []byte(s)}
}
