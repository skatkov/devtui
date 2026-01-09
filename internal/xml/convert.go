package xml

import (
	"bytes"
	"encoding/json"

	"github.com/clbanning/mxj/v2"
)

// XMLToJSON converts XML content to JSON format.
func XMLToJSON(xmlContent string) (string, error) {
	mv, err := mxj.NewMapXml([]byte(xmlContent))
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.MarshalIndent(mv, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// JSONToXML converts JSON content to XML format.
func JSONToXML(jsonContent string) (string, error) {
	mv, err := mxj.NewMapJson([]byte(jsonContent))
	if err != nil {
		return "", err
	}

	xmlBytes, err := mv.XmlIndent("", "  ")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.Write(xmlBytes)
	return buf.String(), nil
}
