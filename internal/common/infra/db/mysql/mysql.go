package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLDB struct {
	DB *gorm.DB
}

type MySQLConnectionOptions struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func NewMySQLDBConnection() *MySQLDB {
	return &MySQLDB{}
}

func (d *MySQLDB) Connect(options MySQLConnectionOptions) error {
	dsn := options.User + ":" + options.Password + "@tcp(" + options.Host + ":" + options.Port + ")/" + options.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	d.DB = db
	return nil
}

func (d *MySQLDB) Close() error {
	db, err := d.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
