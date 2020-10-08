package business

import (
	"encoding/json"
	"reflect"

	"github.com/jinzhu/gorm"
)

type CronTask struct {
	gorm.Model
	CardID     uint
	Card       Card
	Name       string
	Expression string // cron expression
	Kind       string
	Options    string `gorm:"type:json1"`
}

func (t *CronTask) SetTask(task interface{}) error {
	encoded, err := json.Marshal(task)
	if err != nil {
		return err
	}
	t.Kind = reflect.TypeOf(task).Elem().Name()
	t.Options = string(encoded)
	return nil
}
