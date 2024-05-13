// Package nvim implements a Nvim client.
//
// The Nvim type implements the client. To connect to a running instance of
// Nvim, create a *Nvim value using the New, Dial or NewChildProcess functions.
// Call the Close() method to release the resources used by the client.
//
// Use the Batch type to execute a sequence of Nvim API calls atomically. The
// Nvim NewBatch method creates new *Batch values.
package nvim
