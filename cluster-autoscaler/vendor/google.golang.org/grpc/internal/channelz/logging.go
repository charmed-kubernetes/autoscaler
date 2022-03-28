/*
 *
 * Copyright 2020 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package channelz

import (
	"fmt"

<<<<<<< HEAD
	"google.golang.org/grpc/grpclog"
)

var logger = grpclog.Component("channelz")

// Info logs and adds a trace event if channelz is on.
func Info(l grpclog.DepthLoggerV2, id int64, args ...interface{}) {
	if IsOn() {
		AddTraceEvent(l, id, 1, &TraceEventDesc{
			Desc:     fmt.Sprint(args...),
			Severity: CtInfo,
		})
	} else {
		l.InfoDepth(1, args...)
	}
}

// Infof logs and adds a trace event if channelz is on.
func Infof(l grpclog.DepthLoggerV2, id int64, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if IsOn() {
		AddTraceEvent(l, id, 1, &TraceEventDesc{
			Desc:     msg,
			Severity: CtInfo,
		})
	} else {
		l.InfoDepth(1, msg)
	}
}

// Warning logs and adds a trace event if channelz is on.
func Warning(l grpclog.DepthLoggerV2, id int64, args ...interface{}) {
	if IsOn() {
		AddTraceEvent(l, id, 1, &TraceEventDesc{
=======
	"google.golang.org/grpc/internal/grpclog"
)

// Info logs through grpclog.Info and adds a trace event if channelz is on.
func Info(id int64, args ...interface{}) {
	if IsOn() {
		AddTraceEvent(id, 1, &TraceEventDesc{
			Desc:     fmt.Sprint(args...),
			Severity: CtINFO,
		})
	} else {
		grpclog.InfoDepth(1, args...)
	}
}

// Infof logs through grpclog.Infof and adds a trace event if channelz is on.
func Infof(id int64, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if IsOn() {
		AddTraceEvent(id, 1, &TraceEventDesc{
			Desc:     msg,
			Severity: CtINFO,
		})
	} else {
		grpclog.InfoDepth(1, msg)
	}
}

// Warning logs through grpclog.Warning and adds a trace event if channelz is on.
func Warning(id int64, args ...interface{}) {
	if IsOn() {
		AddTraceEvent(id, 1, &TraceEventDesc{
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
			Desc:     fmt.Sprint(args...),
			Severity: CtWarning,
		})
	} else {
<<<<<<< HEAD
		l.WarningDepth(1, args...)
	}
}

// Warningf logs and adds a trace event if channelz is on.
func Warningf(l grpclog.DepthLoggerV2, id int64, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if IsOn() {
		AddTraceEvent(l, id, 1, &TraceEventDesc{
=======
		grpclog.WarningDepth(1, args...)
	}
}

// Warningf logs through grpclog.Warningf and adds a trace event if channelz is on.
func Warningf(id int64, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if IsOn() {
		AddTraceEvent(id, 1, &TraceEventDesc{
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
			Desc:     msg,
			Severity: CtWarning,
		})
	} else {
<<<<<<< HEAD
		l.WarningDepth(1, msg)
	}
}

// Error logs and adds a trace event if channelz is on.
func Error(l grpclog.DepthLoggerV2, id int64, args ...interface{}) {
	if IsOn() {
		AddTraceEvent(l, id, 1, &TraceEventDesc{
=======
		grpclog.WarningDepth(1, msg)
	}
}

// Error logs through grpclog.Error and adds a trace event if channelz is on.
func Error(id int64, args ...interface{}) {
	if IsOn() {
		AddTraceEvent(id, 1, &TraceEventDesc{
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
			Desc:     fmt.Sprint(args...),
			Severity: CtError,
		})
	} else {
<<<<<<< HEAD
		l.ErrorDepth(1, args...)
	}
}

// Errorf logs and adds a trace event if channelz is on.
func Errorf(l grpclog.DepthLoggerV2, id int64, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if IsOn() {
		AddTraceEvent(l, id, 1, &TraceEventDesc{
=======
		grpclog.ErrorDepth(1, args...)
	}
}

// Errorf logs through grpclog.Errorf and adds a trace event if channelz is on.
func Errorf(id int64, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if IsOn() {
		AddTraceEvent(id, 1, &TraceEventDesc{
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
			Desc:     msg,
			Severity: CtError,
		})
	} else {
<<<<<<< HEAD
		l.ErrorDepth(1, msg)
=======
		grpclog.ErrorDepth(1, msg)
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
	}
}
