// Package conv

package conv

import (
	"io"

	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

type Decoder struct {
	e encoding.Encoding
}

// NewDecoder new encoder will use HTML escape sequences for runes that are not supported by the character set.
func NewDecoder(charSet string) (*Decoder, error) {
	e, _ := charset.Lookup(charSet)
	if e == nil {
		return nil, errors.Errorf("invalid charset [%s]", charSet)
	}
	return &Decoder{e: e}, nil
}

// GetEncoding get HTML character set encoder
func (d *Decoder) GetEncoding() encoding.Encoding {
	return d.e
}

// GetReader returns a new Reader that wraps r by transforming the bytes read via t. It calls Reset on t.
func (d *Decoder) GetReader(r io.Reader) *transform.Reader {
	return transform.NewReader(r, d.e.NewDecoder())
}

// ByteToString returns a new string with the result of converting b[:n] using t,
// where n <= len(b). If err == nil, n will be len(b). It calls Reset on t.
func (d *Decoder) ByteToString(src []byte) (string, error) {
	dst, _, err := transform.Bytes(d.e.NewDecoder(), src)
	if err != nil {
		return "", err
	}
	return BytesToString(dst), nil
}

// ByteToByte returns a new byte slice with the result of converting b[:n] using t,
// where n <= len(b). If err == nil, n will be len(b). It calls Reset on t.
func (d *Decoder) ByteToByte(src []byte) ([]byte, error) {
	dst, _, err := transform.Bytes(d.e.NewDecoder(), src)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// StringToByte returns a byte slice with the result of converting s[:n] using t, where
// n <= len(s). If err == nil, n will be len(s). It calls Reset on t.
func (d *Decoder) StringToByte(src string) ([]byte, error) {
	dst, _, err := transform.String(d.e.NewDecoder(), src)
	if err != nil {
		return nil, err
	}
	return StringToBytes(dst), nil
}

// StringToString returns a string with the result of converting s[:n] using t, where
// n <= len(s). If err == nil, n will be len(s). It calls Reset on t.
func (d *Decoder) StringToString(src string) (string, error) {
	dst, _, err := transform.String(d.e.NewDecoder(), src)
	if err != nil {
		return "", err
	}
	return dst, nil
}
