module github.com/betterDuanjiawei/gin-jianyu

go 1.15

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/astaxie/beego v1.12.3
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fvbock/endless v0.0.0-20170109170031-447134032cb6
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ini/ini v1.62.0
	github.com/go-openapi/spec v0.20.2 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/robfig/cron v1.2.0
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.0
	github.com/ugorji/go v1.2.3 // indirect
	github.com/unknwon/com v1.0.1
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	golang.org/x/tools v0.1.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// 首先你要看到我们使用的是完整的外部模块引用路径（github.com/EDDYCJY/go-gin-example/xxx），而这个模块还没推送到远程，是没有办法下载下来的，因此需要用 replace 将其指定读取本地的模块路径，这样子就可以解决本地模块读取的问题。
replace (
	github.com/betterDuanjiawei/gin-jianyu/conf => /Users/v_duanjiawei/go/src/github.com/betterDuanjiawei/gin-jianyu/conf
	github.com/betterDuanjiawei/gin-jianyu/middleware => /Users/v_duanjiawei/go/src/github.com/betterDuanjiawei/gin-jianyu/middleware
	github.com/betterDuanjiawei/gin-jianyu/models => /Users/v_duanjiawei/go/src/github.com/betterDuanjiawei/gin-jianyu/models
	github.com/betterDuanjiawei/gin-jianyu/pkg/e => /Users/v_duanjiawei/go/src/github.com/betterDuanjiawei/gin-jianyu/pkg/e
	github.com/betterDuanjiawei/gin-jianyu/pkg/setting => /Users/v_duanjiawei/go/src/github.com/betterDuanjiawei/gin-jianyu/pkg/setting
	github.com/betterDuanjiawei/gin-jianyu/pkg/util => /Users/v_duanjiawei/go/src/github.com/betterDuanjiawei/gin-jianyu/pkg/util
	github.com/betterDuanjiawei/gin-jianyu/router => /Users/v_duanjiawei/go/src/github.com/betterDuanjiawei/gin-jianyu/router
	github.com/betterDuanjiawei/gin-jianyu/runtime => /Users/v_duanjiawei/go/src/github.com/betterDuanjiawei/gin-jianyu/runtime

)
