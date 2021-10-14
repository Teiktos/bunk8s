package parser

import (
	"bunk8s/launcher/model"
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func ParseConfig(fileName string) model.Bunk8sConfig {

	var bunk8sConfig model.Bunk8sConfig

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading YAML file")
	}

	err = yaml.Unmarshal(yamlFile, &bunk8sConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing YAML file")
	}

	return bunk8sConfig
}
