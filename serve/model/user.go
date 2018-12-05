package model

import "activityjs.io/serve/errors"

// User 用户
type User struct {
	ID       *Identity
	Username string
	Avatar   string
	Level    int
}

func (usr *User) Set(key string, val interface{}) error {

	switch key {
	case "ID":
		usr.ID = val.(*Identity)
	case "Avatar":
		usr.Avatar = val.(string)
	case "Username":
		usr.Username = val.(string)
	case "Level":
		usr.Level = val.(int)
	default:
		return errors.New("无效的 User 字段 " + key)
	}
	return nil
}

func (usr *User) Get(key string) (interface{}, bool) {

	switch key {
	case "ID":
		return usr.ID, true
	case "Avatar":
		return usr.Avatar, true
	case "Username":
		return usr.Username, true
	case "Level":
		return usr.Level, true
	default:
		return nil, false
	}
}

func (usr *User) compute() {

}
