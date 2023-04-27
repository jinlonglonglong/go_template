package dao

import (
	"template/pkg/api/dtos"
	"template/pkg/data"
	token "template/pkg/helpers"
	"template/pkg/models"
)

func SaveOrUpdateUser(dto dtos.UserDto) (*dtos.LoginResp, bool) {
	db := data.MustGetDB("template")
	user := &models.User{}
	db.Where("address = ?", dto.Address).First(user)
	id := user.ID
	if id == 0 {
		result := db.Create(&models.User{
			Address: dto.Address,
		})
		if result.Error != nil {
			return nil, false
		}
	}
	//token
	token, _ := token.GenerateToken(dto.Address)
	return &dtos.LoginResp{
		Address: dto.Address,
		Token:   token,
	}, true
}
