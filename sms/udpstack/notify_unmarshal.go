package udpstack

import (
	"reflect"
	"strconv"
	"strings"
)

const (
	tagField         = "field"
	cmdJoinSymbol    = " "
	notifyJoinSymbol = ";"
	notifyPairSymbol = ":"
)

//noinspection SpellCheckingInspection
func unmarshalNotifyPacket(packet string, r interface{}) error {
	mapping := make(map[string]string)
	packet = strings.TrimSpace(packet)
	packet = strings.TrimSuffix(packet, notifyJoinSymbol)
	for _, block := range strings.Split(packet, notifyJoinSymbol) {
		pair := strings.SplitN(block, notifyPairSymbol, 2)
		if len(pair) != 2 {
			return errParseNotifyError
		}
		key := strings.TrimSpace(pair[0])
		value := strings.TrimSpace(pair[1])
		if value == "" || value == "(null)" || value == "<NULL>" {
			continue
		}
		mapping[key] = value
	}
	v := reflect.ValueOf(r).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		value, ok := mapping[t.Field(i).Tag.Get(tagField)]
		if !ok {
			continue
		}
		field := v.Field(i)
		switch field.Interface().(type) {
		case bool:
			field.SetBool(value == "Y" || value == "1")
		case string:
			field.SetString(value)
		case *string:
			field.Set(reflect.ValueOf(&value))
		case VoIPStatus:
			field.Set(reflect.ValueOf(VoIPStatus(value)))
		case []string:
			separator := ","
			value = strings.Trim(value, separator)
			splitted := strings.Split(value, separator)
			field.Set(reflect.ValueOf(splitted))
		case *int:
			n, _ := strconv.ParseInt(value, 10, 64)
			if n > 0 {
				field.SetInt(n)
			}
		case int, int8, int16, int32, int64:
			n, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return errParseNotifyError
			}
			field.SetInt(n)
		case uint, uint8, uint16, uint32, uint64:
			n, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return errParseNotifyError
			}
			field.SetUint(n)
		case Direction:
			n, _ := strconv.ParseInt(value, 10, 64)
			field.Set(reflect.ValueOf(Direction(n)))
		}
	}
	return nil
}
