package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cppis/elio"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Herald service
type Herald struct {
	app     *elio.App
	ctx     context.Context
	cancel  context.CancelFunc
	prev    time.Time
	mqtter  *Mqtter
	mqttUrl string
}

func NewHerald(app *elio.App) *Herald {
	s := new(Herald)
	if nil != s {
		s.app = app
	}
	return s
}

// String string
func (s *Herald) String() string {
	return fmt.Sprintf("Herald::%p", s)
}

func (s *Herald) Name() string {
	return "herald"
}

func (s *Herald) OnInit(ctx context.Context, cancel context.CancelFunc) error {
	s.ctx = ctx
	s.cancel = cancel

	s.mqttUrl, _ = s.app.Config().GetStringOrDefault(fmt.Sprintf("%s.mqtt.url", s.Name()), "")
	//fmt.Printf("%s", s.mqttUrl)

	s.mqtter = NewMqtter(s.mqttUrl, s.mqttOnConnect)

	elio.AppDebug().Str(elio.LogObject, s.String()).Msg("on init")
	return nil
}

func (s *Herald) OnExit() {
	elio.AppDebug().Str(elio.LogObject, s.String()).Msg("on exit")
}

func (s *Herald) OnOpen(n *elio.Session) error {
	//fmt.Printf("o")

	elio.AppDebug().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, n.String()).Msgf("service:%s on.open", s.Name())

	return nil
}

func (s *Herald) OnClose(n *elio.Session, err error) {
	//fmt.Printf("c")

	s.mqtter.delAll(n)

	elio.AppDebug().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, n.String()).Msgf("service:%s on.close", s.Name())
}

func (s *Herald) OnError(n *elio.Session, err error) {
	//fmt.Printf("e")

	elio.AppError().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, n.String()).Msgf("service:%s on.error", s.Name())
}

const upperBoundRequestPacketLen uint16 = 4096

func (s *Herald) OnRead(n *elio.Session, in []byte) (processed int) {
	elio.AppTrace().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, n.String()).Msgf("service:%s on.read +%d", s.Name(), len(in))
	//fmt.Printf("+%d", len(in))

	n.GetKeepBuffer().Put(in)
	l := n.GetKeepBuffer().Len()

	ok := true

	for (true == ok) && (processed < l) {
		var parsed int
		var counts int

		i := n.GetKeepBuffer().Begin(processed)
		parsed, ok = s.onParse(upperBoundRequestPacketLen, i)
		if (true == ok) && (0 < parsed) {
			//fmt.Printf("+%d", len(i))

			request := elio.T2pParseCommand(i[:parsed])

			commands := elio.T2PParse(request)

			s.runCommand(commands, n)

			processed += parsed
			counts++

		} else {
			if 0 <= parsed {
				//AppDebug().Str(elio.LogObject, a.String()).Str(elio.LogSession, s.String()).
				//	Str(LogPayload, hex.Dump(i)).
				//	Msgf("fd:%d buffer.len:%v process.len:%v process.count:%d parse.ok:ok", s.GetFd(), l, processed, counts)

			} else {
				e := fmt.Sprintf("fd:%d invalid packet", n.GetFd())
				elio.AppError().Str(elio.LogObject, s.String()).Str(elio.LogSession, n.String()).
					Str(elio.LogPayload, hex.Dump(in)).Msgf(e)
				//err := errors.New(e)
				processed = elio.OnReadInvalidLen
			}

			ok = false
		}
	}

	n.GetKeepBuffer().Clear(processed)

	return processed
}

func (s *Herald) OnWrite(n *elio.Session, out []byte) {
	elio.AppTrace().Str(elio.LogObject, s.String()).
		Str(elio.LogSession, n.String()).Msgf("service:%s on.write -%d", s.Name(), len(out))
	//fmt.Printf("-%d", len(out))
}

const (
	// defaultFetchLimit default fetch limit
	defaultFetchLimit int = 2000
)

func (s *Herald) OnLoop(host *elio.IoHost, t time.Time, d time.Duration) {
	//fmt.Printf("on loop with delta:%v\n", d)
	//if t.Sub(e.prev) > 10*time.Second {
	//	fmt.Printf("e")
	//}

	//host.RunDivision(t, r.callbackDivision)

	_, _ = host.Dispatching(t, defaultFetchLimit)

	s.prev = t
}

func (s *Herald) onParse(lenMax uint16, in []byte) (int, bool) {
	parsed, delimiter, ok := elio.T2pOnParse(in)
	if false == ok {
		parsed = 0
		delimiter = 0
	}
	return parsed + delimiter, ok
}

const (
	// argLenEcho echo coomand argument length
	argLenEcho = 2
	// argLenPub pub coomand argument length
	argLenPub = 3
	// argLenSub sub coomand argument length
	argLenSub = 2
	// argLenUnsub unsub coomand argument length
	argLenUnsub = 2
)

// runCommand run command
func (s *Herald) runCommand(commands []string, n *elio.Session) {
	if len(commands) <= 0 {
		return
	}

	l := len(commands)
	var err error

	switch commands[0] {
	case "exit":
		fallthrough
	case "quit":
		fmt.Printf("bye~\n")
		elio.Elio().End()
	case "echo":
		if l < argLenEcho {
			err = fmt.Errorf("invalid arg len:%d with expect:%d", l, argLenEcho)
		} else {
			//commands[1]:	echo message
			s.echo(n, commands[1])
		}
	case "pub":
		if l < argLenPub {
			err = fmt.Errorf("invalid arg len:%d with expect:%d", l, argLenPub)
		} else {
			//commands[1]:	topic
			//commands[2]:	message
			s.mqtter.Pub(n, commands[1], commands[2])
		}
	case "sub":
		if l < argLenSub {
			err = fmt.Errorf("invalid arg len:%d with expect:%d", l, argLenSub)
		} else {
			//commands[1]:	topic
			s.mqtter.Sub(n, commands[1], s.mqttOnSub)
		}
	case "unsub":
		if l < argLenUnsub {
			err = fmt.Errorf("invalid arg len:%d with expect:%d", l, argLenUnsub)
		} else {
			//commands[1]:	topic
			s.mqtter.Unsub(n, commands[1])
		}
	case "list":
		// list topics
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

// echo echo
func (s *Herald) echo(n *elio.Session, m string) {
	n.Write([]byte(m))
}

// callback methods
// mqttOnConnect mqtt on connect
func (s *Herald) mqttOnConnect(c mqtt.Client) {
	elio.AppInfo().Str(elio.LogObject, s.String()).Msg("mqtt on.connect")

	// topics := map[string]byte{
	// 	TopicNotify():                           1,
	// 	TopicEcho(elio.GetAppParams().GetSuid()): 1,
	// }
	// if err := r.echo.SubscribeMulti(topics, processSubscribe); nil != err {
	// 	heatgo.AlogError().Str(elio.LogObject, r.echo.String()).Err(err).Msg("failed to subscribe.multi")
	// } else {
	// 	heatgo.AlogInfo().Str(elio.LogObject, r.echo.String()).Msg("succeed to to subscribe.multi")
	// }
}

// mqttOnSub mqtt on subscribe
func (s *Herald) mqttOnSub(c mqtt.Client, m mqtt.Message) {
	elio.AppInfo().Str(elio.LogObject, s.String()).
		Msgf("mqtt on.sub with topic:%s payload:%s", m.Topic(), string(m.Payload()))

	s.mqtter.OnSub(m.Topic(), string(m.Payload()))
}
