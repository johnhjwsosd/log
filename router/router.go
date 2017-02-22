package router

import (
	"net/http"
	"strconv"

	"errors"

	"../lib"
	"../logger"
	_ "../mongodriver" //...
	"github.com/gin-gonic/gin"
)

var logMap = make(map[string]*logger.Logger)

//GetRouter ...
func GetRouter() *gin.Engine {
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "welcome to logsys")
	})
	router.POST("/", func(c *gin.Context) {
		c.String(http.StatusOK, "welcome to logsys")
	})

	router.GET("/reg", registryLog)
	router.POST("/reg", registryLog)

	router.GET("/wl/", writerLog)
	router.POST("/wl/", writerLog)

	router.GET("/rl/", readLog)
	router.POST("/rl/", readLog)
	return router
}

func writerLog(c *gin.Context) {
	params := getParams(c)
	log := logMap[params["appID"].(string)]
	if log == nil {
		err := verfyRegArgs(params)
		if err != nil {
			c.JSON(404, err.Error())
			return
		}
		regLogObj(params)
	}
	content := getWriteParams(params)
	switch params[""] {
	case "info":
		log.Info(content)
	case "trace":
		log.Trace(content)
	case "debug":
		log.Debug(content)
	case "warn":
		log.Warn(content)
	case "error":
		log.Error(content)
	case "fatal":
		log.Fatal(content)
	default:
		log.Info(content)
	}
	return
}

func readLog(c *gin.Context) {
	params := getParams(c)
	log := logMap[params["appID"].(string)]
	if log == nil {
		err := verfyRegArgs(params)
		if err != nil {
			c.JSON(404, err.Error())
			return
		}
		regLogObj(params)
	}
	where := getReadParams(params)
	res, err := log.Read(where)
	if err != nil {
		c.JSON(404, err.Error())
		return
	}
	c.JSON(200, res)
}

func registryLog(c *gin.Context) {
	params := getParams(c)
	err := verfyRegArgs(params)
	if err != nil {
		c.JSON(404, err.Error())
		return
	}
	appID, err := regLogObj(params)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	appMap := map[string]string{
		"appID": appID,
	}
	c.JSON(200, appMap)
	return
}

func verfyRegArgs(params map[string]interface{}) (err error) {
	if params["stHost"] == nil {
		err = errors.New("缺少参数，存储地址")
		return
	}
	if params["stPort"] == nil {
		err = errors.New("缺少参数，存储端口")
		return
	}
	if params["stName"] == nil {
		err = errors.New("缺少参数，存储名称")
		return
	}
	if params["appName"] == nil {
		err = errors.New("缺少参数，应用名称")
		return
	}
	if params["stType"] == nil {
		err = errors.New("缺少参数，存储类型")
		return
	}
	if params["msName"] == nil {
		err = errors.New("缺少参数，服务名称")
		return
	}
	return nil
}

func getParams(c *gin.Context) (params map[string]interface{}) {
	r := c.Request
	r.ParseForm()
	params = make(map[string]interface{})
	if r.Method == "GET" {
		for k, v := range r.Form {
			if c.Query(k) != "" {
				params[k] = v[0]
			}
		}
	} else {
		for k, v := range r.PostForm {
			params[k] = v[0]
		}
	}
	return
}

func regLogObj(params map[string]interface{}) (string, error) {
	stType := params["stType"].(string)
	stHost := params["stHost"].(string)
	stPort, _ := strconv.Atoi(params["stPort"].(string))
	stName := params["stName"].(string)
	appName := params["appName"].(string)
	msName := params["msName"].(string)
	l, err := logger.NewLogger(stType, stHost, stPort, stName, appName, msName)
	if err != nil {
		return "", err
	}
	uuid := lib.GetUUID()
	logMap[uuid] = l
	return uuid, nil
}

func getWriteParams(params map[string]interface{}) map[string]interface{} {
	c := make(map[string]interface{})
	for k, v := range params {
		if k != "appID" && k != "level" {
			c[k] = v
		}
	}
	return c
}

func getReadParams(params map[string]interface{}) map[string]interface{} {
	c := make(map[string]interface{})
	for k, v := range params {
		if k != "appID" {
			c[k] = v
		}
	}
	return c
}
