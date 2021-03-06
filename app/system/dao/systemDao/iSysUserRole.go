package systemDao

import (
	"baize/app/system/models/systemModels"
)

type IUserRoleDao interface {
	DeleteUserRole(ids []int64)
	BatchUserRole(users []*systemModels.SysUserRole)
	DeleteUserRoleByUserId(userId int64)
	CountUserRoleByRoleId(ids []int64) int
	DeleteUserRoleInfo(userRole *systemModels.SysUserRole)
	DeleteUserRoleInfos(roleId int64 ,userIds []int64)
}
