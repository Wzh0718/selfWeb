package structs

type AutoCode struct {
	Root            string `mapstructure:"root" json:"root" yaml:"root"`
	TransferRestart bool   `mapstructure:"transfer-restart" json:"transfer-restart" yaml:"transfer-restart"`
}
