package server

import (
	"funtabs-server/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type saveDataReq struct {
	Data string `json:"data" binding:"required"`
}

// SaveData POST /api/saveData
// 接收 JSON：{ "data": "localStorage 的 JSON 序列化字符串" }
func SaveData(c *gin.Context) {
	userID := c.GetUint("userID")

	var req saveDataReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "data 不能为空")
		return
	}

	// 有则更新，没有则插入（upsert）
	record := model.UserData{
		UserID: userID,
		Data:   req.Data,
	}
	result := model.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"data", "updated_at"}),
	}).Create(&record)

	if result.Error != nil {
		fail(c, 500, "保存失败")
		return
	}

	ok(c, "保存成功")
}

// GetData GET /api/getData
// 返回 data: "localStorage 的 JSON 序列化字符串"
func GetData(c *gin.Context) {
	userID := c.GetUint("userID")

	var record model.UserData
	err := model.DB.Where("user_id = ?", userID).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		// 用户还没有同步过数据
		ok(c, nil)
		return
	}
	if err != nil {
		fail(c, 500, "获取失败")
		return
	}

	// 直接返回字符串，前端会 JSON.parse
	ok(c, record.Data)
}
