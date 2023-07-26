package entity

import validation "github.com/go-ozzo/ozzo-validation"

type Text struct {
	Text     string `json:"text"   validate:"required"`
	MetaData string `json:"meta_data,omitempty"`
	UserID   string `json:"user_id,omitempty"`
}

func (c *Text) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Text, validation.Required),
		validation.Field(&c.UserID, validation.Required),
		//validation.Field(&c.Number, validation.By(requiredIf(c.EncryptPassword == "")), validation.Length(1, 64)),
	)
	//return validation.ValidateStruct(a, validation.Field(&a.Email, validation.Required, is.Email))
}

// Sanitize очищает поля, для того чтоб они не возвращались в ответе
func (c *Text) Sanitize() {
	c.UserID = ""
}
