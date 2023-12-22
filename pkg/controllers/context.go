package controllers

import (
	"github.com/DanArmor/vtuber-go/pkg/controllers/types"
	"github.com/gin-gonic/gin"
)

func getInitData(c *gin.Context) types.InitData {
	v, _ := c.Get("initData")
	return v.(types.InitData)
}
