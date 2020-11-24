package transform

import (
	"bytes"
	"fmt"

	"github.com/beevik/etree"
)

type XmlGenerator struct {
	doc  *etree.Document
	root *etree.Element
}

func (xml *XmlGenerator) AddTag(tagName string, value string) {
	tag := xml.root.CreateElement(tagName)
	tag.SetText(value)
}

func (xml *XmlGenerator) ToXml() (s string, err error) {
	buf := &bytes.Buffer{}
	xml.doc.Indent(2)
	_, err = xml.doc.WriteTo(buf)
	s = string(buf.Bytes())
	return
}

func NewXmlGenerator(rootTag string) *XmlGenerator {
	doc := etree.NewDocument()
	root := doc.CreateElement(rootTag)
	return &XmlGenerator{doc, root}
}

func Xml2mapWithRoot(body []byte, rootName string) (map[string]string, error) {
	doc := etree.NewDocument()
	err := doc.ReadFromBytes(body)
	if err != nil {
		return nil, err
	}
	root := doc.SelectElement(rootName)
	if root == nil {
		return nil, fmt.Errorf("no root named \"%s\" found", rootName)
	}
	tags := root.ChildElements()
	if tags == nil {
		return nil, fmt.Errorf("no xml children")
	}

	res := make(map[string]string, len(tags))
	for _, tag := range tags {
		res[tag.Tag] = tag.Text()
	}
	return res, nil
}
