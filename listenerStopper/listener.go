package listenerStopper

import (
	"errors"
	"net"
	"time"
)

type Listener struct {
	*net.TCPListener
	stop chan struct{}
}

func (l *Listener) Accept() (net.Conn, error) {
	for {
		select {
		case <-l.stop:
			return nil, errors.New("Listener signaled to stop")

		default:
			l.SetDeadline(time.Now().Add(time.Second))

			newConn, err := l.TCPListener.Accept()
			if err != nil {
				netErr, ok := err.(net.Error)
				if ok && netErr.Timeout() && netErr.Temporary() {
					continue
				}
			}
			return newConn, err
		}
	}
}

func (l *Listener) Stop() {
	l.stop <- struct{}{}
}

func NewListener(l *net.TCPListener) *Listener {
	return &Listener{
		TCPListener: l,
		stop:        make(chan struct{}),
	}
}
