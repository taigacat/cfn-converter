package converter

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

type JoinToSubConverter struct{}

func (c JoinToSubConverter) Convert(data *yaml.Node) (*yaml.Node, error) {
	var err error = nil
	if data.Kind == yaml.MappingNode || data.Kind == yaml.DocumentNode {
		for i, node := range data.Content {
			if node.Tag == "!Join" {
				var stringArray = make([]string, 0)
				separator := node.Content[0].Value
				for _, item := range node.Content[1].Content {
					if item.Tag == "!Ref" {
						stringArray = append(stringArray, fmt.Sprintf("${%s}", item.Value))
					} else {
						stringArray = append(stringArray, item.Value)
					}
				}

				node = &yaml.Node{
					Kind:  yaml.ScalarNode,
					Tag:   "!Sub",
					Value: strings.Join(stringArray, separator),
				}
			}
			data.Content[i], err = c.Convert(node)
			if err != nil {
				panic("Error converting Join to Sub")
			}
		}
	} else if data.Kind == yaml.SequenceNode {
		for i, node := range data.Content {
			data.Content[i], err = c.Convert(node)
			if err != nil {
				panic("Error converting Join to Sub")
			}
		}
	}

	return data, err
}
