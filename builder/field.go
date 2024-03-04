package builder

type Before struct {
	Command []string `yaml:"command"`
}
type Checksum struct {
	File string `yaml:"file"`
}
type Archives struct {
	Enable bool     `yaml:"enable"`
	Name   string   `yaml:"name"`
	Format string   `yaml:"format"`
	Files  []string `yaml:"files"`
}
type After struct {
	Command []string `yaml:"command"`
}
