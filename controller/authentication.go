package controller

import "github.com/kataras/iris/v12"

/**
 * 认证
 */
func Authentication(context iris.Context) {
	context.Next()
}
