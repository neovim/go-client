package rpc

import (
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"testing"
)

func testClientServer(tb testing.TB, opts ...Option) (client, server *Endpoint, cleanup func()) {
	tb.Helper()

	opts = append(opts, WithLogf(tb.Logf))

	serverConn, clientConn := net.Pipe()

	var err error
	server, err = NewEndpoint(serverConn, serverConn, serverConn, opts...)
	if err != nil {
		tb.Fatal(err)
	}

	client, err = NewEndpoint(clientConn, clientConn, clientConn, opts...)
	if err != nil {
		tb.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := server.Serve(); err != nil && !errors.Is(err, io.ErrClosedPipe) {
			tb.Errorf("server: %v", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if err := client.Serve(); err != nil && !errors.Is(err, io.ErrClosedPipe) {
			tb.Errorf("server: %v", err)
		}
		wg.Done()
	}()

	if tb.Failed() {
		tb.FailNow()
	}

	cleanup = func() {
		client.Close()
		server.Close()
		wg.Wait()
	}

	return client, server, cleanup
}

func TestEndpoint(t *testing.T) {
	t.Parallel()

	client, server, cleanup := testClientServer(t)
	defer cleanup()

	addFn := func(a, b int) (int, error) { return a + b, nil }
	if err := server.Register("add", addFn); err != nil {
		t.Fatal(err)
	}

	// call
	var sum int
	if err := client.Call("add", &sum, 1, 2); err != nil {
		t.Fatal(err)
	}
	if sum != 3 {
		t.Fatalf("sum = %d, want %d", sum, 3)
	}

	// notification
	notifCh := make(chan string)
	n1Fn := func(s string) { notifCh <- s }
	if err := server.Register("n1", n1Fn); err != nil {
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

func TestArgs(t *testing.T) {
	t.Parallel()

	client, server, cleanup := testClientServer(t)
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

	argsTests := []struct {
		sm     string
		args   []interface{}
		result []string
	}{
		{
			sm:     "n",
			args:   []interface{}{},
			result: []string{"", ""},
		},
		{
			sm:     "n",
			args:   []interface{}{"a"},
			result: []string{"a", ""},
		},
		{
			sm:     "n",
			args:   []interface{}{"a", "b"},
			result: []string{"a", "b"},
		},
		{
			sm:     "n",
			args:   []interface{}{"a", "b", "c"},
			result: []string{"a", "b"},
		},
		{
			sm:     "v",
			args:   []interface{}{},
			result: []string{"", ""},
		},
		{
			sm:     "v",
			args:   []interface{}{"a"},
			result: []string{"a", ""},
		},
		{
			sm:     "v",
			args:   []interface{}{"a", "b"},
			result: []string{"a", "b"},
		},
		{
			sm:     "v",
			args:   []interface{}{"a", "b", "x1"},
			result: []string{"a", "b", "x1"},
		},
		{
			sm:     "v",
			args:   []interface{}{"a", "b", "x1", "x2"},
			result: []string{"a", "b", "x1", "x2"},
		},
		{
			sm:     "v",
			args:   []interface{}{"a", "b", "x1", "x2", "x3"},
			result: []string{"a", "b", "x1", "x2", "x3"},
		},
		{
			sm:     "a",
			args:   []interface{}{},
			result: []string(nil),
		},
		{
			sm:     "a",
			args:   []interface{}{"x1", "x2", "x3"},
			result: []string{"x1", "x2", "x3"},
		},
	}
	for _, tt := range argsTests {
		t.Run(tt.sm, func(t *testing.T) {
			var result []string
			if err := client.Call(tt.sm, &result, tt.args...); err != nil {
				t.Fatalf("%s(%v) returned error %v", tt.sm, tt.args, err)
			}

			if !reflect.DeepEqual(result, tt.result) {
				t.Fatalf("%s(%v) returned %#v, want %#v", tt.sm, tt.args, result, tt.result)
			}
		})
	}
}

func TestCallAfterClose(t *testing.T) {
	t.Parallel()

	client, server, cleanup := testClientServer(t)

	if err := server.Register("a", func() error {
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	cleanup()

	if err := client.Call("a", nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestExtraArgs(t *testing.T) {
	t.Parallel()

	client, server, cleanup := testClientServer(t)
	defer cleanup()

	if err := server.Register("a", func(hello string) error {
		if hello != "hello" {
			t.Fatal("first arg not equal to 'hello'")
		}
		return nil
	}, "hello"); err != nil {
		t.Fatal(err)
	}

	if err := client.Call("a", nil); err != nil {
		t.Fatal(err)
	}

	if err := server.Register("b", func(hello *string) error {
		if hello != nil {
			t.Fatal("first arg not nil")
		}
		return nil
	}, nil); err != nil {
		t.Fatal(err)
	}

	if err := client.Call("b", nil); err != nil {
		t.Fatal(err)
	}
}

func TestBadFunction(t *testing.T) {
	t.Parallel()

	_, server, cleanup := testClientServer(t)
	defer cleanup()

	if err := server.Register("a", func(hello string) int {
		return 1
	}); err == nil {
		t.Fatal("expected error, got nil")
	}
}
