package command

import (
	"bytes"
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
	IndentSize int    `validate:"required"`
	JoinToSub  bool
}

// New returns a new Command struct
func New() *Command {
	validate = validator.New()

	var (
		sourceFile string
		outputFile string
		indentSize int
		joinToSub  bool
	)

	flag.StringVar(&sourceFile, "src", "", "CloudFormation template file")
	flag.StringVar(&outputFile, "out", "", "Output file")
	flag.BoolVar(&joinToSub, "join2sub", true, "Convert !Join to !Sub")
	flag.IntVar(&indentSize, "indent", 2, "Indent size")
	flag.Parse()

	// Create a new command of the Command struct
	if outputFile == "" {
		outputFile = fmt.Sprintf("%s.converted.yaml", sourceFile)
	}
	command := &Command{
		SourceFile: sourceFile,
		OutputFile: outputFile,
		IndentSize: indentSize,
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
	source, err := os.ReadFile(c.SourceFile)
	if err != nil {
		panic("Error reading file")
	}

	node := &yaml.Node{}
	// load yaml file
	_ = yaml.Unmarshal(source, node)

	// Convert
	var converters = make([]converter.Converter, 0)
	if c.JoinToSub {
		converters = append(converters, converter.JoinToSubConverter{})
	}
	for i, c := range converters {
		_, err = c.Convert(node)
		if err != nil {
			panic(fmt.Sprintf("Error converting [%d]", i))
		}
	}

	// Print yaml
	var out bytes.Buffer
	encoder := yaml.NewEncoder(&out)
	encoder.SetIndent(c.IndentSize)
	err = encoder.Encode(node)
	if err != nil {
		panic(err)
	}
	err = encoder.Close()
	if err != nil {
		return
	}

	// Write to file
	fmt.Println(out.String())
	err = os.WriteFile(c.OutputFile, out.Bytes(), 0644)
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
