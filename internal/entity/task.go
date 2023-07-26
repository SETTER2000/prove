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
	TaskList []Task

	TaskResponse struct {
		*config.Config `json:"-"`
		Number         string `json:"number,omitempty"`
		MetaData       string `json:"meta_data,omitempty"`
		UploadedAt     string `json:"uploaded_at" db:"uploaded_at"`
		UserID         string `json:"user_id,omitempty" db:"user_id"`
	}

	Task struct {
		TaskID      string `json:"task_id,omitempty" :"taskID"`
		TaskName    string `json:"name"  validate:"required" :"taskName"`
		Description string `json:"description"  validate:"required" :"description"`
		Price       string `json:"price" validate:"required"`
		UploadedAt  string `json:"created,omitempty" :"uploadedAt"`
	}
)

func (c *Task) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.TaskName, validation.Required),
		validation.Field(&c.Description, validation.Required),
		validation.Field(&c.Price, validation.Required),
	)
}

func (c *Task) Luna() error {
	number, err := strconv.Atoi(c.TaskName)
	if err != nil {
		return er.ErrValidTask
	}
	// проверка формата номера заказа
	if !luna.Luna(number) { // ...цветы, цветы 😁
		fmt.Println("неверный формат номера карты")
		return er.ErrValidTask
	}
	return nil
}
