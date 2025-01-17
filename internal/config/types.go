package config

// ======================== Config ======================== //

type Config struct {
	PersistFilePath string        `mapstructure:"persist_file_path" yaml:"persist_file_path"`
	Flags           *FlagsConfig  `mapstructure:"flags" yaml:"flags"`
	System          *SystemConfig `mapstructure:"system" yaml:"system"`
	App             *AppConfig    `mapstructure:"app" yaml:"app"`
}

// ======================== Flags ======================== //

type FlagsConfig struct {
	Environment string `mapstructure:"environment" yaml:"environment"`
	DebugMode   bool   `mapstructure:"debug_mode" yaml:"debug_mode"`
}

// ======================== System ======================== //

type SystemConfig struct {
	AppName      string   `mapstructure:"app_name" yaml:"app_name"`
	AppVersion   string   `mapstructure:"app_version" yaml:"app_version"`
	BuildNumber  int      `mapstructure:"build_number" yaml:"build_number"`
	ReleaseDate  string   `mapstructure:"release_date" yaml:"release_date"`
	Contributors []string `mapstructure:"contributors" yaml:"contributors"`
}

// ======================== App ======================== //

type AppConfig struct {
	Logging LoggingConfig `mapstructure:"logging" yaml:"logging"`
	Mqtt    MqttConfig    `mapstructure:"mqtt" yaml:"mqtt"`
}

type LoggingConfig struct {
	Level      string `mapstructure:"level" yaml:"level"`
	FilePath   string `mapstructure:"file_path" yaml:"file_path"`
	MaxSize    int    `mapstructure:"max_size" yaml:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" yaml:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" yaml:"max_age"`
	Compress   bool   `mapstructure:"compress" yaml:"compress"`
	AddTime    bool   `mapstructure:"add_time" yaml:"add_time"`
}

type MqttConfig struct {
	Broker             string `mapstructure:"broker" yaml:"broker"`
	ClientId           string `mapstructure:"client_id" yaml:"client_id"`
	Port               int    `mapstructure:"port" yaml:"port"`
	Topic              string `mapstructure:"topic" yaml:"topic"`
	Qos                byte   `mapstructure:"qos" yaml:"qos"`
	CleanSession       bool   `mapstructure:"clean_session" yaml:"clean_session"`
	KeepAlive          int    `mapstructure:"keep_alive" yaml:"keep_alive"`
	ReconnectOnFailure bool   `mapstructure:"reconnect_on_failure" yaml:"reconnect_on_failure"`
	Username           string `mapstructure:"username" yaml:"username"`
	Password           string `mapstructure:"password" yaml:"password"`
}
