// Copyright 2016 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"testing"
)

func clientServer(t *testing.T, options ...Option) (*Endpoint, *Endpoint, func()) {
	var wg sync.WaitGroup

	options = append(options, WithLogf(t.Logf))

	serverConn, clientConn := net.Pipe()

	server, err := NewEndpoint(serverConn, serverConn, serverConn, options...)
	if err != nil {
		t.Fatal(err)
	}

	client, err := NewEndpoint(clientConn, clientConn, clientConn, options...)
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	go func() {
		err := server.Serve()
		if err != nil && err != io.ErrClosedPipe {
			t.Logf("server: %v", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := client.Serve()
		if err != nil && err != io.ErrClosedPipe {
			t.Logf("server: %v", err)
		}
		wg.Done()
	}()

	cleanup := func() {
		client.Close()
		server.Close()
		wg.Wait()
	}

	return client, server, cleanup
}

func TestEndpoint(t *testing.T) {
	client, server, cleanup := clientServer(t)
	defer cleanup()

	if err := server.Register("add", func(a, b int) (int, error) { return a + b, nil }); err != nil {
		t.Fatal(err)
	}

	// Call.

	var sum int
	if err := client.Call("add", &sum, 1, 2); err != nil {
		t.Fatal(err)
	}

	if sum != 3 {
		t.Errorf("sum = %d, want %d", sum, 3)
	}

	// Notification.

	notifCh := make(chan string)
	if err := server.Register("n1", func(s string) { notifCh <- s }); err != nil {
		t.Fatal(err)
	}

	const n = 10
	for i := 0; i < i; i++ {
		for j := 0; j < n; j++ {
			if err := client.Notify("n1", fmt.Sprintf("notif %d,%d", i, j)); err != nil {
				t.Fatal(err)
			}
		}
		for j := 0; j < n; j++ {
			got := <-notifCh
			want := fmt.Sprintf("notif %d,%d", i, j)
			if got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		}
	}
}

var argsTests = []struct {
	sm     string
	args   []interface{}
	result []string
}{
	{"n", []interface{}{}, []string{"", ""}},
	{"n", []interface{}{"a"}, []string{"a", ""}},
	{"n", []interface{}{"a", "b"}, []string{"a", "b"}},
	{"n", []interface{}{"a", "b", "c"}, []string{"a", "b"}},

	{"v", []interface{}{}, []string{"", ""}},
	{"v", []interface{}{"a"}, []string{"a", ""}},
	{"v", []interface{}{"a", "b"}, []string{"a", "b"}},
	{"v", []interface{}{"a", "b", "x1"}, []string{"a", "b", "x1"}},
	{"v", []interface{}{"a", "b", "x1", "x2"}, []string{"a", "b", "x1", "x2"}},
	{"v", []interface{}{"a", "b", "x1", "x2", "x3"}, []string{"a", "b", "x1", "x2", "x3"}},

	{"a", []interface{}{}, []string(nil)},
	{"a", []interface{}{"x1", "x2", "x3"}, []string{"x1", "x2", "x3"}},
}

func TestArgs(t *testing.T) {
	client, server, cleanup := clientServer(t)
	defer cleanup()

	if err := server.Register("n", func(a, b string) ([]string, error) {
		return append([]string{a, b}), nil
	}); err != nil {
		t.Fatal(err)
	}

	if err := server.Register("v", func(a, b string, x ...string) ([]string, error) {
		return append([]string{a, b}, x...), nil
	}); err != nil {
		t.Fatal(err)
	}

	if err := server.Register("a", func(x ...string) ([]string, error) {
		return x, nil
	}); err != nil {
		t.Fatal(err)
	}

	for _, tt := range argsTests {
		var result []string
		if err := client.Call(tt.sm, &result, tt.args...); err != nil {
			t.Errorf("%s(%v) returned error %v", tt.sm, tt.args, err)
			continue
		}

		if !reflect.DeepEqual(result, tt.result) {
			t.Errorf("%s(%v) returned %#v, want %#v", tt.sm, tt.args, result, tt.result)
		}
	}
}

func TestCallAfterClose(t *testing.T) {
	client, server, cleanup := clientServer(t)
	err := server.Register("a", func() error {
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	cleanup()
	if err := client.Call("a", nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestExtraArgs(t *testing.T) {
	client, server, cleanup := clientServer(t)
	defer cleanup()

	err := server.Register("a", func(hello string) error {
		if hello != "hello" {
			t.Fatal("first arg not equal to 'hello'")
		}
		return nil
	}, "hello")
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Call("a", nil); err != nil {
		t.Fatal(err)
	}

	err = server.Register("b", func(hello *string) error {
		if hello != nil {
			t.Fatal("first arg not nil")
		}
		return nil
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Call("b", nil); err != nil {
		t.Fatal(err)
	}

}

func TestBadFunction(t *testing.T) {
	_, server, cleanup := clientServer(t)
	defer cleanup()

	err := server.Register("a", func(hello string) int {
		return 1
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
