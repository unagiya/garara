package garara

type MailRequestHeader struct {
}
type V1MailRequest struct {
	requestBase
	Delivery []Delivery `xml:"delivery"`
}

type Delivery struct {
	AttrID
	Action    ActionType `xml:"action"`
	RequestID string     `xml:"request_id,omitempty"`
	Setting   Setting    `xml:"setting"`
	Contents  Contents   `xml:"contents"`
	SendList  SendList   `xml:"send_list"`
	//下記はresultで利用
	DeliveryID string `xml:"delivery_id,omitempty"`
}

type Setting struct {
	SendDate     string      `xml:"send_date"`
	FromName     CDATAString `xml:"from_name,omitempty"`
	From         string      `xml:"from,omitempty"`
	EnvelopeFrom string      `xml:"envelope_from,omitempty"`
	Option       Option      `xml:"option"`
	Throttle     int         `xml:"throttle,omitempty"`
	SMime        AttrUse     `xml:"s_mime,omitempty"`
	OpendFlag    AttrUse     `xml:"opend_flag,omitempty"`
}

type Option struct {
	StopTime      string `xml:"stop_time,omitempty"`
	StartTime     string `xml:"start_time,omitempty"`
	LifeTime      string `xml:"lifetime,omitempty"`
	RetryInterval int    `xml:"retry_interval,omitempty"`
}

type Contents struct {
	Subject    *CDATAString   `xml:"subject,omitempty"`
	Body       *AttrPart      `xml:"body,omitempty"`
	Image      []AttrIdString `xml:"image,omitempty"`
	Text       []AttrIdCdata  `xml:"text,omitempty"`
	AttachFile string         `xml:"attach_file,omitempty"`
	Template   *AttrID        `xml:"template,omitempty"`
}

type SendList struct {
	Data    []Data   `xml:"data,omitempty"`
	ExtData *ExtData `xml:"ext_data,omitempty"`
}

type Data struct {
	AttrID
	Address  AttrDevice     `xml:"address"`
	IntText  []AttrIdCdata  `xml:"int_txt,omitempty"`
	ExtText  []AttrIdString `xml:"ext_txt,omitempty"`
	ExtImage []AttrIdString `xml:"ext_img,omitempty"`
	KeyField string         `xml:"key_field,omitempty"`
}

type ExtData struct {
	AttrType
	ListID        string      `xml:"list_id,omitempty"`
	Query         CDATAString `xml:"query,omitempty"`
	ExtractTiming TimingType  `xml:"extract_timing,omitempty"`
}

type V1MailResurt struct {
	requestBase
	Delivery
}