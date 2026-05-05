package protocol

import (
	"errors"
	"strings"
)

var (
	ErrorEmptyCommand = errors.New("Emoty command")
)

type Command struct {
	Name string
	Args []string
}

func ParseCommand(input string) (Command, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return Command{}, ErrorEmptyCommand
	}

	splitCommand := strings.Fields(input)

	commandName := strings.ToUpper(splitCommand[0])
	args := splitCommand[1:]

	return Command{
		Name: commandName,
		Args: args,
	}, nil
}
