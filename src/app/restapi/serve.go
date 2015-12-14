/**
 * Copyright 2015 @ z3q.net.
 * name : rest_server.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package restapi

import (
	"github.com/jsix/gof"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"go2o/src/core/variable"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	API_DOMAIN   string
	API_HOST_CHK bool = false // 必须匹配Host
	PathPrefix        = "/go2o_api_v1"
	sto          gof.Storage
)

func Run(app gof.App, port int) {
	log.Println("** [ Go2o][ API][ Booted] - Api server running on port " + strconv.Itoa(port))
	API_DOMAIN = app.Config().GetString(variable.ApiDomain)
	sto = app.Storage()
	s := echo.New()
	s.Use(mw.Recover())
	s.Use(beforeRequest())
	s.Hook(splitPath) // 获取新的路径,在请求之前发生
	registerRoutes(s)
	s.Run(":" + strconv.Itoa(port)) //启动服务
}

func registerRoutes(s *echo.Echo) {
	pc := &partnerC{}
	mc := &MemberC{}
	gc := &getC{}

	s.Get("/", ApiTest)
	s.Get("/get/invite_qr", gc.Invite_qr) // 获取二维码
	s.Post("/mm_login", mc.Login)         // 会员登陆接口
	s.Post("/mm_register", mc.Register)   // 会员注册接口
	s.Post("/partner/get_ad", pc.Get_ad)  // 商户广告接口
	//s.Post("/member/*",mc)  // 会员接口
}

func beforeRequest() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			host := ctx.Request().URL.Host
			path := ctx.Request().URL.Path
			// todo: path compare
			if API_HOST_CHK && host != API_DOMAIN {
				return ctx.String(http.StatusNotFound, "no such file")
			}

			if path != "/" {
				//检查商户接口权限
				ctx.Request().ParseForm()
				if !chkPartnerApiSecret(ctx) {
					return ctx.String(http.StatusOK, "{error:'secret incorrent'}")
				}
				//检查会员会话
				if strings.HasPrefix(path, "/member") && !checkMemberToken(ctx) {
					return ctx.String(http.StatusOK, "{error:'incorrent session'}")
				}
			}

			return h(ctx)
		}
	}
}

func splitPath(w http.ResponseWriter, r *http.Request) {
	preLen := len(PathPrefix)
	if len(r.URL.Path) > preLen {
		r.URL.Path = r.URL.Path[preLen:]
	}
}
