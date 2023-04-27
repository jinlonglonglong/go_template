package dtos

type UserDto struct {
	Address string `json:"address"`
}

type LoginResp struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}
