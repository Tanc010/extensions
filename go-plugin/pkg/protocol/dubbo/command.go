/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dubbo

import (
	"mosn.io/api"
	"mosn.io/pkg/buffer"
)

type Header struct {
	Magic   []byte
	Flag    byte
	Status  byte
	Id      uint64
	DataLen uint32

	IsEvent         bool // true: heartbeat or readonly event
	IsTwoWay        bool // true: send request and expect response, false: just request without response
	Direction       int  // 1 mean req
	SerializationId int  // 2 mean hessian
	CommonHeader
}

type Frame struct {
	Header
	rawData []byte // raw data
	payload []byte // raw payload

	data    api.IoBuffer // wrapper of data
	content api.IoBuffer // wrapper of payload
}

var _ api.XFrame = &Frame{}

// ~ XFrame
func (r *Frame) GetRequestId() uint64 {
	return r.Header.Id
}

func (r *Frame) SetRequestId(id uint64) {
	r.Header.Id = id
}

func (r *Frame) IsHeartbeatFrame() bool {
	return r.Header.IsEvent
}

// dubbo frame returns default timeout
// TODO: use dubbo timeout?
func (r *Frame) GetTimeout() int32 {
	return 0
}

func (r *Frame) GetStreamType() api.StreamType {
	switch r.Direction {
	case EventRequest:
		return api.Request
	case EventResponse:
		return api.Response
	default:
		return api.Request
	}
}

func (r *Frame) GetHeader() api.HeaderMap {
	return r
}

func (r *Frame) GetData() api.IoBuffer {
	return r.content
}

func (r *Frame) SetData(data api.IoBuffer) {
	r.content = data
	r.payload = data.Bytes()
	r.DataLen = uint32(data.Len())
}

func (r *Frame) GetStatusCode() uint32 {
	return uint32(r.Header.Status)
}

func (r *Frame) Clone() api.HeaderMap {
	clone := &Frame{
		rawData: make([]byte, len(r.rawData)),
		payload: make([]byte, len(r.payload)),
	}
	clone.Header = r.Header
	copy(clone.rawData, r.rawData)
	copy(clone.payload, r.payload)
	clone.data = buffer.NewIoBufferBytes(clone.rawData)
	clone.content = buffer.NewIoBufferBytes(clone.payload)
	return clone
}
