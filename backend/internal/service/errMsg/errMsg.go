package errMsg

const (
	SUCCESS = (0)
	ERROR   = (500)
)

// Admin模块
const (
	ErrorAdminPswError     = 10001
	ErrorAdminUserNotFound = 10002
	ErrorAdminJWT          = 10003
	ErrorTokenInvalid      = 9999
)

var codeMsg = map[uint]string{
	ErrorAdminPswError:     "密码错误",
	ErrorAdminUserNotFound: "用户不存在",
	ErrorTokenInvalid:      "Token已失效",
}

func GetErrMsg(code uint) string {
	return codeMsg[code]
}
