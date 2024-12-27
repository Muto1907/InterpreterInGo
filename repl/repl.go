package repl

import (
	"fmt"
	"io"

	"github.com/chzyer/readline"

	"github.com/Muto1907/interpreterInGo/evaluator"
	"github.com/Muto1907/interpreterInGo/lexer"
	"github.com/Muto1907/interpreterInGo/object"
	"github.com/Muto1907/interpreterInGo/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {

	cfg := &readline.Config{
		Prompt:                 PROMPT,
		HistoryFile:            "tmp/chimp_repl_history",
		InterruptPrompt:        "^C",
		EOFPrompt:              "exit",
		DisableAutoSaveHistory: false,
		EnableMask:             false,
	}

	rl, err := readline.NewEx(cfg)
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	env := object.NewEnvironment()
	eval := evaluator.NewEval()

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				return
			} else {
				continue
			}
		} else if err == io.EOF {
			return
		} else if err != nil {
			fmt.Fprintln(out, "Error reading line:", err)
			return
		}

		if line == "" {
			continue
		}
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}
		evaluated := eval.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect()+"\n")
		}
		//io.WriteString(out, fmt.Sprintf("Heap Size after eval: %d\n", len(eval.Heap)))

	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
