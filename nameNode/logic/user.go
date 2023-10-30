package logic

// func UserLogin(username string, password string) (user *model.User, err error) { //校验用户名密码
// 	return mysql.UserLogin(username, password)
// }

// func UserRegister(username string, password string, email string) (ok bool, err error) {
// 	return mysql.UserRegister(username, password, email)
// }

// func SendSighUpCode(SignUpUser *model.UserRegister) error {
// 	code := util.GenValidateCode(6)
// 	if err := redis.AddSignUpCode(SignUpUser.Username, SignUpUser.Email, code); err != nil {
// 		return err
// 	}
// 	return nil
// }

// //验证用户注册验证码
// func VerifySignUpCode(SignUpUser *model.UserRegister) (bool, error) {
// 	ok, err := redis.VerifySignUpCode(SignUpUser.Username, SignUpUser.Email, SignUpUser.Code)
// 	if err != nil {
// 		return false, err
// 	}
// 	//清理验证码
// 	if ok {
// 		redis.DeleteSignUpKey(SignUpUser.Username, SignUpUser.Email)
// 	}
// 	return ok, err
// }
