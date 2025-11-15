package db

import (
	"fmt"
	"log"
	"oj-back/app/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	var err error
	const dsn string = "host=postgres user=postgres password=secretpassword dbname=mydatabase sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	fmt.Println("数据库连接成功")
	migrator := DB.Migrator()
	// 清空所有表
	if err := migrator.DropTable(&model.Problem{}, &model.TestCase{}); err != nil {
		log.Fatalf("删除表失败: %v", err)
	}
	// 重新创建所有表
	if err := migrator.CreateTable(&model.Problem{}, &model.TestCase{}); err != nil {
		log.Fatalf("创建表失败: %v", err)
	}
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_problems_tags_gin ON problems USING GIN (tags);`)
	model.InsertData(DB)
	return DB
}
