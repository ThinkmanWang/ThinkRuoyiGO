package RoleController

import (
	"baize/app/common/commonController"
	"baize/app/common/commonLog"
	"baize/app/common/commonModels"
	"baize/app/system/models/systemModels"
	"baize/app/system/service/systemService"
	"baize/app/system/service/systemService/systemServiceImpl"
	"baize/app/utils/slicesUtils"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/util/gconv"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

var iRole systemService.IRoleService = systemServiceImpl.GetRoleService()

func RoleList(c *gin.Context) {
	loginUser := commonController.GetCurrentLoginUser(c)
	role := new(systemModels.SysRoleDQL)
	c.ShouldBind(role)
	role.SetLimit(c)
	role.SetDataScope(loginUser, "d", "")
	list, count := iRole.SelectRoleList(role)

	c.JSON(http.StatusOK, commonModels.SuccessListData(list, count))

}

func RoleExport(c *gin.Context) {
	loginUser := commonController.GetCurrentLoginUser(c)
	role := new(systemModels.SysRoleDQL)
	c.ShouldBind(role)
	role.SetDataScope(loginUser, "d", "")
	data := iRole.RoleExport(role)
	commonController.DataPackageExcel(c,data)

}
func RoleGetInfo(c *gin.Context) {
	roleId, err := strconv.ParseInt(c.Param("roleId"), 10, 64)
	if err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		c.JSON(http.StatusOK, commonModels.ParameterError())
		return
	}
	sysUser := iRole.SelectRoleById(roleId)

	c.JSON(http.StatusOK, commonModels.SuccessData(sysUser))
}
func RoleAdd(c *gin.Context) {
	commonLog.SetLog(c, "角色管理", "INSERT")
	loginUser := commonController.GetCurrentLoginUser(c)
	sysRole := new(systemModels.SysRoleDML)
	if err := c.ShouldBindJSON(sysRole); err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		c.JSON(http.StatusOK, commonModels.ParameterError())
		return
	}
	if iRole.CheckRoleNameUnique(sysRole) {
		c.JSON(http.StatusOK, commonModels.Waring("新增角色'"+sysRole.RoleName+"'失败，角色名称已存在"))
		return
	}
	if iRole.CheckRoleKeyUnique(sysRole) {
		c.JSON(http.StatusOK, commonModels.Waring("新增角色'"+sysRole.RoleKey+"'失败，角色权限已存在"))
		return
	}
	sysRole.SetCreateBy(loginUser.User.UserName)
	iRole.InsertRole(sysRole)
	c.JSON(http.StatusOK, commonModels.Success())

}
func RoleEdit(c *gin.Context) {
	commonLog.SetLog(c, "角色管理", "UPDATE")
	loginUser := commonController.GetCurrentLoginUser(c)
	sysRole := new(systemModels.SysRoleDML)
	if err := c.ShouldBindJSON(sysRole); err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		c.JSON(http.StatusOK, commonModels.ParameterError())
		return
	}
	if iRole.CheckRoleNameUnique(sysRole) {
		c.JSON(http.StatusOK, commonModels.Waring("新增角色'"+sysRole.RoleName+"'失败，角色名称已存在"))
		return
	}
	if iRole.CheckRoleKeyUnique(sysRole) {
		c.JSON(http.StatusOK, commonModels.Waring("新增角色'"+sysRole.RoleKey+"'失败，角色权限已存在"))
		return
	}
	sysRole.SetUpdateBy(loginUser.User.UserName)
	iRole.UpdateRole(sysRole)
	c.JSON(http.StatusOK, commonModels.Success())

}
func RoleDataScope(c *gin.Context) {
	commonLog.SetLog(c, "角色管理", "UPDATE")
	loginUser := commonController.GetCurrentLoginUser(c)
	sysRole := new(systemModels.SysRoleDML)
	c.ShouldBindJSON(sysRole)
	sysRole.SetUpdateBy(loginUser.User.UserName)
	iRole.AuthDataScope(sysRole)
	c.JSON(http.StatusOK, commonModels.Success())

}
func RoleChangeStatus(c *gin.Context) {
	commonLog.SetLog(c, "角色管理", "UPDATE")
	loginUser := commonController.GetCurrentLoginUser(c)
	sysRole := new(systemModels.SysRoleDML)
	c.ShouldBindJSON(sysRole)

	sysRole.SetUpdateBy(loginUser.User.UserName)
	iRole.UpdateRoleStatus(sysRole)

	c.JSON(http.StatusOK, commonModels.Success())
}
func RoleRemove(c *gin.Context) {
	commonLog.SetLog(c, "角色管理", "DELETE")
	var s slicesUtils.Slices = strings.Split(c.Param("rolesIds"), ",")
	ids := s.StrSlicesToInt()
	if iRole.CountUserRoleByRoleId(ids) {
		c.JSON(http.StatusOK, commonModels.Waring("角色已分配，不能删除"))
		return
	}
	iRole.DeleteRoleByIds(ids)
	c.JSON(http.StatusOK, commonModels.Success())
}
func AllocatedList(c *gin.Context) {
	loginUser := commonController.GetCurrentLoginUser(c)
	user := new(systemModels.SysRoleAndUserDQL)
	c.ShouldBind(user)
	user.SetLimit(c)
	user.SetDataScope(loginUser, "d", "u")
	list, count := iRole.SelectAllocatedList(user)
	c.JSON(http.StatusOK, commonModels.SuccessListData(list, count))

}
func UnallocatedList(c *gin.Context) {
	loginUser := commonController.GetCurrentLoginUser(c)
	user := new(systemModels.SysRoleAndUserDQL)
	c.ShouldBind(user)
	user.SetLimit(c)
	user.SetDataScope(loginUser, "d", "u")
	list, count := iRole.SelectUnallocatedList(user)
	c.JSON(http.StatusOK, commonModels.SuccessListData(list, count))

}
func InsertAuthUser(c *gin.Context) {
	commonLog.SetLog(c, "角色管理", "GRANT")
	var userIds slicesUtils.Slices = strings.Split(c.Query("userIds"), ",")
	roleId := c.Query("roleId")
	iRole.InsertAuthUsers(gconv.Int64(roleId), userIds.StrSlicesToInt())
	c.JSON(http.StatusOK, commonModels.Success())
	return
}
func CancelAuthUser(c *gin.Context) {
	commonLog.SetLog(c, "角色管理", "GRANT")
	userRole := new(systemModels.SysUserRole)
	c.ShouldBindJSON(userRole)
	iRole.DeleteAuthUserRole(userRole)
	c.JSON(http.StatusOK, commonModels.Success())
	return
}
func CancelAuthUserAll(c *gin.Context) {
	commonLog.SetLog(c, "角色管理", "GRANT")
	var userIds slicesUtils.Slices = strings.Split(c.Query("userIds"), ",")
	roleId := c.Query("roleId")
	iRole.DeleteAuthUsers(gconv.Int64(roleId), userIds.StrSlicesToInt())
	c.JSON(http.StatusOK, commonModels.Success())
	return
}
