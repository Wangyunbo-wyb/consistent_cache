package consistent_cache

import (
	"context"
	"errors"
)

var (
	ErrorDataNotExist = errors.New("data not exist") //数据在缓存中不存在
	ErrorCacheMiss    = errors.New("cache miss")     //数据库不存在该数据
	ErrorDBMiss       = errors.New("db miss")
)

const NullData = "Err_Syntax_Null_Data"

// Cache 缓存模块的抽象接口定义
type Cache interface {
	// Enable 启用某个 key 对应读流程写缓存机制（默认情况下为启用状态）
	Enable(ctx context.Context, key string, delayMilis int64) error
	// Disable 禁用某个 key 对应读流程写缓存机制，为了防止出现意外给一个兜底过期时间。
	Disable(ctx context.Context, key string, expireSeconds int64) error
	// Get 读取 key 对应缓存。如果数据在cache中不存在，需要返回指定错误：ErrorCacheMiss
	Get(ctx context.Context, key string) (string, error)
	// Del 删除 key 对应缓存
	Del(ctx context.Context, key string) error
	// PutWhenEnable 校验某个 key 对应读流程写缓存机制是否启用，倘若启用则写入缓存（默认情况下为启用状态）
	PutWhenEnable(ctx context.Context, key, value string, expireSeconds int64) (bool, error)
}

// DB 数据库模块的抽象接口定义
type DB interface {
	// Put 数据写入数据库
	Put(ctx context.Context, obj Object) error
	// Get 从数据库读取数据
	Get(ctx context.Context, obj Object) error
}

// Object 每次读写操作时，操作的一笔数据记录
type Object interface {
	// KeyColumn 获取 key 对应的字段名
	KeyColumn() string
	// Key 获取 key 对应的值
	Key() string

	// Write 将 object 序列化成字符串
	Write() (string, error)
	// Read 读取字符串内容，反序列化到 object 实例中
	Read(body string) error
}

// Logger 日志打印输出模块
type Logger interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}
