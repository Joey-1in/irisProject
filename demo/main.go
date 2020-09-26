package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/boltdb"
)

var (
	USERNAME = "userName"
	ISLOGIN  = "isLogin"
)

//自定义的struct
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

//自定义的结构体
type Student struct {
	//XMLName xml.Name `xml:"student"`
	StuName string `xml:"stu_name"`
	StuAge  int    `xml:"stu_age"`
}

type Coniguration struct {
	AppName string `json:"appname"`
	Port    int    `json:"port"`
}

func main() {
	app := iris.New()

	app.Get("/getRequest", func(context context.Context) {
		//处理get请求，请求的url为：/getRequest
		path := context.Path()
		app.Logger().Info(path)
		context.WriteString("请求路径：" + path)
	})
	//GET请求
	// http://localhost:8080/hello
	app.Handle("GET", "/hello", func(context context.Context) {
		// 返回html
		// context.HTML("<h1> Hello world. </h1>")
		// 返回json数据格式
		context.JSON(iris.Map{"message": "success", "name": "Hello world", "requestCode": 200})
	})

	//      http://localhost:8002?date=20190310&city=beijing
	//GET： http://localhost:8002/weather/2019-03-10/beijing
	//      http://localhost:8002/weather/2019-03-11/beijing
	//      http://localhost:8002/weather/2019-03-11/tianjin

	app.Get("/weather/{date}/{city}", func(context context.Context) {
		path := context.Path()
		date := context.Params().Get("date")
		city := context.Params().Get("city")
		context.WriteString(path + "  , " + date + " , " + city)
	})
	/**
	 * 1、Get 正则表达式 路由
	 * 使用：context.Params().Get("name") 获取正则表达式变量
	 *
	 */
	// 请求1：/hello/1  /hello/2  /hello/3 /hello/10000
	//正则表达式：{name}
	app.Get("/hello/{name}", func(context context.Context) {
		//获取变量
		path := context.Path()

		app.Logger().Info(path)
		//获取正则表达式变量内容值
		name := context.Params().Get("name")
		context.HTML("<h1>" + name + "</h1>")
	})

	/**
	 * 2、自定义正则表达式变量路由请求 {unit64:uint64}进行变量类型限制
	 */
	// 123
	// davie
	app.Get("/api/users/{userid:uint64}", func(context context.Context) {
		userID, err := context.Params().GetUint("userid")

		if err != nil {
			//设置请求状态码，状态码可以自定义
			context.JSON(map[string]interface{}{
				"requestcode": 201,
				"message":     "bad request",
			})
			return
		}

		context.JSON(map[string]interface{}{
			"requestcode": 200,
			"user_id":     userID,
		})
	})

	//自定义正则表达式路由请求 bool
	app.Get("/api/users/{isLogin:bool}", func(context context.Context) {
		isLogin, err := context.Params().GetBool("isLogin")
		if err != nil {
			context.StatusCode(iris.StatusNonAuthoritativeInfo)
			return
		}
		if isLogin {
			context.WriteString(" 已登录 ")
		} else {
			context.WriteString(" 未登录 ")
		}

		//正则表达式所支持的数据类型
		context.Params()
	})

	//2.处理Get请求 并接受参数
	app.Get("/userinfo", func(context context.Context) {
		path := context.Path()
		app.Logger().Info(path)
		//获取get请求所携带的参数
		userName := context.URLParam("username")
		app.Logger().Info(userName)

		pwd := context.URLParam("pwd")
		app.Logger().Info(pwd)
		//返回html数据格式
		context.HTML("<h1>" + userName + "," + pwd + "</h1>")
	})

	//3.处理Post请求 form表单的字段获取
	app.Post("/postLogin", func(context context.Context) {
		path := context.Path()
		app.Logger().Info(path)
		//context.PostValue方法来获取post请求所提交的for表单数据
		name := context.PostValue("name")
		pwd := context.PostValue("pwd")
		app.Logger().Info(name, "  ", pwd)
		// 返回HTML数据格式
		// context.HTML(name)
		// 返回json数据格式
		context.JSON(iris.Map{"message": "success", "name": name, "requestCode": 200})
	})

	////POST请求
	app.Handle("POST", "/postHello", func(context context.Context) {
		// context.HTML("<h1> This is post request </h1>")
		context.JSON(iris.Map{"message": "success", "name": "Post Hello world", "requestCode": 200})
	})

	//4、处理Post请求 Json格式数据
	/**
	 * Postman工具选择[{"key":"Content-Type","value":"application/json","description":""}]
	 * 请求内容：{"name": "Joey.1in","age": 28}
	 */
	app.Post("/postJson", func(context context.Context) {

		//1.path
		path := context.Path()
		app.Logger().Info("请求URL：", path)

		//2.Json数据解析
		var person Person
		//context.ReadJSON()
		if err := context.ReadJSON(&person); err != nil {
			panic(err.Error())
		}
		// 返回HTML数据格式
		// context.Writef("Received: %#+v\n", person)
		// 返回json数据格式
		context.JSON(iris.Map{"message": "success", "data": person, "requestCode": 200})
	})

	//5.处理Post请求 Xml格式数据
	/**
	 * 请求配置：Content-Type到application/xml（可选但最好设置）
	 * 请求内容：
	 *
	 *  <student>
	 *		<stu_name>Joey.1in</stu_name>
	 *		<stu_age>28</stu_age>
	 *	</student>
	 *
	 */

	app.Post("/postXml", func(context context.Context) {

		//1.Path
		path := context.Path()
		app.Logger().Info("请求URL：", path)

		//2.XML数据解析
		var student Student
		if err := context.ReadXML(&student); err != nil {
			panic(err.Error())
		}
		//输出：
		context.Writef("Received：%#+v\n", student)
	})

	// 模块user
	userParty := app.Party("/user", func(context context.Context) {
		// 处理下一级请求
		context.Next()
	})
	userParty.Done(func(context context.Context) {
		context.Application().Logger().Info("我执行了")
	})
	// /user/register
	userParty.Get("/register", func(context context.Context) {
		app.Logger().Info("用户模块下的注册方法")
	})

	// /user/register
	userParty.Get("/info", func(context context.Context) {
		app.Logger().Info("用户模块下的查询方法")
		// 这里会触发userParty.Done
		context.Next()
	})

	//一、通过程序代码对应用进行全局配置
	app.Configure(iris.WithConfiguration(iris.Configuration{
		//如果设置为true，当人为中断程序执行时，则不会自动正常将服务器关闭。如果设置为true，需要自己自定义处理。
		DisableInterruptHandler: false,
		//该配置项表示更正并将请求的路径重定向到已注册的路径
		//比如：如果请求/home/ 但找不到此Route的处理程序，然后路由器检查/home处理程序是否存在，如果是，（permant）将客户端重定向到正确的路径/home。
		//默认为false
		DisablePathCorrection: false,
		//
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		TimeFormat:                        "Mon,02 Jan 2006 15:04:05 GMT",
		Charset:                           "utf-8",
	}))
	//二、通过读取tml配置文件读取服务配置
	//注意：要在run方法运行之前执行
	// app.Configure(iris.WithConfiguration(iris.TOML("irisProject/demo/configs/iris.tml")))

	//三、通过读取yaml配置文件读取服务配置
	//同样要在run方法运行之前执行
	// app.Configure(iris.WithConfiguration(iris.YAML("irisProject/demo/configs/iris.yml")))

	//四、通过json配置文件进行应用配置
	file, _ := os.Open("D:/Program/Goprogram/irisProject/demo/configs/config.json")
	defer file.Close()

	fmt.Println(file)
	decoder := json.NewDecoder(file)
	conf := Coniguration{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(conf.Port)

	// mvc
	mvc.New(app).Handle(new(UserController))
	mvc.New(app).Handle(new(OrderController))

	// session
	sessionID := "mySession"

	//1、创建session并进行使用
	sess := sessions.New(sessions.Config{
		Cookie: sessionID,
	})

	// 用户登录功能
	app.Post("/login", func(context context.Context) {
		path := context.Path()
		app.Logger().Info(" 请求Path：", path)
		userName := context.PostValue("name")
		passwd := context.PostValue("pwd")

		if userName == "davie" && passwd == "pwd123" {
			session := sess.Start(context)

			//用户名
			session.Set(USERNAME, userName)
			//登录状态
			session.Set(ISLOGIN, true)

			context.WriteString("账户登录成功 ")

		} else {
			session := sess.Start(context)
			session.Set(ISLOGIN, false)
			context.WriteString("账户登录失败，请重新尝试")
		}
	})

	// 用户退出登录功能
	app.Get("/logout", func(context context.Context) {
		path := context.Path()
		app.Logger().Info(" 退出登录 Path :", path)
		session := sess.Start(context)
		//删除session
		session.Delete(ISLOGIN)
		session.Delete(USERNAME)
		context.WriteString("退出登录成功")
	})

	app.Get("/query", func(context context.Context) {
		path := context.Path()
		app.Logger().Info(" 查询信息 path :", path)
		session := sess.Start(context)

		isLogin, err := session.GetBoolean(ISLOGIN)
		if err != nil {
			context.WriteString("账户未登录,请先登录 ")
			return
		}

		if isLogin {
			app.Logger().Info(" 账户已登录 ")
			context.WriteString("账户已登录")
		} else {
			app.Logger().Info(" 账户未登录 ")
			context.WriteString("账户未登录")
		}

	})

	//2、session和db绑定使用
	db, err := boltdb.New("sessions.db", 0600)
	if err != nil {
		panic(err.Error())
	}

	//程序中断时，将数据库关闭
	iris.RegisterOnInterrupt(func() {
		defer db.Close()
	})

	//session和db绑定
	sess.UseDatabase(db)

	app.Run(iris.Addr(":8000"), iris.WithoutServerError(iris.ErrServerClosed))
}

type UserController struct{}
type OrderController struct{}

func (us *UserController) Get() string {
	return "UserController"
}
func (us *UserController) GetInfo() string {
	return "UserController GetInfo"
}
func (us *UserController) Post() string {
	return "OrderController"
}
