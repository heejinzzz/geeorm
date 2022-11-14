package session

import (
	"github.com/heejinzzz/geeorm/log"
	"reflect"
)

type HookType string

// Hook Type constants
const (
	BeforeQuery  HookType = "BeforeQuery"
	AfterQuery   HookType = "AfterQuery"
	BeforeUpdate HookType = "BeforeUpdate"
	AfterUpdate  HookType = "AfterUpdate"
	BeforeDelete HookType = "BeforeDelete"
	AfterDelete  HookType = "AfterDelete"
	BeforeInsert HookType = "BeforeInsert"
	AfterInsert  HookType = "AfterInsert"
)

// CallMethod calls the registered hooks
func (s *Session) CallMethod(method HookType, value interface{}) {
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(string(method))
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(string(method))
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		log.Infof("Find %s Hook", method)
		if v := fm.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}
