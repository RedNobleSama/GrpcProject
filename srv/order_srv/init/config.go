/**
    @auther: oreki
    @date: 2022年09月28日 3:39 PM
    @note: 图灵老祖保佑,永无BUG
**/

package init

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Ip   string `mapstructure:"ip"`
	Port int64  `mapstructure:"port"`
}

type Config struct {
	ServerInfo ServerConfig
}

// GetEnvInfo 获取电脑环境变量
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

// WatchConfig viper监控文件变化
func WatchConfig(v *viper.Viper) {
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&ServerConfig{})
		fmt.Println(ServerConfig{})
	})
}
