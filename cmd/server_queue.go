package main

import "github.com/goph/serverz"

// newServerQueue returns a new server queue with all the registered servers.
func newServerQueue(app *application) *serverz.Queue {
	queue := serverz.NewQueue()

	debugServer := newDebugServer(app)
	queue.Prepend(debugServer, nil)

	httpServer := newHTTPServer(app)
	queue.Append(httpServer, nil)

	return queue
}
