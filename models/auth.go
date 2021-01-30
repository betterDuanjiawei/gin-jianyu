package models

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) bool {
	var auth Auth
	// 这里的 db是 models包里的,在同一个包可以直接使用, 又因为用不到 models,所以不需要组合 models
	db.Select("id").Where(Auth{
		Username: username,
		Password: password,
	}).First(&auth)
	if auth.ID > 0 {
		return true
	}
	return false
}
