package actors

import "github.com/endiangroup/specstack/config"

func NewDeveloper(configReader config.ConfigReader) Developer {
	return Developer{
		ConfigReader: configReader,
	}
}

type Developer struct {
	ConfigReader config.ConfigReader
}

func (d Developer) ListConfiguration() (string, error) {
	return d.ConfigReader.List()
}
