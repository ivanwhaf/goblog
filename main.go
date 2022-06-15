package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goblog/config"
	"goblog/handlers"
	"goblog/routers"
	"goblog/stores"
)

//
//                       _oo0oo_
//                      o8888888o
//                      88" . "88
//                      (| -_- |)
//                      0\  =  /0
//                    ___/`---'\___
//                  .' \\|     |// '.
//                 / \\|||  :  |||// \
//                / _||||| -:- |||||- \
//               |   | \\\  -  /// |   |
//               | \_|  ''\---/''  |_/ |
//               \  .-\__  '-'  ___/-. /
//             ___'. .'  /--.--\  `. .'___
//          ."" '<  `.___\_<|>_/___.' >' "".
//         | | :  `- \`.;`\ _ /`;.`/ - ` : | |
//         \  \ `_.   \_ __\ /__ _/   .-` /  /
//     =====`-.____`.___ \_____/___.-`___.-'=====
//                       `=---='
//
//     ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//
//               佛祖保佑         永无BUG
//

func main() {
	gin.SetMode(gin.ReleaseMode)
	config.LoadConfig()

	Init()

	router := routers.InitRouter()
	router.LoadHTMLGlob("templates/*")
	router.Static("static", "static")

	err := router.Run(":" + config.GetConfig().Server.Port)
	if err != nil {
		panic(err)
	}
	fmt.Println("go blog start.")
}

func Init() {
	stores.Setup()
	handlers.InitImageBuffer()
	handlers.TokenRoutineStart()
}
