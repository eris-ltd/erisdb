// Copyright 2017 Monax Industries Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loggers

import (
	"github.com/monax/eris-db/logging/structure"
	kitlog "github.com/go-kit/kit/log"
)

const (
	InfoChannelName  = "Info"
	TraceChannelName = "Trace"

	InfoLevelName  = InfoChannelName
	TraceLevelName = TraceChannelName
)

type infoTraceLogger struct {
	infoLogger  *kitlog.Context
	traceLogger *kitlog.Context
}

// InfoTraceLogger maintains two independent concurrently-safe channels of
// logging. The idea behind the independence is that you can ignore one channel
// with no performance penalty. For more fine grained filtering or aggregation
// the Info and Trace loggers can be decorated loggers that perform arbitrary
// filtering/routing/aggregation on log messages.
type InfoTraceLogger interface {
	// Send a log message to the default channel
	kitlog.Logger

	// Send an log message to the Info channel, formed of a sequence of key value
	// pairs. Info messages should be operationally interesting to a human who is
	// monitoring the logs. But not necessarily a human who is trying to
	// understand or debug the system. Any handled errors or warnings should be
	// sent to the Info channel (where you may wish to tag them with a suitable
	// key-value pair to categorise them as such).
	Info(keyvals ...interface{}) error

	// Send an log message to the Trace channel, formed of a sequence of key-value
	// pairs. Trace messages can be used for any state change in the system that
	// may be of interest to a machine consumer or a human who is trying to debug
	// the system or trying to understand the system in detail. If the messages
	// are very point-like and contain little structure, consider using a metric
	// instead.
	Trace(keyvals ...interface{}) error

	// A logging context (see go-kit log's Context). Takes a sequence key values
	// via With or WithPrefix and ensures future calls to log will have those
	// contextual values appended to the call to an underlying logger.
	// Values can be dynamic by passing an instance of the kitlog.Valuer interface
	// This provides an interface version of the kitlog.Context struct to be used
	// For implementations that wrap a kitlog.Context. In addition it makes no
	// assumption about the name or signature of the logging method(s).
	// See InfoTraceLogger

	// Establish a context by appending contextual key-values to any existing
	// contextual values
	With(keyvals ...interface{}) InfoTraceLogger

	// Establish a context by prepending contextual key-values to any existing
	// contextual values
	WithPrefix(keyvals ...interface{}) InfoTraceLogger
}

// Interface assertions
var _ InfoTraceLogger = (*infoTraceLogger)(nil)
var _ kitlog.Logger = (InfoTraceLogger)(nil)

func NewInfoTraceLogger(infoLogger, traceLogger kitlog.Logger) InfoTraceLogger {
	// We will never halt the progress of a log emitter. If log output takes too
	// long will start dropping log lines by using a ring buffer.
	// We also guard against any concurrency bugs in underlying loggers by feeding
	// them from a single channel
	logger := kitlog.NewContext(NonBlockingLogger(VectorValuedLogger(
		MultipleChannelLogger(
			map[string]kitlog.Logger{
				InfoChannelName:  infoLogger,
				TraceChannelName: traceLogger,
			}))))
	return &infoTraceLogger{
		infoLogger: logger.With(
			structure.ChannelKey, InfoChannelName,
			structure.LevelKey, InfoLevelName,
		),
		traceLogger: logger.With(
			structure.ChannelKey, TraceChannelName,
			structure.LevelKey, TraceLevelName,
		),
	}
}

func NewNoopInfoTraceLogger() InfoTraceLogger {
	noopLogger := kitlog.NewNopLogger()
	return NewInfoTraceLogger(noopLogger, noopLogger)
}

func (l *infoTraceLogger) With(keyvals ...interface{}) InfoTraceLogger {
	return &infoTraceLogger{
		infoLogger:  l.infoLogger.With(keyvals...),
		traceLogger: l.traceLogger.With(keyvals...),
	}
}

func (l *infoTraceLogger) WithPrefix(keyvals ...interface{}) InfoTraceLogger {
	return &infoTraceLogger{
		infoLogger:  l.infoLogger.WithPrefix(keyvals...),
		traceLogger: l.traceLogger.WithPrefix(keyvals...),
	}
}

func (l *infoTraceLogger) Info(keyvals ...interface{}) error {
	// We send Info and Trace log lines down the same pipe to keep them ordered
	return l.infoLogger.Log(keyvals...)
}

func (l *infoTraceLogger) Trace(keyvals ...interface{}) error {
	return l.traceLogger.Log(keyvals...)
}

// If logged to as a plain kitlog logger presume the message is for Trace
// This favours keeping Info reasonably quiet. Note that an InfoTraceLogger
// aware adapter can make its own choices, but we tend to thing of logs from
// dependencies as less interesting than logs generated by us or specifically
// routed by us.
func (l *infoTraceLogger) Log(keyvals ...interface{}) error {
	l.Trace(keyvals...)
	return nil
}
