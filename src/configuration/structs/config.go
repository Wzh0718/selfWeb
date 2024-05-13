package structs

type Configuration struct {
	JWT      JWT      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap      Zap      `mapstructure:"zap" json:"zap" yaml:"zap"`
	AutoCode AutoCode `mapstructure:"autocode" json:"autocode" yaml:"autocode"`
	// 跨域配置
	Cors Cors `mapstructure:"cors" json:"cors" yaml:"cors"`
	// 系统配置
	System System `mapstructure:"system" json:"system" yaml:"system"`
	Local  Local  `mapstructure:"local" json:"local" yaml:"local"`
}
