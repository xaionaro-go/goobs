package goobs_test

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	goobs "github.com/andreykaipov/goobs"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func Test_client(t *testing.T) {
	t.Run("wrong password", func(t *testing.T) {
		_, err := goobs.New(
			"localhost:"+os.Getenv("OBS_PORT"),
			goobs.WithPassword("wrongpassword"),
			goobs.WithRequestHeader(http.Header{"User-Agent": []string{"goobs-e2e/0.0.0"}}),
		)
		assert.Error(t, err)
		assert.IsType(t, &websocket.CloseError{}, err)
		assert.Equal(t, err.(*websocket.CloseError).Code, 4009)
	})

	t.Run("server isn't running", func(t *testing.T) {
		_, err := goobs.New(
			"localhost:42069",
			goobs.WithPassword("wrongpassword"),
			goobs.WithRequestHeader(http.Header{"User-Agent": []string{"goobs-e2e/0.0.0"}}),
		)
		assert.Error(t, err)
		assert.IsType(t, &net.OpError{}, err)
	})

	t.Run("right password", func(t *testing.T) {
		client, err := goobs.New(
			"localhost:"+os.Getenv("OBS_PORT"),
			goobs.WithPassword("goodpassword"),
			goobs.WithRequestHeader(http.Header{"User-Agent": []string{"goobs-e2e/0.0.0"}}),
		)
		assert.NoError(t, err)
		t.Cleanup(func() {
			client.Disconnect()
		})
		time.Sleep(1 * time.Second)
	})
}

func Test_multi_goroutine(t *testing.T) {
	for i := 1; i <= 10; i++ {
		t.Run(fmt.Sprintf("goroutine-%d", i), func(t *testing.T) {
			t.Parallel()

			client, err := goobs.New(
				"localhost:"+os.Getenv("OBS_PORT"),
				goobs.WithPassword("goodpassword"),
				goobs.WithRequestHeader(http.Header{"User-Agent": []string{"goobs-e2e/0.0.0"}}),
			)
			assert.NoError(t, err)
			t.Cleanup(func() {
				client.Disconnect()
			})
			wg := sync.WaitGroup{}
			for i := 0; i < 5_000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, err := client.Scenes.GetSceneList()
					assert.NoError(t, err)
				}()
			}
			wg.Wait()
		})
	}
}

func Test_profile(t *testing.T) {
	client, err := goobs.New(
		"localhost:"+os.Getenv("OBS_PORT"),
		goobs.WithPassword("goodpassword"),
		goobs.WithRequestHeader(http.Header{"User-Agent": []string{"goobs-e2e/0.0.0"}}),
	)
	assert.NoError(t, err)
	t.Cleanup(func() {
		client.Disconnect()
	})

	for n := 0; n < 5_000; n++ {
		_, err := client.Scenes.GetSceneList()
		assert.NoError(t, err)
	}
}