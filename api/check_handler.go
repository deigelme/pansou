package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"pansou/model"
	"pansou/service"
)

var (
	checkService     *service.CheckService
	checkServiceOnce sync.Once
)

func getCheckService() *service.CheckService {
	checkServiceOnce.Do(func() {
		checkService = service.NewCheckService()
	})
	return checkService
}

func CheckHandler(c *gin.Context) {
	var req model.CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "无效的检测请求: "+err.Error()))
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "items不能为空"))
		return
	}

	// Limit batch size to avoid overloading the check service with too many items at once
	const maxItems = 50
	if len(req.Items) > maxItems {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "items数量超过限制(最多50个)"))
		return
	}

	response := getCheckService().Check(req.Items)
	c.JSON(http.StatusOK, response)
}
