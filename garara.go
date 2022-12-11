package garara

import "encoding/xml"

type requestBase struct {
	XMLName xml.Name `xml:"mail"`
}

type CDATAString struct {
	CDATA string `xml:",cdata"`
}

/*
func hoge() {
	v := AttrIdCdata{ID: "1", Value: CDATAString{CDATA: "ほげ"}}
}
*/
