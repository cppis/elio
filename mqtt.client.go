package elio

import (
	"encoding/hex"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	defaultToSubscribe      = 500 * time.Millisecond
	defaultToPublish        = 500 * time.Millisecond
	defaultToquiesceMs uint = 1000
)

// MqttClient mqtt client
type MqttClient struct {
	client      mqtt.Client
	qos         byte
	toSubscribe time.Duration
	toPublish   time.Duration
	toquiesce   uint
	broker      string
}

// NewMqttClient new mqtt client
func NewMqttClient(broker string, onConnect mqtt.OnConnectHandler) (m *MqttClient) {
	if m = new(MqttClient); nil != m {
		opts := mqtt.NewClientOptions()
		m.broker = broker
		opts.AddBroker(broker)
		//opts.SetWriteTimeout(1 * time.Second)
		//opts.SetClientID(*id)
		//opts.SetUsername(*user)
		//opts.SetPassword(*password)
		//opts.SetCleanSession(*cleansess)	// default false
		//opts.SetAutoReconnect(true)
		if nil != onConnect {
			opts.SetOnConnectHandler(onConnect)
		}
		//if nil != onLost {
		//	opts.SetConnectionLostHandler(onLost)
		//}
		//if DebugEnabled() {
		// 	mqtt.DEBUG = log.New(os.Stdout, "DEBUG ", 0)
		// 	mqtt.WARN = log.New(os.Stdout, "WARN ", 0)
		// 	mqtt.CRITICAL = log.New(os.Stdout, "CRITICAL ", 0)
		// 	mqtt.ERROR = log.New(os.Stdout, "ERROR ", 0)
		// }

		m.client = mqtt.NewClient(opts)
		m.qos = 1
		m.toSubscribe = defaultToSubscribe
		m.toPublish = defaultToPublish
		m.toquiesce = defaultToquiesceMs
	}

	return m
}

// func onLost(c mqtt.Client, err error) {
// 	c.Disconnect(100)
// }

// String object to string
func (m *MqttClient) String() string {
	return fmt.Sprintf("MqttClient::%p", m)
}

// Connect connect
func (m *MqttClient) Connect() (err error) {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		err = fmt.Errorf("connect failed with timeout:%v", m.toPublish)
		AppError().Str(LogObject, m.String()).Err(token.Error()).
			Msgf("failed to connect to %s", m.broker)
	}

	return err
}

// Close close
func (m *MqttClient) Close() {
	m.client.Disconnect(m.toquiesce)
}

// SafeClose safe close
func (m *MqttClient) SafeClose() {
	m.client.Disconnect(m.toquiesce)
	m.client = nil
}

// IsConnected is connected
func (m *MqttClient) IsConnected() bool {
	return m.client.IsConnected()
}

// Publish publish
//	token.Wait(): Can also use '<-t.Done()' in releases > 1.2.0
func (m *MqttClient) Publish(topic string, retained bool, payload []byte) (err error) {
	token := m.client.Publish(topic, m.qos, retained, payload)
	go func() {
		if false == token.WaitTimeout(m.toPublish) {
			err = fmt.Errorf("failed to publish:%s retain:%v payload.len:%d", topic, retained, len(payload))
		} else {
			err = token.Error()
		}

		//_ = token.Wait()
		//err = token.Error()
		if nil != err {
			AppError().Str(LogObject, m.String()).Err(err).
				Msgf("failed to publish:%s with len:%d", topic, len(payload))
		} else {
			AppTrace().Str(LogObject, m.String()).Str(LogPayload, hex.Dump(payload)).
				Msgf("succeed to publish:%s with len:%d", topic, len(payload))
		}
	}()

	return nil
}

// Subscribe subscribe
func (m *MqttClient) Subscribe(topic string, callback mqtt.MessageHandler) (err error) {
	token := m.client.Subscribe(topic, m.qos, callback)
	_ = token.Wait()
	err = token.Error()
	if nil != err {
		AppError().Str(LogObject, m.String()).Err(err).
			Msgf("failed to subscribe:%s", topic)
	} else {
		AppTrace().Str(LogObject, m.String()).
			Msgf("succeed to subscribe:%s", topic)
	}

	return err
}

// SubscribeMulti subscribe multi
func (m *MqttClient) SubscribeMulti(topics map[string]byte, callback mqtt.MessageHandler) (err error) {
	token := m.client.SubscribeMultiple(topics, callback)
	_ = token.Wait()
	err = token.Error()
	if nil != err {
		AppError().Str(LogObject, m.String()).Err(err).
			Msgf("failed to subscribemulti:%v", topics)
	} else {
		AppTrace().Str(LogObject, m.String()).
			Msgf("succeed to subscribemulti:%v", topics)
	}

	return err
}

// Unsubscribe unsubscribe
func (m *MqttClient) Unsubscribe(topics ...string) (err error) {
	if false == m.IsConnected() {
		return nil
	}

	token := m.client.Unsubscribe(topics...)
	_ = token.Wait()
	err = token.Error()
	if nil != err {
		AppError().Str(LogObject, m.String()).Err(err).
			Msgf("failed to unsubscribe:%s", topics)
	} else {
		AppTrace().Str(LogObject, m.String()).
			Msgf("succeed to unsubscribe:%s", topics)
	}

	return err
}

// PublishNoWait publish no wait
func (m *MqttClient) PublishNoWait(topic string, retained bool, payload []byte) {
	token := m.client.Publish(topic, m.qos, retained, payload)
	go func() {
		var err error
		if false == token.WaitTimeout(m.toPublish) {
			err = fmt.Errorf("publish:%s failed with timeout:%v", topic, m.toPublish)
		} else {
			err = token.Error()
		}

		if nil != err {
			AppError().Str(LogObject, m.String()).Err(token.Error()).
				Msgf("failed to publish:%s with len:%d", topic, len(payload))

			m.Close()

		} else {
			AppTrace().Str(LogObject, m.String()).Str(LogPayload, hex.Dump(payload)).
				Msgf("succeed to publish:%s with len:%d", topic, len(payload))
		}
	}()
}

// SubscribeNoWait subscribe async no wait
func (m *MqttClient) SubscribeNoWait(topic string, callback mqtt.MessageHandler) {
	token := m.client.Subscribe(topic, m.qos, callback)
	go func() {
		var err error
		if false == token.WaitTimeout(m.toSubscribe) {
			err = fmt.Errorf("publish:%s failed with timeout:%v", topic, m.toPublish)
		} else {
			err = token.Error()
		}

		if nil != err {
			//fmt.Printf("subscribe:%s failed with error:%s\n", topic, token.Error().Error())
			AppError().Str(LogObject, m.String()).Err(token.Error()).
				Msgf("failed to subscribe:%s", topic)

			m.Close()
		}
	}()
}

// SubscribeMultiNoWait subscribe multi no wait
func (m *MqttClient) SubscribeMultiNoWait(topics map[string]byte, callback mqtt.MessageHandler) {
	token := m.client.SubscribeMultiple(topics, callback)
	go func() {
		var err error
		if false == token.WaitTimeout(m.toSubscribe) {
			err = fmt.Errorf("publish:%v failed with timeout:%v", topics, m.toPublish)
		} else {
			err = token.Error()
		}

		if nil != err {
			//fmt.Printf("subscribemulti:%s failed with error:%s\n", topic, token.Error().Error())
			AppError().Str(LogObject, m.String()).Err(err).
				Msgf("failed to subscribemulti:%v", topics)

			m.Close()
		}
	}()
}

// UnsubscribeNoWait unsubscribe no wait
func (m *MqttClient) UnsubscribeNoWait(topics ...string) {
	if false == m.IsConnected() {
		AppDebug().Str(LogObject, m.String()).
			Msgf("cannot unsubscribe:%v invalid connection", topics)
		return
	}

	token := m.client.Unsubscribe(topics...)
	go func() {
		var err error
		if false == token.WaitTimeout(m.toSubscribe) {
			err = fmt.Errorf("unsubscribe:%v failed with timeout:%v", topics, m.toSubscribe)

		} else {
			err = token.Error()
		}

		if nil != err {
			//fmt.Printf("unsubscribe:%s failed with error:%v\n", topic, err)
			AppError().Str(LogObject, m.String()).Err(err).
				Msgf("unsubscribe:%v failed", topics)
			m.Close()
		}
	}()
}
