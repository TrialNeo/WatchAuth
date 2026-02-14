package errMsg

const (
	SUCCESS          = (0)
	ERROR            = 500
	ERRORDataBaseErr = 501
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
	ErrorAdminAppVerUsed   = 20002
)

var codeMsg = map[uint]string{
	ErrorAdminPswError:     "密码错误",
	ErrorAdminUserNotFound: "用户不存在",
	ErrorTokenInvalid:      "Token已失效",
	ErrorAdminAppDelDBFail: "数据异常",
	ErrorAdminAppNotFound:  "应用不存在",
	ERRORDataBaseErr:       "内部数据库错误",
	ErrorAdminAppVerUsed:   "该版本号已被使用",
}

func GetErrMsg(code uint) string {
	return codeMsg[code]
}
