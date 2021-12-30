package nats

import (
	"log"

	"git.xenonstack.com/util/continuous-security-backend/config"
	"github.com/nats-io/nats.go"
)

// InitConnection is a function to initalize a nats connection with setup options
func InitConnection() {

	nc, err := nats.Connect(config.Conf.NatsServer.URL, setupOptions()...)
	if err != nil {
		log.Fatalln("nats error -> ", err)
		return
	}
	log.Println("nats connected successfully")
	config.NC = nc
	go Subscribe()
	go GitSubscribe()
}

//setupOptions is used to declear and create setup for NATS connection
func setupOptions() []nats.Option {
	opts := make([]nats.Option, 0)

	opts = append(opts, nats.Name("akirastack-website-scanner"))

	opts = append(opts, nats.UserInfo(config.Conf.NatsServer.Username, config.Conf.NatsServer.Password))
	opts = append(opts, nats.Token(config.Conf.NatsServer.Token))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		if !nc.IsClosed() {
			log.Println("Disconnected due to: ", err, ", will attempt reconnects in ", nats.DefaultReconnectWait.Seconds(), "s")
		}
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		if !nc.IsClosed() {
			log.Fatal("Exiting: no servers available")
		} else {
			log.Fatal("Exiting")
		}
	}))
	return opts
}
