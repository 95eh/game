package cmn

const (
	RedisGlobal = "global"
)

// 角色组
const (
	RoleGuest int64 = 1 << iota
	RolePlayer
	RoleAdmin
)
