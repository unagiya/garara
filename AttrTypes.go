package garara

type AttrDevice struct {
	Device DeviceType `xml:"device,attr"`
	Value  string     `xml:",innerxml"`
}

type AttrUse struct {
	Use UseType `xml:"use,attr"`
}

type AttrPart struct {
	Part PartType `xml:"part,attr"`
	CDATAString
}

type AttrType struct {
	Type string `xml:"type,attr"`
}

type AttrID struct {
	ID int `xml:"id,attr"`
}

type AttrIdString struct {
	ID    int    `xml:"id,attr"`
	Value string `xml:",innerxml"`
}

type AttrIdCdata struct {
	ID int `xml:"id,attr"`
	CDATAString
}

type AttrCodeString struct {
	Code  string `xml:"code,attr"`
	Value string `xml:",innerxml"`
}
