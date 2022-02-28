package entity

type Storage struct {
	Driver   string `yaml:"driver"`
	Filename string `yaml:"filename"`
}

type Path struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type Config struct {
	Storage []Storage `yaml:"storage"`
	Paths   []Path    `yaml:"paths"`
}
