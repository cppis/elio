package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cppis/elio"
)

// Sample sample service
type Sample struct {
	app    *elio.App
	ctx    context.Context
	cancel context.CancelFunc
	prev   time.Time
}

func NewSample(app *elio.App) *Sample {
	e := new(Sample)
	if nil != e {
		e.app = app
	}
	return e
}

// String string
func (s *Sample) String() string {
	return fmt.Sprintf("Sample::%p", s)
}

func (s *Sample) Name() string {
	return "sample"
}

func (s *Sample) OnInit(ctx context.Context, cancel context.CancelFunc) error {
	s.ctx = ctx
	s.cancel = cancel

	fmt.Printf("%s on init\n", s.Name())
	return nil
}

func (s *Sample) OnExit() {
	fmt.Printf("%s on exit\n", s.Name())
}

func (s *Sample) OnOpen(n *elio.Session) error {
	//fmt.Printf("%s on open\n", e.Name())
	n.GetIo().Host.PostToQueue(&EventOpen{
		session: n,
		host:    n.GetIo().Host,
	})

	fmt.Printf("o")

	elio.AppDebug().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, n.String()).Msgf("service:%s on.open", s.Name())

	return nil
}

func (s *Sample) OnClose(n *elio.Session, err error) {
	//fmt.Printf("%s on close\n", e.Name())

	n.GetIo().Host.PostToQueue(&EventClose{
		session: n,
		host:    n.GetIo().Host,
	})

	fmt.Printf("c")

	elio.AppDebug().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, n.String()).Msgf("service:%s on.close", s.Name())
}

func (s *Sample) OnError(n *elio.Session, err error) {
	//fmt.Printf("on error\n")
	elio.AppError().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, n.String()).Msgf("service:%s on.error", s.Name())
}

const (
	defaultPacketMaxLen = 4096 * 8
)

func (s *Sample) OnRead(n *elio.Session, in []byte) (processed int) {
	//fmt.Printf("on read %d\n", len(in))
	//e.app.End()
	//s.GetIo().End()
	//return len(in)
	l := len(in)

	ok := true

	for (true == ok) && (processed < l) {
		var parsed int
		var counts int

		parsed, ok = s.onParse(defaultPacketMaxLen, in[processed:])
		if (true == ok) && (0 < parsed) {
			fmt.Printf("+%d", len(in))

			request := elio.T2pParseCommand(in)

			commands := elio.T2PParse(request)

			s.runCommand(commands, n)

			processed += parsed
			counts++

		} else {
			if 0 <= parsed {
				//AlogDebug().Str(elf.LogObject, a.String()).Str(elf.LogSession, s.String()).
				//	Str(LogPayload, hex.Dump(i)).
				//	Msgf("fd:%d buffer.len:%v process.len:%v process.count:%d parse.ok:ok", s.GetFd(), l, processed, counts)

			} else {
				m := fmt.Sprintf("fd:%d invalid packet size", n.GetFd())
				elio.AppError().Str(elio.LogObject, s.String()).
					Str(elio.LogSession, n.String()).
					Str(elio.LogPayload, hex.Dump(in)).Msgf(m)
				//err := errors.New(e)
			}

			ok = false
		}
	}

	return processed
}

func (s *Sample) OnWrite(n *elio.Session, out []byte) {
	fmt.Printf("on write\n")
}

const (
	// defaultFetchLimit default fetch limit
	defaultFetchLimit int = 2000
)

func (s *Sample) OnLoop(host *elio.IoHost, t time.Time, d time.Duration) {
	//fmt.Printf("on loop with delta:%v\n", d)
	if t.Sub(s.prev) > 3*time.Second {
		fmt.Printf("e")
	}

	//host.RunDivision(t, r.callbackDivision)

	_, _ = host.Dispatching(t, defaultFetchLimit)

	s.prev = t
}

func (s *Sample) onParse(lenMax uint16, in []byte) (int, bool) {
	parsed, delimiter, ok := elio.T2pOnParse(in)
	if false == ok {
		parsed = 0
		delimiter = 0
	}
	return parsed + delimiter, ok
}

const (
	// argumentLenQuery query coomand argument length
	argumentLenQuery = 2
	// argumentLenShow show coomand argument length
	argumentLenShow = 2
	// argumentLenEcho echo coomand argument length
	argumentLenEcho = 2
)

// runCommand run command
func (s *Sample) runCommand(commands []string, n *elio.Session) {
	if len(commands) <= 0 {
		return
	}

	elio.Elio().MetricMeter(fmt.Sprintf(elio.MetricAppIoInCountf, s.Name())).Mark(1)

	l := len(commands)
	var err error

	switch commands[0] {
	case "exit":
		fallthrough
	case "quit":
		fmt.Printf("bye~\n")
		elio.Elio().End()
	// case "show":
	// 	//show := server.Show()
	// 	//show += "\n\n"
	// 	//s.Write([]byte(show))
	// 	if l < argumentLenShow {
	// 		err = fmt.Errorf("invalid arg len:%d, expect:%d", l, argumentLenShow)
	// 	} else {
	// 		r.Show(s, commands[1])
	// 	}
	case "query":
		if l < argumentLenQuery {
			err = fmt.Errorf("invalid arg len:%d, expect:%d", l, argumentLenQuery)
		} else {
			s.Query(n, commands[1])
		}
	case "echo":
		if l < argumentLenEcho {
			err = fmt.Errorf("invalid arg len:%d, expect:%d", l, argumentLenEcho)
		} else {
			//commands[1]:	echo message
			s.Echo(n, commands[1])
		}
	default:
		err = fmt.Errorf("invalid command:%s", commands[0])
		break
	}

	if nil != err {
		//c := s.GetContext().(*house.AppContext)
		out := fmt.Sprintf("%s\nresult: %s\n\n", commands[0], err.Error())
		n.Write([]byte(out))
	}
}
