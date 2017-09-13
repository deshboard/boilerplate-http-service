package main

import "github.com/goph/serverz"

// newServerQueue returns a new server queue with all the registered servers.
func newServerQueue(a *application) *serverz.Queue {
	queue := serverz.NewQueue()

	debugServer := newDebugServer(a)
	queue.Prepend(debugServer, nil)

	httpServer := newHTTPServer(appCtx)
	queue.Append(httpServer, nil)

	return queue
}
