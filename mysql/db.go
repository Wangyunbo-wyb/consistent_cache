package mysql

import (
	"consistent_cache"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// 判断操作模型是否声明了表名
type tabler interface {
	TableName() string
}

// DB 数据库模块的抽象接口定义
type DB struct {
	db *gorm.DB
}

func NewDB(dsn string) *DB {
	return &DB{db: getDB(dsn)}
}

// Put 数据写入数据库
func (d *DB) Put(ctx context.Context, obj consistent_cache.Object) error {
	db := d.db
	//倘若obj显示声明了表名，则进行应用
	tabler, ok := obj.(tabler)
	if ok {
		db = db.Table(tabler.TableName())
	}

	// 此处通过两个非原子性动作实现 upsert 效果：
	// 1 尝试创建记录
	// 2 倘若发生唯一键冲突，则改为执行更新操作
	err := db.WithContext(ctx).Create(obj).Error
	if err == nil {
		return nil
	}

	// 判断是否为唯一键冲突，若是的话，则改为更新操作
	if IsDuplicateEntryErr(err) {
		return db.WithContext(ctx).Debug().Where(fmt.Sprintf("`%s` = ?", obj.KeyColumn()), obj.Key()).Updates(obj).Error
	}
	// 其他错误直接返回
	return err
}

// Get 从数据库读取数据
func (d *DB) Get(ctx context.Context, obj consistent_cache.Object) error {
	db := d.db
	//倘若obj显示声明了表名，则进行应用
	tabler, ok := obj.(tabler)
	if ok {
		db = db.Table(tabler.TableName())
	}
	//select 语句读取通过唯一键检索数据记录
	err := db.WithContext(ctx).Where(fmt.Sprintf("`%s` = ?", obj.KeyColumn()), obj.Key()).First(obj).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return consistent_cache.ErrorDBMiss
	}
	return err
}
