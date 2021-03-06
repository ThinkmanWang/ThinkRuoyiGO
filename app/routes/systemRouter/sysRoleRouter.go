package systemRouter

import (
	"baize/app/common/middlewares"
	"baize/app/system/controller/RoleController"
	"github.com/gin-gonic/gin"
)

func InitSysRoleRouter(router *gin.RouterGroup) {
	role := router.Group("/system/role")
	role.GET("/list", middlewares.HasPermission("system:role:list"), RoleController.RoleList)
	role.POST("/export", middlewares.HasPermission("system:role:export"), RoleController.RoleExport)
	role.GET("/:roleId", middlewares.HasPermission("system:role:query"), RoleController.RoleGetInfo)
	role.POST("", middlewares.HasPermission("system:role:add"), RoleController.RoleAdd)
	role.PUT("", middlewares.HasPermission("system:role:edit"), RoleController.RoleEdit)
	role.PUT("/dataScope", middlewares.HasPermission("system:role:edit"), RoleController.RoleDataScope)
	role.PUT("/changeStatus", middlewares.HasPermission("system:role:edit"), RoleController.RoleChangeStatus)
	role.DELETE("/:userIds", middlewares.HasPermission("system:role:remove"), RoleController.RoleRemove)
	role.GET("/authUser/allocatedList", middlewares.HasPermission("system:role:list"), RoleController.AllocatedList)
	role.GET("/authUser/unallocatedList", middlewares.HasPermission("system:role:list"), RoleController.UnallocatedList)
	role.PUT("/authUser/selectAll", middlewares.HasPermission("system:role:edit"), RoleController.InsertAuthUser)
	role.PUT("/authUser/cancelAll", middlewares.HasPermission("system:role:edit"), RoleController.CancelAuthUserAll)
	role.PUT("/authUser/cancel", middlewares.HasPermission("system:role:edit"), RoleController.CancelAuthUser)

}
