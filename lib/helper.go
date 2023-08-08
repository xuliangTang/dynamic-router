package lib

import "github.com/gin-gonic/gin"

// 判断路由是否存在
func hasRoute(method, path string, info gin.RoutesInfo) bool {
	for _, r := range info {
		if r.Method == method && r.Path == path {
			return true
		}
	}
	return false
}
