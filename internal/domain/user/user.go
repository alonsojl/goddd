package user

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	nameMinLength = 1
	nameMaxLength = 30
	phoneLength   = 10
	ageMin        = 1
	ageMax        = 100
	adult         = 18
)

func NameRules() []validation.Rule {
	return []validation.Rule{
		validation.Required,
		validation.Length(nameMinLength, nameMaxLength),
	}
}

func FirstNameRules() []validation.Rule {
	return []validation.Rule{
		validation.Required,
		validation.Length(nameMinLength, nameMaxLength),
	}
}

func LastNameRules() []validation.Rule {
	return []validation.Rule{
		validation.Required,
		validation.Length(nameMinLength, nameMaxLength),
	}
}

func EmailRules() []validation.Rule {
	return []validation.Rule{
		validation.Required,
		is.Email,
	}
}

func PhoneRules() []validation.Rule {
	return []validation.Rule{
		validation.Required,
		is.Int,
		validation.Length(phoneLength, phoneLength),
	}
}

func AgeRules() []validation.Rule {
	return []validation.Rule{
		validation.Required,
		validation.Min(ageMin),
		validation.Max(ageMax),
	}
}

func PositionRules() []validation.Rule {
	positions := []interface{}{"Go developer", "Rust developer", "Python developer"}
	return []validation.Rule{
		validation.In(positions...),
	}
}

type Params struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Age       int16  `json:"age"`
	Position  string `json:"position"`
}

func (p *Params) Validate() error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Name, NameRules()...),
		validation.Field(&p.FirstName, FirstNameRules()...),
		validation.Field(&p.LastName, LastNameRules()...),
		validation.Field(&p.Email, EmailRules()...),
		validation.Field(&p.Age, AgeRules()...),
		validation.Field(&p.Phone, validation.When(p.Age > adult, PhoneRules()...)),
		validation.Field(&p.Position, PositionRules()...),
	)
	if err != nil {
		return err
	}

	return nil
}
