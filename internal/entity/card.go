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
	CardList []CardResponse

	CardResponse struct {
		*config.Config `json:"-"`
		Number         string `json:"number,omitempty"`
		MetaData       string `json:"meta_data,omitempty"`
		UploadedAt     string `json:"uploaded_at" db:"uploaded_at"`
		UserID         string `json:"user_id,omitempty" db:"user_id"`
	}
	Card struct {
		Number   string `json:"number"  validate:"required"`
		MetaData string `json:"meta_data,omitempty"`
		UserID   string `json:"user_id,omitempty"`
	}
)

func (c *Card) Validate() error {
	return validation.ValidateStruct(
		c, validation.Field(&c.Number, validation.Required),
		//validation.Field(&c.Number, validation.By(requiredIf(c.EncryptPassword == "")), validation.Length(1, 64)),
	)
	//return validation.ValidateStruct(a, validation.Field(&a.Email, validation.Required, is.Email))
}

func (c *Card) Luna() error {
	number, err := strconv.Atoi(c.Number)
	if err != nil {
		return er.ErrValidCard
	}
	// проверка формата номера заказа
	if !luna.Luna(number) { // ...цветы, цветы 😁
		fmt.Println("неверный формат номера карты")
		return er.ErrValidCard
	}
	return nil
}

// Sanitize очищает поля, для того чтоб они не возвращались в ответе
func (c *Card) Sanitize() {
	c.UserID = ""
}
