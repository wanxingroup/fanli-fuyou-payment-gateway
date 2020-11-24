package transform

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/gin-gonic/gin/binding"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var XmlGbk = xmlGBKBinding{}

type xmlGBKBinding struct {
}

func (xmlGBKBinding) Bind(b []byte, obj interface{}) error {
	r, _, err := transform.Bytes(simplifiedchinese.GBK.NewDecoder(), b)
	if err != nil {
		return err
	}
	decoder := xml.NewDecoder(bytes.NewReader(r))
	decoder.CharsetReader = func(c string, i io.Reader) (reader io.Reader, err error) {
		return charset.NewReaderLabel(c, i)
	}
	if err = decoder.Decode(obj); err != nil {
		return err
	}
	return binding.Validator.ValidateStruct(obj)
}
