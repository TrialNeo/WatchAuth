package dao

import "Diggpher/global"

// BindDao 对外接口，dao绑定
func BindDao() {
	global.DataBase = global.DataBase.Debug()
	//下面数据库初始化失败就不能继续用了。。
	err := global.DataBase.AutoMigrate(
		new(Admin),
		new(App), new(Version),
		new(Machine), new(MachineInfo), new(UsedApp), new(MachineLog),
		new(License),
		new(User),
		new(Agent),
	)
	if err != nil {
		panic(err)
	}

}
