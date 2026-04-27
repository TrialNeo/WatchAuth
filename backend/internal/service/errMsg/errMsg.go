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

// agent模块
const (
	ERRORAgentNotFound       = 50000
	ERRORAgentParentNotFound = 50001
	ERRORAgentParentFrozen   = 50002
	ERRORAgentLevelLimit     = 50003
	ERRORAgentHasChildren    = 50004
)

// system_config模块
const (
	ERRORConfigNotFound = 60000
)

// announcement模块
const (
	ERRORAnnouncementNotFound = 70000
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

	ERRORAgentNotFound:       "代理不存在",
	ERRORAgentParentNotFound: "上级代理不存在",
	ERRORAgentParentFrozen:   "上级代理已冻结",
	ERRORAgentLevelLimit:     "代理层级已达上限(最多3级)",
	ERRORAgentHasChildren:    "该代理下存在子代理，无法删除",

	ERRORConfigNotFound:        "配置项不存在",
	ERRORAnnouncementNotFound:  "公告不存在",
}

func GetErrMsg(code uint) string {
	return codeMsg[code]
}
