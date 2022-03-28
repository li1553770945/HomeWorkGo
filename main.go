package main

import (
	"HomeWorkGo/dao"
	"HomeWorkGo/models"
	"HomeWorkGo/routers"
	"HomeWorkGo/setting"
	"fmt"
)

func main() {

	// 加载配置文件
	if err := setting.Init("conf/config.ini"); err != nil {
		fmt.Printf("load config from file failed, err:%v\n", err)
		return
	}
	// 创建数据库
	// 连接数据库
	err := dao.InitRedis(setting.Conf.RedisConfig)
	if err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	err = dao.InitMySQL(setting.Conf.MySQLConfig)
	if err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	err = dao.DB.SetupJoinTable(&models.GroupModel{}, "Members", &models.GroupMemberModel{})
	if err != nil {
		fmt.Printf("migrate failed, err:%v\n", err)
		return
	}
	err = dao.DB.AutoMigrate(&models.UserModel{}, &models.GroupModel{}, &models.HomeWorkModel{}, &models.SubmissionModel{})
	if err != nil {
		fmt.Printf("migrate failed, err:%v\n", err)
		return
	}

	// 注册路由
	r := routers.SetupRouter()

	if err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port)); err != nil {
		fmt.Printf("server startup failed, err:%v\n", err)
	}
}
