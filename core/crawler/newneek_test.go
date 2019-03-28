package crawler

import (
	"context"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const (
	testAddr = "127.0.0.1"
	testPort = "2470"
)

func TestNewNeek_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	http.Handle("/", http.FileServer(http.Dir("testdata")))
	testServer := &http.Server{Addr: ":" + testPort}
	go func() { testServer.ListenAndServe() }()
	defer testServer.Shutdown(ctx)

	crawler, err := NewNeekCrawler("http://"+testAddr+":"+testPort+"/newneek.html", nil)
	require.NoError(t, err)

	articles, err := crawler.Get()
	require.NoError(t, err)
	require.Equal(t, 6, len(articles))
}
