package httpmodels

import (
	"github.com/esmailemami/eshop/app/validations"
	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/models"
	"github.com/esmailemami/eshop/services/numeric"
	"github.com/esmailemami/eshop/services/sanitize"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type ExampleReqModel struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (model *ExampleReqModel) ApplyDataTransformation() {
	model.Name = numeric.TransformFa2En(model.Name)
	model.Code = numeric.TransformFa2En(model.Code)
}

func (model ExampleReqModel) ValidateCreate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(
			&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
		),
	)
}

func (model ExampleReqModel) ValidateUpdate(id uuid.UUID) error {
	return validation.ValidateStruct(
		&model,
		validation.Field(
			&model.Name,
			validation.Required.Error(consts.Required),
			validation.By(validations.ClearText()),
		),
		validation.Field(&model.Code,
			validation.Required.Error(consts.Required),
			validation.By(validations.Code()),
		),
	)
}

func (model ExampleReqModel) ToDBModel(rootModel models.Example) models.Example {
	rootModel.Name = sanitize.AsClearText(model.Name)
	rootModel.Code = sanitize.AsCode(model.Code)

	return rootModel
}
