package controllers

import (
	"net/http"

	"github.com/esmailemami/eshop/app"
	"github.com/esmailemami/eshop/app/consts"
	"github.com/esmailemami/eshop/app/errors"
	appmodels "github.com/esmailemami/eshop/app/models"
	"github.com/esmailemami/eshop/app/services/settings"
	"github.com/esmailemami/eshop/app/services/settings/processor"
)

// GetSetting godoc
// @Summary Settings
// @Description Settings
// @Tags Settings
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param item path string true "item to get list" Enums(systemSetting)
// @Success 200 {object} []processor.SettingItem
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/settings/{item} [get]
func GetSettings(ctx *app.HttpContext) error {
	var items []processor.SettingItem

	key := ctx.GetPathParam("item")

	switch key {
	case "systemSetting":
		items = settings.GetSystemSettingItems()
	default:
		return errors.NewBadRequestError("Invalid item", nil)
	}

	return ctx.JSON(items, http.StatusOK)
}

// Update Setting godoc
// @Tags Settings
// @Accept json
// @Produce json
// @Security Bearer
// @Param item path string true "item to get list" Enums(systemSetting)
// @Param Setting   body  appmodels.SettingsReqModel  true  "Setting model"
// @Success 200 {object} helpers.SuccessResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /admin/settings/{item}  [post]
func UpdateSetting(ctx *app.HttpContext) error {
	var inputModel appmodels.SettingsReqModel

	err := ctx.BlindBind(&inputModel)
	if err != nil {
		return errors.NewBadRequestError(consts.BadRequest, err)
	}
	var (
		key = ctx.GetPathParam("item")
	)

	switch key {
	case "systemSetting":
		err = settings.UpdateSystemSettings(inputModel.Column, inputModel.Value)
	default:
		return errors.NewBadRequestError("Invalid item", nil)
	}

	if err != nil {
		return errors.NewBadRequestError(err.Error(), err)
	}

	return ctx.QuickResponse(consts.OperationDone, http.StatusOK)
}
