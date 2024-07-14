package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

// ServeHTTP handles the logic of running  server in a goroutine and waiting for signal to gracefully stop the server
// on ctx.Done signal a request to shut down the server is sent, so that no new requests will be served.
func ServeHTTP(ctx context.Context, appPort string, handler http.Handler) {
	if appPort[0] != ':' {
		appPort = ":" + appPort
	}

	srv := &http.Server{Addr: appPort, Handler: handler}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen:%s\n", err)
		}
	}()

	log.Printf("server started on port %s", appPort)

	<-ctx.Done()

	log.Printf("graceful shutdown request received")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%s", err.Error())
	}

	log.Println("application stopped accepting requests")
}
