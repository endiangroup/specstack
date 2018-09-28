package actors

import "github.com/endiangroup/specstack/config"

func NewDeveloper(configReader config.Reader) Developer {
	return Developer{
		ConfigReader: configReader,
	}
}

type Developer struct {
	ConfigReader config.Reader
}

func (d Developer) ListConfiguration() (string, error) {
	return d.ConfigReader.List()
}
