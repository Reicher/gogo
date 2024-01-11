package api

type CommandType int

const (
	Greet CommandType = iota
	StartGame
	GetBoard
	MakeMove
	Pass
	Resign
)

var commandNames = map[CommandType]string{
	Greet:    "Greet",    // player name as data
	MakeMove: "MakeMove", // row, column, color as data
	Pass:     "Pass",     // no data
	Resign:   "Resign",   // no data
	GetBoard: "GetBoard", // no data
}

func (c CommandType) String() string {
	return commandNames[c]
}

func StringToCommandType(s string) (CommandType, bool) {
	for k, v := range commandNames {
		if v == s {
			return k, true
		}
	}
	return 0, false // or return an error
}

type Command struct {
	Type CommandType
	Data string
}
