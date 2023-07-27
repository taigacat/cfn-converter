package converter

import "gopkg.in/yaml.v3"

type Converter interface {
	Convert(*yaml.Node) (*yaml.Node, error)
}
