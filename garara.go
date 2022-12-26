package garara

import "encoding/xml"

type MailRoot struct {
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
