package model

type User struct { //数据库User表
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type UserLogin struct { //用户登录参数
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegister struct { //用户注册参数
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Code     string `json:"code" defult:""`
}
