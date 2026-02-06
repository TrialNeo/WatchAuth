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

// admin_App模块
const (
	ErrorAdminAppNotFound  = 20000
	ErrorAdminAppDelDBFail = 20001
)

var codeMsg = map[uint]string{
	ErrorAdminPswError:     "密码错误",
	ErrorAdminUserNotFound: "用户不存在",
	ErrorTokenInvalid:      "Token已失效",
	ErrorAdminAppDelDBFail: "数据异常",
	ErrorAdminAppNotFound:  "未找到应用",
}

func GetErrMsg(code uint) string {
	return codeMsg[code]
}
