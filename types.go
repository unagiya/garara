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

func (d DeviceType) String() string {
	switch d {
	case NONE:
		return "NONE"
	case PC:
		return "PC"
	case PORTABLE_PHONE:
		return "PORTABLE_PHONE"
	case I_PHONE:
		return "I_PHONE"
	case ANDROID_PHONE:
		return "ANDROID_PHONE"
	case WINDOWS_PHONE:
		return "WINDOWS_PHONE"
	default:
		return "Unknown"
	}
}

type TimingType string

const (
	REQUEST  TimingType = "request"
	DELIVERY TimingType = "delivery"
)

type ExcludeType string

const (
	INVALID = "0" //無効
	VALID   = "1" //有効
)
