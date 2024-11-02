package commandqueue

type CommandQueue struct {
	commands []string
}

func Make() CommandQueue {
	return CommandQueue{commands: make([]string, 0)}
}

func (q *CommandQueue) Add(command string) {
	q.commands = append(q.commands, command)
}

func (q *CommandQueue) None() bool {
	return len(q.commands) == 0
}

func (q *CommandQueue) Pop() string {
	if q.None() {
		return ""
	}

	command := q.commands[0]
	q.commands = q.commands[1:]
	return command
}
