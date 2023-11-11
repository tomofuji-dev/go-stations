package env

import (
	"fmt"

	"github.com/spf13/viper"
)

type Env struct {
	Port          string `mapstructure:"PORT"`
	DBPath        string `mapstructure:"DB_PATH"`
	BasicAuthId   string `mapstructure:"BASIC_AUTH_USER_ID"`
	BasicAuthPass string `mapstructure:"BASIC_AUTH_PASSWORD"`
}

// デフォルト値
func (e *Env) DefaultValues() map[string]interface{} {
	return map[string]interface{}{
		"PORT":    ":8080",
		"DB_PATH": ".sqlite3/todo.db",
	}
}

// 必須フィールド
func (e *Env) RequiredFields() []string {
	return []string{"BASIC_AUTH_USER_ID", "BASIC_AUTH_PASSWORD"}
}

func GetEnv() (*Env, error) {
	env := Env{}

	viper.SetConfigFile(".env")

	// デフォルト値の設定
	for key, value := range env.DefaultValues() {
		viper.SetDefault(key, value)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		return nil, err
	}

	// 必須フィールドの確認
	for _, key := range env.RequiredFields() {
		if !viper.IsSet(key) || viper.GetString(key) == "" {
			return nil, fmt.Errorf(".env error: %s is not set", key)
		}
	}

	return &env, nil
}
