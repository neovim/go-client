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

// Package nvim implements a Nvim client.
//
// See the ./plugin package for additional functionality required for writing
// Nvim plugins.
//
// The Nvim type implements the client. To connect to a running instance of
// Nvim, create a *Nvim value using the New or NewEmbedded functions and call
// the Serve() method to process RPC messages. Call the Close() method to
// release the resources used by the client.
//
// Use the Batch type to execute a sequence of Nvim API calls atomically. The
// Nvim NewBatch method creates new *Batch values.
package nvim
