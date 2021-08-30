package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/cppis/elio"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	defaultToSubscribe      = 500 * time.Millisecond
	defaultToPublish        = 500 * time.Millisecond
	defaultToquiesceMs uint = 1000
)

type sessionSet map[*elio.Session]struct{}

type Mqtter struct {
	mqtt       *elio.MqttClient
	url        string
	onConnect  mqtt.OnConnectHandler
	topics     map[string]sessionSet
}

func NewMqtter(url string, onConnect mqtt.OnConnectHandler) *Mqtter {
	q := new(Mqtter)
	if nil != q {
		q.url = url
		q.onConnect  = onConnect
		q.topics	 = make(map[string]sessionSet)
	}
	return q
}

// public methods
// String object to string
func (q *Mqtter) String() string {
	return fmt.Sprintf("Mqtter::%p", q)
}

// Mqtt mqtt
func (q *Mqtter) Mqtt() *elio.MqttClient {
	var err error
	if nil == q.mqtt {
		elio.AppInfo().Str(elio.LogObject, q.String()).
			Msgf("begin to connect to broker:%s with mqtt.client:%s", q.url, q.mqtt.String())

		q.mqtt = elio.NewMqttClient(q.url, q.onConnect)
		err = q.mqtt.Connect()
		if nil == err {
			elio.AppInfo().Str(elio.LogObject, q.String()).
				Msgf("succeed to connect to broker:%s with mqtt.client:%s", q.url, q.mqtt.String())
		}
	}

	if nil == err {
		return q.mqtt
	}

	elio.AppError().Str(elio.LogObject, q.String()).Err(err).
		Msgf("failed to connect to broker:%s with mqtt.client:%s", q.url, q.mqtt.String())

	return nil
}

// Pub pub
func (q *Mqtter) Pub(n *elio.Session, t string, p string) {
	err := q.Publish(t, p)
	if nil != err {
		//elio.AppError().Str(elio.LogObject, q.String()).Err(err).
		//	Msgf("failed to publish to mqtt.client:%s", q.mqtt.String())

	} else {
		elio.AppDebug().Str(elio.LogObject, q.String()).
			Msgf("succeed to publish to mqtt.client:%s with topic:%s payload:%s", q.mqtt.String(), t, p)

	}
}

// Sub sub
func (q *Mqtter) Sub(n *elio.Session, t string, c mqtt.MessageHandler) {
	err := q.Subscribe(t, c)
	if nil != err {
		//elio.AppError().Str(elio.LogObject, q.String()).Err(err).
		//	Msgf("failed to subscribe to mqtt.client:%s", q.mqtt.String())

	} else {
		elio.AppDebug().Str(elio.LogObject, q.String()).
			Msgf("succeed to subscribe to mqtt.client:%s with topic:%s", q.mqtt.String(), t)

		q.sessionSet[n] = struct{}{}
	}
}

// Unsub unsub
func (q *Mqtter) Unsub(n *elio.Session, t string) {
	err := q.Unsubscribe(t)
	if nil != err {
		//elio.AppError().Str(elio.LogObject, q.String()).Err(err).
		//	Msgf("failed to unsubscribe to mqtt.client:%s", q.mqtt.String())

	} else {
		//elio.AppTrace().Str(elio.LogObject, q.String()).
		//	Msgf("succeed to unsubscribe to mqtt.client:%s", q.mqtt.String())

		delete(q.sessionSet, n)
	}
}

func (q *Mqtter) Publish(t string, p string) (err error) {
	mqtt := q.Mqtt()
	err = mqtt.Publish(t, false, []byte(p))
	if nil != err {
		elio.AppError().Str(elio.LogObject, q.String()).Err(err).
			Msgf("failed to publish to %s", t)

	} else {
	}

	return err
}

// Subscribe subscribe from topic
func (q *Mqtter) Subscribe(topic string, callback mqtt.MessageHandler) error {
	mqtt := q.Mqtt()
	if nil != mqtt {
		return mqtt.Subscribe(topic, callback)
	}

	return fmt.Errorf("mqtt is not prepared")
}

// SubscribeMulti subscribe from topic
func (q *Mqtter) SubscribeMulti(topics map[string]byte, callback mqtt.MessageHandler) error {
	mqtt := q.Mqtt()
	if nil != mqtt {
		return mqtt.SubscribeMulti(topics, callback)
	}

	t := fmt.Sprintf("mqtt is not prepared")
	return errors.New(t)
}

// Unsubscribe unsubscribe from topic
func (q *Mqtter) Unsubscribe(topic string) error {
	mqtt := q.Mqtt()
	if nil != mqtt {
		return mqtt.Unsubscribe(topic)
	}

	return fmt.Errorf("mqtt is not prepared")
}

// private methods
// getMqtt get mqtt client
func (q *Mqtter) getMqtt() *elio.MqttClient {
	return q.mqtt
}

// setMqtt set mqtt client
func (q *Mqtter) setMqtt(c *elio.MqttClient) {
	q.mqtt = c
}

func (q *Mqtter) addSet(t string, n *elio.Session) {
	_, ok := q.topics[t]
	if false == ok {
		q.topics[t] = make(sessionSet)
	}

	q.topics[t][n] = struct{}{}
}

func (q *Mqtter) delSet(t string, n *elio.Session) (ok bool) {
	_, ok = q.topics[t]
	if true == ok {
		_, ok = q.topics[t][n]
		if true == ok {
			delete(q.topics[t], n)
		}
	}

	return ok
}

func (q *Mqtter) findSet(t string, n *elio.Session) (ok bool) {
	_, ok = q.topics[t]
	if true == ok {
		_, ok = q.topics[t][n]
	}

	return ok
}
