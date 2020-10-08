package udpstack

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

type CellInfo struct{ LAC, CellID string }

func ParseCellInfo(input string) *CellInfo {
	re := regexp.MustCompile(`LAC:(\w+),CELL ID:(\w+)`)
	if matched := re.FindStringSubmatch(input); len(matched) == 3 {
		return &CellInfo{
			LAC:    matched[1],
			CellID: matched[2],
		}
	}
	return nil
}

func onReplyNotify(message interface{}, err error) string {
	var (
		v    = reflect.ValueOf(message).Elem()
		t    = v.Type()
		name = t.Field(0).Tag.Get(tagField)
		id   = v.Field(0).String()
	)
	if err != nil {
		log.Error().
			Str("action", "onReplyNotify").
			Str("name", name).
			Str("id", "id").
			Err(err)
	}
	if name == "AT" || name == "GSM" {
		return ""
	} else if name == "req" {
		status := 0
		if err != nil {
			status = -1
		}
		return fmt.Sprintf("reg:%s;status:%d;", id, status)
	}
	values := []string{name, id, "OK"}
	if err == ErrDisabledModule {
		values[2] = "DISABLE"
	} else if err != nil {
		values[2] = "ERROR"
		values = append(values, err.Error())
	}
	return strings.Join(values, cmdJoinSymbol)
}
