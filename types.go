package garara

type ActionType string

const (
	RESERVE    ActionType = "reserve"
	GET_ERROR  ActionType = "get_error"
	GET_RESULT ActionType = "get_result"
	GET_STATUS ActionType = "get_status"
	DELETE     ActionType = "delete"
)

type UseType int

const (
	UNUSE UseType = iota
	USE
)

type PartType string

const (
	TEXT PartType = "text"
	HTML PartType = "html"
)

type DeviceType int

const (
	NONE DeviceType = iota
	PC
	PORTABLE_PHONE
	I_PHONE
	ANDROID_PHONE
	WINDOWS_PHONE
)

type TimingType string

const (
	REQUEST  TimingType = "request"
	DELIVERY TimingType = "delivery"
)
