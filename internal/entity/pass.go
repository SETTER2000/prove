package entity

import validation "github.com/go-ozzo/ozzo-validation"

type Pass struct {
	Login    string `json:"login"   validate:"required"`
	Password string `json:"password"   validate:"required"`
	MetaData string `json:"meta_data,omitempty"`
	UserID   string `json:"user_id,omitempty"`
}

func (c *Pass) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Login, validation.Required),
		validation.Field(&c.Password, validation.Required),
		validation.Field(&c.UserID, validation.Required),
		//validation.Field(&c.Number, validation.By(requiredIf(c.EncryptPassword == "")), validation.Length(1, 64)),
	)
	//return validation.ValidateStruct(a, validation.Field(&a.Email, validation.Required, is.Email))
}

// Sanitize очищает поля, для того чтоб они не возвращались в ответе
func (c *Pass) Sanitize() {
	c.UserID = ""
}
