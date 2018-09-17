package utils

import (
	"golang.org/x/text/encoding/charmap"

	"bytes"
	"fmt"
	"golang.org/x/text/transform"
	"io"
	"unicode/utf16"
	"unicode/utf8"
)

func Utf8ToWindows1251(b []byte) ( r []byte, err error) {
	buf := new(bytes.Buffer)
	wToWin1251 := transform.NewWriter(buf, charmap.Windows1251.NewEncoder())
	_, err = io.Copy(wToWin1251, bytes.NewReader(b))
	if err == nil {
		r = buf.Bytes()
	}
	return
}

func Windows1251ToUtf8(b []byte) ( []byte, error) {
	buf := new(bytes.Buffer)
	win1251toUtf8 := transform.NewWriter(buf, charmap.Windows1251.NewDecoder())
	if _, err := io.Copy(win1251toUtf8, bytes.NewReader(b)); err != nil {
		return nil, err
	}
	r := buf.Bytes()
	return r, nil

}

func EncodeUtf16(b []byte) (utf16bytes []byte) {

	var runes []rune
	for {
		if len(b) == 0 {
			break
		}
		rune, size := utf8.DecodeRune(b)
		runes = append(runes, rune)
		b = b[size:]
	}

	for _, x := range utf16.Encode(runes) {
		utf16bytes = append(utf16bytes, byte(x), byte(x>>8))
	}

	return
}

func DecodeUTF16(b []byte) (string, error) {

	if len(b)%2 != 0 {
		return "", fmt.Errorf("Must have even length byte slice")
	}

	u16s := make([]uint16, 1)

	ret := &bytes.Buffer{}

	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String(), nil
}
