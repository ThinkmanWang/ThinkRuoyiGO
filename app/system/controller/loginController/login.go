package loginController

import (
	"baize/app/common/commonController"
	"baize/app/common/commonModels"
	"baize/app/monitor/monitorService"
	"baize/app/monitor/monitorService/monitorServiceImpl"
	"baize/app/system/models/loginModels"
	"baize/app/system/service/loginService"
	"baize/app/system/service/loginService/loginServiceImpl"
	"baize/app/system/service/systemService"
	"baize/app/system/service/systemService/systemServiceImpl"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var iUserOnline monitorService.ItUserOnlineService = monitorServiceImpl.GetUserOnlineService()
var iLogin loginService.ILoginService = loginServiceImpl.GetLoginService()
var iMenu systemService.IMenuService = systemServiceImpl.GetMenuService()

func Login(c *gin.Context) {
	var login loginModels.LoginBody
	if err := c.ShouldBindJSON(&login); err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		c.JSON(http.StatusOK, commonModels.ParameterError())
		return
	}
	data := iLogin.Login(&login, c)

	c.JSON(http.StatusOK, data)

}
func GetInfo(c *gin.Context) {
	loginUser := commonController.GetCurrentLoginUser(c)

	data := make(map[string]interface{})
	data["user"] = loginUser.User
	data["roles"] = loginUser.RolePerms
	data["permissions"] = loginUser.Permissions
	c.JSON(http.StatusOK, commonModels.SuccessData(data))

}
func GetRouters(c *gin.Context) {
	loginUser := commonController.GetCurrentLoginUser(c)
	menus := iMenu.SelectMenuTreeByUserId(loginUser.User.UserId)
	buildMenus := iMenu.BuildMenus(menus)
	c.JSON(http.StatusOK, commonModels.SuccessData(buildMenus))

}
func Logout(c *gin.Context) {
	loginUser := commonController.GetCurrentLoginUser(c)
	if loginUser != nil {
		iUserOnline.ForceLogout(loginUser.Token)
	}
	c.JSON(http.StatusOK, commonModels.Success())

}
