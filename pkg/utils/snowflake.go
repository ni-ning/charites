package utils

import "charites/global"

func GenInt64() int64 {
	// 坑：前端展示不了 int64，需要String()
	return global.SnowNode.Generate().Int64()
}
