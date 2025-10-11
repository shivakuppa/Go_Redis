package test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	conn, err := dialer.Dial("tcp", ":6379")
	require.NoError(t, err)

	defer conn.Close() //nolint:errcheck // OK for testing

	_, err = conn.Write([]byte("Hello from client"))
	require.NoError(t, err)

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	require.NoError(t, err)

	t.Log(string(buffer[:n]))
	t.Fail()
}