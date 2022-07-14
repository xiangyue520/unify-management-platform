package initialize

import (
	"encoding/json"
	"fmt"
	common_index_router "mayfly-go/internal/common/router"
	devops_router "mayfly-go/internal/devops/router"
	"mayfly-go/internal/sys/domain/entity"
	sys_router "mayfly-go/internal/sys/router"
	"mayfly-go/pkg/config"
	"mayfly-go/pkg/middleware"
	"mayfly-go/pkg/model"
	"mayfly-go/pkg/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// server配置
	serverConfig := config.Conf.Server
	gin.SetMode(serverConfig.Model)

	var router = gin.New()
	router.MaxMultipartMemory = 8 << 20

	// 没有路由即 404返回
	router.NoRoute(func(g *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				g.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": fmt.Sprintf("请求出错 '%s'", err)})
			}
		}()
		method := g.Request.Method
		path := g.Request.URL.Path
		// 查找是否配置了该路径的mock数据，存在则直接返回mock数据
		mockData := &entity.MockData{}
		mockData.Path = path
		mockData.Method = method
		err := model.GetBy(mockData)
		if err != nil || mockData.Status == -1 {
			if mockData.RealUrl != "" {
				reqUrl := g.Request.URL
				realUrl, _ := url.ParseRequestURI(mockData.RealUrl)
				realUrl.Path = reqUrl.Path
				realUrl.RawQuery = reqUrl.RawQuery

				httputil.NewSingleHostReverseProxy(realUrl).ServeHTTP(g.Writer, g.Request)
				return
			}
			g.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": fmt.Sprintf("not found '%s:%s'", method, path)})
			return
		}

		vars := make(map[string]map[string]interface{})

		rowQuery := g.Request.URL.RawQuery
		if rowQuery != "" {
			queryMap := make(map[string]interface{})
			pars := strings.Split(rowQuery, "&")
			for _, par := range pars {
				// 然后分解 sessionid=222222 和 token=3333
				parkv := strings.Split(par, "=")
				// 等号前面是key,后面是value
				queryMap[parkv[0]] = parkv[1]
			}
			vars["query"] = queryMap
		}

		reqBody := make(map[string]interface{})
		if strings.HasSuffix(g.ContentType(), "application/json") {
			g.BindJSON(&reqBody)
		} else {
			g.Request.ParseMultipartForm(100000000)
			postForm := g.Request.PostForm
			if len(postForm) > 0 {
				for k, v := range postForm {
					reqBody[k] = v[0]
				}
			}
		}
		vars["body"] = reqBody

		mockRes := utils.TemplateParse(mockData.Data, vars)

		var res map[string]interface{}
		err = json.Unmarshal([]byte(mockRes), &res)
		if err != nil {
			g.String(http.StatusOK, mockRes)
		} else {
			g.JSON(http.StatusOK, res)
		}
	})

	// 设置静态资源
	if staticConfs := serverConfig.Static; staticConfs != nil {
		for _, scs := range *staticConfs {
			router.Static(scs.RelativePath, scs.Root)
		}

	}
	// 设置静态文件
	if staticFileConfs := serverConfig.StaticFile; staticFileConfs != nil {
		for _, sfs := range *staticFileConfs {
			router.StaticFile(sfs.RelativePath, sfs.Filepath)
		}
	}
	// 是否允许跨域
	if serverConfig.Cors {
		router.Use(middleware.Cors())
	}

	// 设置路由组
	api := router.Group("/api")
	{
		common_index_router.InitIndexRouter(api)

		sys_router.InitCaptchaRouter(api)
		sys_router.InitAccountRouter(api) // 注册account路由
		sys_router.InitResourceRouter(api)
		sys_router.InitRoleRouter(api)
		sys_router.InitSystemRouter(api)
		sys_router.InitMockRouter(api)
		sys_router.InitSyslogRouter(api)

		devops_router.InitProjectRouter(api)
		devops_router.InitDbRouter(api)
		devops_router.InitDbSqlExecRouter(api)
		devops_router.InitRedisRouter(api)
		devops_router.InitMachineRouter(api)
		devops_router.InitMachineScriptRouter(api)
		devops_router.InitMachineFileRouter(api)
		devops_router.InitMongoRouter(api)
	}

	return router
}
