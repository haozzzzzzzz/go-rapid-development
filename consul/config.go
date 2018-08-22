package consul

type ClientConfigFormat struct {
	Address string `json:"address" yaml:"address" validate:"required"`
}
