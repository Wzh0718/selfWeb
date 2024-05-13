package structs

type Configuration struct {
	JWT      JWT      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap      Zap      `mapstructure:"zap" json:"zap" yaml:"zap"`
	AutoCode AutoCode `mapstructure:"autocode" json:"autocode" yaml:"autocode"`
}
