package entity

import (
	"fmt"
	"github.com/SETTER2000/prove/config"
	"github.com/SETTER2000/prove/pkg/er"
	"github.com/SETTER2000/prove/scripts/luna"
	validation "github.com/go-ozzo/ozzo-validation"
	"strconv"
)

type (
	GroupList []Group

	GroupResponse struct {
		*config.Config `json:"-"`
		Number         string `json:"number,omitempty"`
		MetaData       string `json:"meta_data,omitempty"`
		UploadedAt     string `json:"uploaded_at" db:"uploaded_at"`
		UserID         string `json:"user_id,omitempty" db:"user_id"`
	}

	Group struct {
		GroupID    string `json:"group_id,omitempty"`
		GroupName  string `json:"name"  validate:"required"`
		UploadedAt string `json:"created,omitempty"`
	}
)

func (c *Group) Validate() error {
	return validation.ValidateStruct(c, validation.Field(&c.GroupName, validation.Required))
}

func (c *Group) Luna() error {
	number, err := strconv.Atoi(c.GroupName)
	if err != nil {
		return er.ErrValidGroup
	}
	// проверка формата номера заказа
	if !luna.Luna(number) { // ...цветы, цветы 😁
		fmt.Println("неверный формат номера карты")
		return er.ErrValidGroup
	}
	return nil
}
