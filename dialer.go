package simplelb

import (
	"context"
	"net"
	"time"
)

// DialContextFn was defined to make code more readable.
// https://gist.github.com/c4milo/275abc6eccbfd88ad56ca7c77947883a
type DialContextFn func(ctx context.Context, network, address string) (net.Conn, error)

// TimeoutDialContext implements our own dialer in order to set read and write idle timeouts.
func TimeoutDialContext(rwtimeout, ctimeout time.Duration) DialContextFn {
	dialer := &net.Dialer{Timeout: ctimeout}

	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		c, err := dialer.DialContext(ctx, network, addr)
		if err != nil {
			return nil, err
		}

		if rwtimeout > 0 {
			c = &tcpConn{TCPConn: c.(*net.TCPConn), timeout: rwtimeout}
		}

		return c, nil
	}
}

// tcpConn is our own net.Conn which sets a read and write deadline and resets them each
// time there is read or write activity in the connection.
type tcpConn struct {
	*net.TCPConn
	timeout time.Duration
}

// Read implements the Conn Read method.
func (c *tcpConn) Read(b []byte) (int, error) {
	err := c.TCPConn.SetDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return 0, err
	}

	return c.TCPConn.Read(b)
}

// Write implements the Conn Write method.
func (c *tcpConn) Write(b []byte) (int, error) {
	err := c.TCPConn.SetDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return 0, err
	}

	return c.TCPConn.Write(b)
}
