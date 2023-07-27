package command

import (
	"cfn-converter/converter"
	"flag"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"os"
)

var validate *validator.Validate

type Command struct {
	SourceFile string `validate:"required"`
	OutputFile string `validate:"required"`
	JoinToSub  bool
}

// New returns a new Command struct
func New() *Command {
	validate = validator.New()

	var (
		sourceFile string
		outputFile string
		joinToSub  bool
	)

	flag.StringVar(&sourceFile, "src", "", "CloudFormation template file")
	flag.StringVar(&outputFile, "out", "", "Output file")
	flag.BoolVar(&joinToSub, "join2sub", false, "Convert !Join to !Sub")
	flag.Parse()

	// Create a new command of the Command struct
	if outputFile == "" {
		outputFile = fmt.Sprintf("%s.converted.yaml", sourceFile)
	}
	command := &Command{
		SourceFile: sourceFile,
		OutputFile: outputFile,
		JoinToSub:  joinToSub,
	}

	// Validate struct
	if err := validate.Struct(command); err != nil {
		command.PrintUsage()
		os.Exit(1)
	}

	fmt.Printf("%+v\n", command)
	return command
}

// Run executes the command
func (c Command) Run() {
	println("Running command...")

	// load yaml file
	bytes, err := os.ReadFile(c.SourceFile)

	node := &yaml.Node{}
	// load yaml file
	_ = yaml.Unmarshal(bytes, node)

	// Convert
	var converters = make([]converter.Converter, 0)
	if c.JoinToSub {
		converters = append(converters, converter.JoinToSubConverter{})
	}
	for i, c := range converters {
		node, err = c.Convert(node)
		if err != nil {
			panic(fmt.Sprintf("Error converting [%d]", i))
		}
	}

	// Print yaml
	out, err := yaml.Marshal(node)
	if err != nil {
		panic(err)
	}

	// Write to file
	fmt.Println(string(out))
	err = os.WriteFile(c.OutputFile, out, 0644)
	if err != nil {
		panic(err)
	}
}

// PrintUsage prints the command usage
func (c Command) PrintUsage() {
	println("Usage: command <SourceFile>")
	println("Flags:")
	flag.PrintDefaults()
}
