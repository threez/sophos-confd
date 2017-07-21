// Copyright 2016 Vincent Landgraf. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package confd

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func connHelper() *Conn {
	conn := NewAnonymousConn()
	conn.Logger = log.New(os.Stdout, "confd ", log.LstdFlags)
	conn.Options.Name = "confd-package-test"
	t := conn.Transport.(*tcpTransport)
	t.Timeout = time.Second * 1
	return conn
}

func systemConnHelper() *Conn {
	conn := connHelper()
	conn.Options.Username = "system"
	return conn
}

func TestInvalidURL(t *testing.T) {
	_, err := NewConn("%")
	assert.Error(t, err)
}

func TestConnFailed(t *testing.T) {
	conn, err := NewConn("http://127.0.0.1:50001")
	assert.NoError(t, err)
	_, err = conn.SimpleRequest("get_SID")
	assert.Error(t, err)
}

func TestInvalidCmd(t *testing.T) {
	conn := connHelper()
	_, err := conn.SimpleRequest("foobar")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "FATAL [FUNCTION_UNKNOWN] No public "+
		"function 'foobar' is provided by this Confd.")
}

func TestSafeURL(t *testing.T) {
	conn, err := NewConn("http://user:pass@127.0.0.1:5000/")
	assert.NoError(t, err)
	assert.Equal(t, "http://user:********@127.0.0.1:5000/", conn.safeURL())
}

func TestSID(t *testing.T) {
	conn := connHelper()
	defer func() { _ = conn.Close() }()
	assert.True(t, conn.Options.SID == nil)

	_, err := conn.SimpleRequest("get_SID")
	assert.NoError(t, err)
	assert.True(t, conn.Options.SID != nil)
	old := conn.Options.SID
	assert.NoError(t, conn.Close())
	_, err = conn.SimpleRequest("get_SID")
	assert.NoError(t, err)
	assert.True(t, conn.Options.SID == old)
}
