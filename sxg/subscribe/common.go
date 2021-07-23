package subscribe

func AssertInt(i interface{}) int {
	switch v := i.(type) {
	case nil:
		return 0
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	default:
		return 0
	}
}
