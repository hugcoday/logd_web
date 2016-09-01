package validate

import "gopkg.in/go-playground/validator.v8"

var Vd *validator.Validate

func Init() {
	// config validate
	config := &validator.Config{TagName: "validate"}

	Vd = validator.New(config)
}
