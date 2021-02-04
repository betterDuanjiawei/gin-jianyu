# gin-jianyu记录
## 外部包和知识

### cron
*  go get -u github.com/robfig/cron
* 介绍
```
Cron 表达式格式
字段名
是否必填
允许的值
允许的特殊字符
秒（Seconds）
Yes
0-59
* / , -
分（Minutes）
Yes
0-59
* / , -
时（Hours）
Yes
0-23
* / , -
一个月中的某天（Day of month）
Yes
1-31
* / , - ?
月（Month）
Yes
1-12 or JAN-DEC
* / , -
星期几（Day of week）
Yes
0-6 or SUN-SAT
* / , - ?
Cron表达式表示一组时间，使用 6 个空格分隔的字段
可以留意到 Golang 的 Cron 比 Crontab 多了一个秒级，以后遇到秒级要求的时候就省事了
Cron 特殊字符
1、星号 ( * )
星号表示将匹配字段的所有值
2、斜线 ( / )
斜线用户 描述范围的增量，表现为 “N-MAX/x”，first-last/x 的形式，例如 3-59/15 表示此时的第三分钟和此后的每 15 分钟，到59分钟为止。即从 N 开始，使用增量直到该特定范围结束。它不会重复
3、逗号 ( , )
逗号用于分隔列表中的项目。例如，在 Day of week 使用“MON，WED，FRI”将意味着星期一，星期三和星期五
4、连字符 ( - )
连字符用于定义范围。例如，9 - 17 表示从上午 9 点到下午 5 点的每个小时
5、问号 ( ? )
不指定值，用于代替 “ * ”，类似 “ _ ” 的存在，不难理解
预定义的 Cron 时间表
输入
简述
相当于
@yearly (or @annually)
1月1日午夜运行一次
0 0 0 1 1 *
@monthly
每个月的午夜，每个月的第一个月运行一次
0 0 0 1  
@weekly
每周一次，周日午夜运行一次
0 0 0   0
@daily (or @midnight)
每天午夜运行一次
0 0 0   *
@hourly
每小时运行一次
0 0    
```
* 使用场景
```
软删除，同时也引入了另外一个问题
就是我怎么硬删除，我什么时候硬删除？这个往往与业务场景有关系，大致为
另外有一套硬删除接口
定时任务清理（或转移、backup）无效数据
在这里我们选用第二种解决方案来进行实践
```


### gorm
* go get -u github.com/jinzhu/gorm
* go get -u github.com/go-sql-driver/mysql
* gorm callback
```
gorm所支持的回调方法：
创建：BeforeSave、BeforeCreate、AfterCreate、AfterSave
更新：BeforeSave、BeforeUpdate、AfterUpdate、AfterSave
删除：BeforeDelete、AfterDelete
查询：AfterFind
```

```
// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
    if !scope.HasError() {
        nowTime := time.Now().Unix()
        if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
            if createTimeField.IsBlank {
                createTimeField.Set(nowTime)
            }
        }

        if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
            if modifyTimeField.IsBlank {
                modifyTimeField.Set(nowTime)
            }
        }
    }
}
在这段方法中，会完成以下功能
检查是否有含有错误（db.Error）
scope.FieldByName 通过 scope.Fields() 获取所有字段，判断当前是否包含所需字段
for _, field := range scope.Fields() {
  if field.Name == name || field.DBName == name {
      return field, true
  }
  if field.DBName == dbName {
      mostMatchedField = field
  }
}
field.IsBlank 可判断该字段的值是否为空
func isBlank(value reflect.Value) bool {
  switch value.Kind() {
  case reflect.String:
      return value.Len() == 0
  case reflect.Bool:
      return !value.Bool()
  case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
      return value.Int() == 0
  case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
      return value.Uint() == 0
  case reflect.Float32, reflect.Float64:
      return value.Float() == 0
  case reflect.Interface, reflect.Ptr:
      return value.IsNil()
  }

  return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
若为空则 field.Set 用于给该字段设置值，参数为 interface{}
2、updateTimeStampForUpdateCallback
// updateTimeStampForUpdateCallback will set `ModifyTime` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
    if _, ok := scope.Get("gorm:update_column"); !ok {
        scope.SetColumn("ModifiedOn", time.Now().Unix())
    }
}
scope.Get(...) 根据入参获取设置了字面值的参数，例如本文中是 gorm:update_column ，它会去查找含这个字面值的字段属性
scope.SetColumn(...) 假设没有指定 update_column 的字段，我们默认在更新回调设置 ModifiedOn 的值
```

```
1、 我们的Article是如何关联到Tag？
func GetArticle(id int) (article Article) {
    db.Where("id = ?", id).First(&article)
    db.Model(&article).Related(&article.Tag)

    return 
}
能够达到关联，首先是gorm本身做了大量的约定俗成
Article有一个结构体成员是TagID，就是外键。gorm会通过类名+ID的方式去找到这两个类之间的关联关系
Article有一个结构体成员是Tag，就是我们嵌套在Article里的Tag结构体，我们可以通过Related进行关联查询
2、 Preload是什么东西，为什么查询可以得出每一项的关联Tag？
func GetArticles(pageNum int, pageSize int, maps interface {}) (articles []Article) {
    db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

    return
}
Preload就是一个预加载器，它会执行两条SQL，分别是SELECT * FROM blog_articles;和SELECT * FROM blog_tag WHERE id IN (1,2,3,4);，那么在查询出结构后，gorm内部处理对应的映射逻辑，将其填充到Article的Tag中，会特别方便，并且避免了循环查询
那么有没有别的办法呢，大致是两种
gorm的Join
循环Related
综合之下，还是Preload更好
```
* 注意硬删除要使用 Unscoped()，这是 GORM 的约定 `db.Unscoped().Where("deleted_on != ?", 0).Delete(&Article{})`
### go 热更新
[从PHP迁移至Golang - 热更新篇](https://segmentfault.com/a/1190000017228287)
### endless
* 安装  go get -u github.com/fvbock/endless
* Zero downtime restarts for golang HTTP and HTTPS servers. (for golang 1.3+)
```
我们借助 fvbock/endless 来实现 Golang HTTP/HTTPS 服务重新启动的零停机
endless server 监听以下几种信号量：
syscall.SIGHUP：触发 fork 子进程和重新启动
syscall.SIGUSR1/syscall.SIGTSTP：被监听，但不会触发任何动作
syscall.SIGUSR2：触发 hammerTime
syscall.SIGINT/syscall.SIGTERM：触发服务器关闭（会完成正在运行的请求）
endless 正正是依靠监听这些信号量，完成管控的一系列动作
```
* 你想想，每次更新发布、或者修改配置文件等，只需要给该进程发送SIGTERM信号，而不需要强制结束应用，是多么便捷又安全的事！
  问题
* endless 热更新是采取创建子进程后，将原进程退出的方式，这点不符合守护进程的要求

### grace
* go get -u github.com/facebookgo/grace/gracehttp
```
在实际的生产环境中推荐使用以上开源库，关于热更新开源库的使用非常方便，下面是facebook的grace库的例子：
引入github.com/facebookgo/grace/gracehttp包

func main() {
    app := gin.New()// 项目中时候的是gin框架
    router.Route(app)
    var server *http.Server
    server = &http.Server{
        Addr:    ":8080",
        Handler: app,
    }
    gracehttp.Serve(server)
}
```

### http.Server - Shutdown()
* golang > 1.8
```
// http.Server - Shutdown() 热更新版本
func main() {
	router := routers.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("listen failed, err:%v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
```

### swag
* go get -u github.com/swaggo/swag/cmd/swag
* swag -v
* 安装 gin-swagger 
```
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/gin-swagger/swaggerFiles
```
### redisgo
* go get -u github.com/gomodule/redigo/redis"
* [goredis文档](https://pkg.go.dev/github.com/gomodule/redigo/redis)
* [Redis 命令参考](http://doc.redisfans.com/index.html)
```
package gredis

import (
    "encoding/json"
    "time"

    "github.com/gomodule/redigo/redis"

    "github.com/EDDYCJY/go-gin-example/pkg/setting"
)

var RedisConn *redis.Pool

func Setup() error {
    RedisConn = &redis.Pool{
        MaxIdle:     setting.RedisSetting.MaxIdle,
        MaxActive:   setting.RedisSetting.MaxActive,
        IdleTimeout: setting.RedisSetting.IdleTimeout,
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", setting.RedisSetting.Host)
            if err != nil {
                return nil, err
            }
            if setting.RedisSetting.Password != "" {
                if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
                    c.Close()
                    return nil, err
                }
            }
            return c, err
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            _, err := c.Do("PING")
            return err
        },
    }

    return nil
}

func Set(key string, data interface{}, time int) error {
    conn := RedisConn.Get()
    defer conn.Close()

    value, err := json.Marshal(data)
    if err != nil {
        return err
    }

    _, err = conn.Do("SET", key, value)
    if err != nil {
        return err
    }

    _, err = conn.Do("EXPIRE", key, time)
    if err != nil {
        return err
    }

    return nil
}

func Exists(key string) bool {
    conn := RedisConn.Get()
    defer conn.Close()

    exists, err := redis.Bool(conn.Do("EXISTS", key))
    if err != nil {
        return false
    }

    return exists
}

func Get(key string) ([]byte, error) {
    conn := RedisConn.Get()
    defer conn.Close()

    reply, err := redis.Bytes(conn.Do("GET", key))
    if err != nil {
        return nil, err
    }

    return reply, nil
}

func Delete(key string) (bool, error) {
    conn := RedisConn.Get()
    defer conn.Close()

    return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
    conn := RedisConn.Get()
    defer conn.Close()

    keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
    if err != nil {
        return err
    }

    for _, key := range keys {
        _, err = Delete(key)
        if err != nil {
            return err
        }
    }

    return nil
}
```
## 内部包和知识
### time.Now().Format()
* 2006-01-02 15:04:05 // 据说是go诞生之日, 记忆方法:6-1-2-3-4-5
```
// 	golang的时间格式 默认采用的是RFC333 RFC3339     = "2006-01-02T15:04:05Z07:00"
// 这个格式还不能随便写,要不然生成的文件时间不对 正确的固定的格式:// 这是个奇葩,必须是这个时间点, 据说是go诞生之日, 记忆方法:6-1-2-3-4-5
//    return []byte(fmt.Sprintf(`"%s"`, t.Format("2006-01-02 15:04:05"))), nil
//TimeFormat  = "20200601"
TimeFormat = "20060102"
```

### file 相关
```
1、 file.go：
package logging

import (
    "os"
    "time"
    "fmt"
    "log"
)

var (
    LogSavePath = "runtime/logs/"
    LogSaveName = "log"
    LogFileExt = "log"
    TimeFormat = "20060102"
)

func getLogFilePath() string {
    return fmt.Sprintf("%s", LogSavePath)
}

func getLogFileFullPath() string {
    prefixPath := getLogFilePath()
    suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

    return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func openLogFile(filePath string) *os.File {
    _, err := os.Stat(filePath)
    switch {
        case os.IsNotExist(err):
            mkDir()
        case os.IsPermission(err):
            log.Fatalf("Permission :%v", err)
    }

    handle, err := os.OpenFile(filePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Fail to OpenFile :%v", err)
    }

    return handle
}

func mkDir() {
    dir, _ := os.Getwd()
    err := os.MkdirAll(dir + "/" + getLogFilePath(), os.ModePerm)
    if err != nil {
        panic(err)
    }
}
os.Stat：返回文件信息结构描述文件。如果出现错误，会返回*PathError
type PathError struct {
  Op   string
  Path string
  Err  error
}
os.IsNotExist：能够接受ErrNotExist、syscall的一些错误，它会返回一个布尔值，能够得知文件不存在或目录不存在
os.IsPermission：能够接受ErrPermission、syscall的一些错误，它会返回一个布尔值，能够得知权限是否满足
os.OpenFile：调用文件，支持传入文件名称、指定的模式调用文件、文件权限，返回的文件的方法可以用于I/O。如果出现错误，则为*PathError。
const (
    // Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
    O_RDONLY int = syscall.O_RDONLY // 以只读模式打开文件
    O_WRONLY int = syscall.O_WRONLY // 以只写模式打开文件
    O_RDWR   int = syscall.O_RDWR   // 以读写模式打开文件
    // The remaining values may be or'ed in to control behavior.
    O_APPEND int = syscall.O_APPEND // 在写入时将数据追加到文件中
    O_CREATE int = syscall.O_CREAT  // 如果不存在，则创建一个新文件
    O_EXCL   int = syscall.O_EXCL   // 使用O_CREATE时，文件必须不存在
    O_SYNC   int = syscall.O_SYNC   // 同步IO
    O_TRUNC  int = syscall.O_TRUNC  // 如果可以，打开时
)
os.Getwd：返回与当前目录对应的根路径名
os.MkdirAll：创建对应的目录以及所需的子目录，若成功则返回nil，否则返回error
os.ModePerm：const定义ModePerm FileMode = 0777
```

### log 相关

```
log.New：创建一个新的日志记录器。out定义要写入日志数据的IO句柄。prefix定义每个生成的日志行的开头。flag定义了日志记录属性
func New(out io.Writer, prefix string, flag int) *Logger {
  return &Logger{out: out, prefix: prefix, flag: flag}
}

log.LstdFlags：日志记录的格式属性之一，其余的选项如下
const (
  Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
  Ltime                         // the time in the local time zone: 01:23:23
  Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
  Llongfile                     // full file name and line number: /a/b/c/d.go:23
  Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
  LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
  LstdFlags     = Ldate | Ltime // initial values for the standard logger
)
```

### golang 使用 iota
* [golang 使用 iota](https://studygolang.com/articles/2192)

### md5的使用方法
* [golang中字符串MD5生成方式](https://studygolang.com/articles/13463)
```
方案一
func md5V(str string) string  {
    h := md5.New()
    h.Write([]byte(str))
    return hex.EncodeToString(h.Sum(nil))
}
方案二
func md5V2(str string) string {
    data := []byte(str)
    has := md5.Sum(data)
    md5str := fmt.Sprintf("%x", has)
    return md5str
}
方案三
func md5V3(str string) string {
    w := md5.New()
    io.WriteString(w, str)
    md5str := fmt.Sprintf("%x", w.Sum(nil))
    return md5str
}
```

### 获取后缀
* file.GetExt() 内部调用的 path.Ext() 获取文件名
* path.Ext() 获取路径的



