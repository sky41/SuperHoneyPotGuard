package main

import (
	"fmt"
	"log"
	"superhoneypotguard/database"
	"superhoneypotguard/models"
	"time"
)

// 提交目的：创建日志清理工具，定期清理operation_logs表历史数据
// 提交内容：实现按时间范围清理日志的功能
// 提交时间：2026-01-19

func main() {
	database.InitDB()

	// 清理30天前的日志
	cutoffTime := time.Now().AddDate(0, 0, -30)

	result := database.DB.Where("created_at < ?", cutoffTime).Delete(&models.OperationLog{})
	if result.Error != nil {
		log.Printf("清理日志失败: %v", result.Error)
		return
	}

	fmt.Printf("成功清理 %d 条 %d 天前的操作日志\n", result.RowsAffected, 30)
}