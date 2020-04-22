package sig

import (
	"errors"
	"fmt"
	"github.com/lithdew/bytesutil"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	RegexCipherFactory = regexp.MustCompile(`(?s)var (\w+)={(\w+:function\([\w,]*?\){.*?},?)+}`)
	RegexCipherSteps   = regexp.MustCompile(`(?s)\w+?=function\(\w+?\){\w+?=\w+?\.split\(""\);((?:\w+?\.\w+?\(\w+?(?:,\d+?)?\);)*)return \w+?.join\(""\)};`)
	RegexCipherMethod  = regexp.MustCompile(`(?s)\w+?\.(\w+?)\(\w+?(?:,(\d+?))?\)`)

	RegexSliceOp   = regexp.MustCompile(`(?s)function\(\w+,\w+\){\w+\.splice\(0,\w+\)}`)
	RegexReverseOp = regexp.MustCompile(`(?s)function\(\w+\){\w+\.reverse\(\)}`)
	RegexSwapOp    = regexp.MustCompile(`(?s)function\(\w+,\w+\){var \w+=\w+\[0];\w+\[0]=\w+\[\w+%\w+\.length];\w+\[\w+%\w+\.length]=\w+}`)
)

type Cipher []Step

func (cipher Cipher) Decode(s string) string {
	sig := bytesutil.Slice(s)

	for _, step := range cipher {
		sig = step(sig)
	}

	return string(sig)
}

type Step func(s []byte) []byte

type StepType uint8

const (
	SliceOp StepType = iota
	ReverseOp
	SwapOp
)

func (c StepType) String() string {
	switch c {
	case SliceOp:
		return "slice"
	case ReverseOp:
		return "reverse"
	case SwapOp:
		return "swap"
	default:
		panic(fmt.Sprintf("unknown cipher step type: %d", c))
	}
}

func (c StepType) Instruction(param int) Step {
	switch c {
	case SliceOp:
		return func(s []byte) []byte {
			return s[param:]
		}
	case ReverseOp:
		return func(s []byte) []byte {
			for i := len(s)/2 - 1; i >= 0; i-- {
				opp := len(s) - 1 - i
				s[i], s[opp] = s[opp], s[i]
			}
			return s
		}
	case SwapOp:
		return func(s []byte) []byte {
			s[0], s[param%len(s)] = s[param%len(s)], s[0]
			return s
		}
	default:
		panic(fmt.Sprintf("unknown cipher step type: %d", c))
	}
}

func LookupCipher(f CipherFactory, script string) (Cipher, error) {
	matches := RegexCipherSteps.FindStringSubmatch(script)
	if matches == nil {
		return nil, errors.New("could not find cipher steps")
	}

	entries := strings.FieldsFunc(matches[1], func(r rune) bool { return r == ';' })
	steps := make(Cipher, 0, len(entries))

	for _, entry := range entries {
		args := RegexCipherMethod.FindStringSubmatch(entry)
		if args == nil {
			return nil, fmt.Errorf("failed to parse arguments of cipher method call %q", entry)
		}

		step, exists := f.methods[args[1]]
		if !exists {
			return nil, fmt.Errorf("script calls cipher method %q which is not registered in cipher factory", args[1])
		}

		param := 0

		if args[2] != "" {
			p, err := strconv.ParseInt(args[2], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("failed to decode argument of cipher method call %q: %w", entry, err)
			}

			param = int(p)
		}

		steps = append(steps, step.Instruction(param))
	}

	return steps, nil
}

type CipherFactory struct {
	id      string
	methods map[string]StepType
}

func LookupCipherFactory(script string) (CipherFactory, error) {
	var f CipherFactory

	matches := RegexCipherFactory.FindStringSubmatch(script)
	if matches == nil {
		return f, errors.New("could not find cipher factory")
	}

	f.id = matches[1]
	f.methods = make(map[string]StepType)

	buf := matches[2]

	for len(buf) > 0 {
		i := 0
		for unicode.IsSpace(rune(buf[i])) {
			i++
		}
		if i >= len(buf) {
			return f, errors.New("failed to remove whitespace while decoding cipher method name")
		}

		buf = buf[i:]

		i = strings.IndexByte(buf, ':')
		if i == -1 {
			return f, errors.New("failed to find end of cipher method name")
		}

		method := buf[:i]
		buf = buf[i+1:]

		i = strings.IndexByte(buf, '\n')
		if i == -1 {
			i = len(buf) - 1
		}

		body := buf[:i+1]

		switch {
		case RegexSliceOp.MatchString(body):
			f.methods[method] = SliceOp
		case RegexReverseOp.MatchString(body):
			f.methods[method] = ReverseOp
		case RegexSwapOp.MatchString(body):
			f.methods[method] = SwapOp
		default:
			return f, fmt.Errorf("do not recognize cipher method %q with body %q", method, body)
		}

		buf = buf[i+1:]
	}

	return f, nil
}
