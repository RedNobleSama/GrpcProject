/**
    @auther: oreki
    @date: 2022/4/26
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"user_api/global/utils"
)

func NewLogger() (*zap.Logger, error) {
	timeNow := utils.GetNowTime()
	cfg := zap.NewProductionConfig()                          // 初始化日志，生产环境
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 时间格式
	cfg.EncoderConfig.TimeKey = "time"                        // 时间key
	cfg.OutputPaths = []string{
		fmt.Sprintf("./user_api/log/%s.log", timeNow), // 日志文件路径
		"stdout",
	}
	return cfg.Build()
}

func LoggerInit() {
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger) // 替换全局日志
}
