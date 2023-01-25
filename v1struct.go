package garara

import "encoding/xml"

type V1MailRequest struct {
	MailRoot
	Delivery []Delivery `xml:"delivery"`
}

type Delivery struct {
	AttrID
	Action    ActionType `xml:"action"`
	RequestID string     `xml:"request_id,omitempty"`
	Setting   Setting    `xml:"setting"`
	Contents  Contents   `xml:"contents"`
	SendList  SendList   `xml:"send_list"`
	//下記メール送信では不使用
	DeliverID string `xml:"deliver_id,omitempty"`
	Term      Term   `xml:"term,omitempty"`
	Mode      string `xml:"mode"`
}

type Term struct {
	From string `xml:"from,omitempty"`
	To   string `xml:"to,omitempty"`
}

type Setting struct {
	SendDate     string      `xml:"send_date"`
	FromName     CDATAString `xml:"from_name,omitempty"`
	From         string      `xml:"from,omitempty"`
	EnvelopeFrom string      `xml:"envelope_from,omitempty"`
	Option       Option      `xml:"option"`
	Throttle     int         `xml:"throttle,omitempty"`
	SMime        AttrUse     `xml:"s_mime,omitempty"`
	OpenedFlag   AttrUse     `xml:"opened_flag,omitempty"`
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
	AttachFile []string       `xml:"attach_file,omitempty"`
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

type V1MailResult struct {
	MailRoot
	Delivery []ResDelivery  `xml:"delivery"`
	Result   AttrCodeString `xml:"result"`
	Errors   V1Errors       `xml:"errors"`
}

type ResDelivery struct {
	AttrID
	DeliverID  string         `xml:"deliver_id"`
	RequestID  string         `xml:"request_id,omitempty"`
	From       string         `xml:"from,omitempty"`
	To         string         `xml:"to,omitempty"`
	ExecCount  string         `xml:"exec_cnt,omitempty"`
	Status     AttrCodeString `xml:"status,omitempty"`
	Action     string         `xml:"action,omitempty"`
	SentCount  int            `xml:"sent_cnt,omitempty"`
	ErrorCount int            `xml:"error_cnt,omitempty"`
	Result     AttrCodeString `xml:"result"`
	Errors     V1Errors       `xml:"errors"`
}

type AuthMail struct {
	MailRoot
	Auth Auth `xml:"auth"`
}

type Auth struct {
	Site    AttrID      `xml:"site"`
	Service AttrID      `xml:"service"`
	Name    CDATAString `xml:"name"`
	Pass    CDATAString `xml:"pass"`
}

type V1AuthRequest struct {
	AuthMail
	Delivery []Delivery `xml:"delivery"`
}

type V1Errors struct {
	Error []V1Error `xml:"error"`
	Err   error
}

type V1Error struct {
	AttrCodeString
	Kind string `xml:"kind,attr"`
}

type ClickLog struct {
	ClickDate   string
	URL         string
	UrlID       string
	CarrierCode string
	DeviceCode  string
	MailAddress string
	UserAgent   string
	AddInfo     string
	DeliverID   string
	OpenedFlag  string
}

type ErrorFilter struct {
	XMLName   xml.Name          `xml:"data"`
	Result    ErrorFilterResult `xml:"result"`
	ListCount int               `xml:"list_cnt"`
	List      []ErrorUser       `xml:"list>user"`
}

type ErrorFilterResult struct {
	Code    string `xml:"result_code"`
	Message string `xml:"result_message"`
}

type ErrorUser struct {
	MailAddress CDATAString `xml:"mail_addr"`
	BounceCount float32     `xml:"bounce_cnt"`
	RegistDate  string      `xml:"regist_date"`
	UpdateDate  string      `xml:"update_date"`
	ExcludeFlg  ExcludeType `xml:"exclude_flg"`
}

type ErrorList struct {
	Result   AttrCodeString  `xml:"result"`
	Errors   V1Errors        `xml:"errors"`
	Delivery []ErrorDelivery `xml:"delivery"`
}

type ErrorDelivery struct {
	AttrID
	DeliveryID string          `xml:"delivery_id"`
	From       string          `xml:"from"`
	To         string          `xml:"to"`
	ErrorList  []ErrorListData `xml:"error_list>data"`
	Result     AttrCodeString  `xml:"result"`
	Errors     V1Errors        `xml:"errors"`
}

type ErrorListData struct {
	AttrCharacterID
	SentDate    string        `xml:"sent_date"`
	ToAddress   CDATAString   `xml:"to_addr"`
	FromAddress CDATAString   `xml:"from_addr"`
	KeyField    CDATAString   `xml:"key_field"`
	Message     AttrCodeCdata `xml:"message"`
}

type ResultList struct {
	MailRoot
	Delivery []ResultDelivery `xml:"delivery"`
}

type ResultDelivery struct {
	AttrID
	DeliverID string         `xml:"deliver_id"`
	SentList  []SentData     `xml:"sent_list>data"`
	Result    AttrCodeString `xml:"result"`
}

type SentData struct {
	AttrCharacterID
	SentDate    string        `xml:"sent_date"`
	ToAddress   CDATAString   `xml:"to_addr"`
	FromAddress CDATAString   `xml:"from_addr"`
	KeyField    CDATAString   `xml:"key_field"`
	Status      string        `xml:"status"`
	Message     AttrCodeCdata `xml:"message"`
}
