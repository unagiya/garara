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
	AttrID
	Value string `xml:",innerxml"`
}

type AttrIdCdata struct {
	AttrID
	CDATAString
}

type AttrCode struct {
	Code string `xml:"code,attr"`
}

type AttrCodeString struct {
	AttrCode
	Value string `xml:",innerxml"`
}

type AttrCodeCdata struct {
	AttrCode
	CDATAString
}

type AttrCharacterID struct {
	ID string `xml:"id,attr"`
}
