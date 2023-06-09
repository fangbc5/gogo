package mysql

import (
	"github.com/fangbc5/gogo/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var mysqlClient *gorm.DB

// Init 数据库
func Init(opts ...Option) {
	options := NewOptions(opts...)
	dsn := options.Username + ":" + options.Password + "@tcp(" + options.Address + ":" + options.Port + ")" + "/" + options.Database + "?charset=utf8&parseTime=True&loc=Local"
	//创建连接
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		SkipDefaultTransaction: false, //跳过默认事务
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 复数形式 User的表名应该是users
			TablePrefix:   "",   //表名前缀 User的表名应该是t_users
		},
		DisableForeignKeyConstraintWhenMigrating: true, //设置成为逻辑外键(在物理数据库上没有外键，仅体现在代码上)
	})
	if err != nil {
		log.Panicf("database load err %v\n", err)
	}
	//设置连接池
	sqldb, err := db.DB()
	if err != nil {
		log.Panicf("database load err %v\n", err)
	}
	sqldb.SetMaxIdleConns(options.MaxIdleConns)
	sqldb.SetMaxOpenConns(options.MaxOpenConns)
	sqldb.SetConnMaxIdleTime(options.ConnMaxIdleTime)
	sqldb.SetConnMaxLifetime(options.ConnMaxLifetime)

	//自动生成表
	//err = db.AutoMigrate(&model.User{})
	//if err != nil {
	//	log.Println(err)
	//}

	//全局db对象
	mysqlClient = db
	if utils.IsNotNull(mysqlClient) {
		log.Println("mysql connect success")
	}
}

func GetGormApi() *gorm.DB {
	return mysqlClient
}
