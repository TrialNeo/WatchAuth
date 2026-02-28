package errMsg

const (
	SUCCESS            = 0
	ERROR              = 500
	ERRORDataBaseErr   = 501
	ERRORInvalidParams = 502
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

// admin_machine模块
const (
	ErrorMachineNotFound = 40000
)

var codeMsg = map[uint]string{
	ERROR:              "服务器内部错误",
	ERRORInvalidParams: "非法参数",
	ERRORDataBaseErr:   "服务器数据库错误",

	ErrorAdminPswError:     "密码错误",
	ErrorAdminUserNotFound: "用户不存在",
	ErrorTokenInvalid:      "Token已失效",
	ErrorAdminAppDelDBFail: "数据异常",
	ErrorAdminAppNotFound:  "应用不存在",
	ErrorAdminAppVerUsed:   "该版本号已被使用",

	ErrorMachineNotFound: "机器不存在",
}

func GetErrMsg(code uint) string {
	return codeMsg[code]
}
