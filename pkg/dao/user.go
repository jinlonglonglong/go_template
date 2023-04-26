package dao

import (
	"template/pkg/api/dtos"
	"template/pkg/data"
	"template/pkg/models"
)

func SaveOrUpdateUser(dto dtos.UserDto) (user models.User, ret bool) {
	db := data.MustGetDB("template")
	db.Where("address = ?", dto.Address).First(&user)
	id := user.ID
	if id == 0 {
		result := db.Create(&models.User{
			Address: dto.Address,
		})
		if result.Error != nil {
			ret = false
			return
		}
	}
	ret = true
	return
}
