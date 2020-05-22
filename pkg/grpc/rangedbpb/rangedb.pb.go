// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rangedb.proto

package rangedbpb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type EventsRequest struct {
	StartingWithEventNumber uint64   `protobuf:"varint,1,opt,name=startingWithEventNumber,proto3" json:"startingWithEventNumber,omitempty"`
	XXX_NoUnkeyedLiteral    struct{} `json:"-"`
	XXX_unrecognized        []byte   `json:"-"`
	XXX_sizecache           int32    `json:"-"`
}

func (m *EventsRequest) Reset()         { *m = EventsRequest{} }
func (m *EventsRequest) String() string { return proto.CompactTextString(m) }
func (*EventsRequest) ProtoMessage()    {}
func (*EventsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_40660d2deccf2909, []int{0}
}

func (m *EventsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventsRequest.Unmarshal(m, b)
}
func (m *EventsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventsRequest.Marshal(b, m, deterministic)
}
func (m *EventsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventsRequest.Merge(m, src)
}
func (m *EventsRequest) XXX_Size() int {
	return xxx_messageInfo_EventsRequest.Size(m)
}
func (m *EventsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EventsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EventsRequest proto.InternalMessageInfo

func (m *EventsRequest) GetStartingWithEventNumber() uint64 {
	if m != nil {
		return m.StartingWithEventNumber
	}
	return 0
}

type EventsByStreamRequest struct {
	StreamName              string   `protobuf:"bytes,1,opt,name=streamName,proto3" json:"streamName,omitempty"`
	StartingWithEventNumber uint64   `protobuf:"varint,2,opt,name=startingWithEventNumber,proto3" json:"startingWithEventNumber,omitempty"`
	XXX_NoUnkeyedLiteral    struct{} `json:"-"`
	XXX_unrecognized        []byte   `json:"-"`
	XXX_sizecache           int32    `json:"-"`
}

func (m *EventsByStreamRequest) Reset()         { *m = EventsByStreamRequest{} }
func (m *EventsByStreamRequest) String() string { return proto.CompactTextString(m) }
func (*EventsByStreamRequest) ProtoMessage()    {}
func (*EventsByStreamRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_40660d2deccf2909, []int{1}
}

func (m *EventsByStreamRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventsByStreamRequest.Unmarshal(m, b)
}
func (m *EventsByStreamRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventsByStreamRequest.Marshal(b, m, deterministic)
}
func (m *EventsByStreamRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventsByStreamRequest.Merge(m, src)
}
func (m *EventsByStreamRequest) XXX_Size() int {
	return xxx_messageInfo_EventsByStreamRequest.Size(m)
}
func (m *EventsByStreamRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EventsByStreamRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EventsByStreamRequest proto.InternalMessageInfo

func (m *EventsByStreamRequest) GetStreamName() string {
	if m != nil {
		return m.StreamName
	}
	return ""
}

func (m *EventsByStreamRequest) GetStartingWithEventNumber() uint64 {
	if m != nil {
		return m.StartingWithEventNumber
	}
	return 0
}

type Record struct {
	AggregateType        string   `protobuf:"bytes,1,opt,name=AggregateType,proto3" json:"AggregateType,omitempty"`
	AggregateID          string   `protobuf:"bytes,2,opt,name=AggregateID,proto3" json:"AggregateID,omitempty"`
	GlobalSequenceNumber uint64   `protobuf:"varint,3,opt,name=GlobalSequenceNumber,proto3" json:"GlobalSequenceNumber,omitempty"`
	StreamSequenceNumber uint64   `protobuf:"varint,4,opt,name=StreamSequenceNumber,proto3" json:"StreamSequenceNumber,omitempty"`
	InsertTimestamp      uint64   `protobuf:"varint,5,opt,name=InsertTimestamp,proto3" json:"InsertTimestamp,omitempty"`
	EventID              string   `protobuf:"bytes,6,opt,name=EventID,proto3" json:"EventID,omitempty"`
	EventType            string   `protobuf:"bytes,7,opt,name=EventType,proto3" json:"EventType,omitempty"`
	Data                 string   `protobuf:"bytes,8,opt,name=Data,proto3" json:"Data,omitempty"`
	Metadata             string   `protobuf:"bytes,9,opt,name=Metadata,proto3" json:"Metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Record) Reset()         { *m = Record{} }
func (m *Record) String() string { return proto.CompactTextString(m) }
func (*Record) ProtoMessage()    {}
func (*Record) Descriptor() ([]byte, []int) {
	return fileDescriptor_40660d2deccf2909, []int{2}
}

func (m *Record) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Record.Unmarshal(m, b)
}
func (m *Record) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Record.Marshal(b, m, deterministic)
}
func (m *Record) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Record.Merge(m, src)
}
func (m *Record) XXX_Size() int {
	return xxx_messageInfo_Record.Size(m)
}
func (m *Record) XXX_DiscardUnknown() {
	xxx_messageInfo_Record.DiscardUnknown(m)
}

var xxx_messageInfo_Record proto.InternalMessageInfo

func (m *Record) GetAggregateType() string {
	if m != nil {
		return m.AggregateType
	}
	return ""
}

func (m *Record) GetAggregateID() string {
	if m != nil {
		return m.AggregateID
	}
	return ""
}

func (m *Record) GetGlobalSequenceNumber() uint64 {
	if m != nil {
		return m.GlobalSequenceNumber
	}
	return 0
}

func (m *Record) GetStreamSequenceNumber() uint64 {
	if m != nil {
		return m.StreamSequenceNumber
	}
	return 0
}

func (m *Record) GetInsertTimestamp() uint64 {
	if m != nil {
		return m.InsertTimestamp
	}
	return 0
}

func (m *Record) GetEventID() string {
	if m != nil {
		return m.EventID
	}
	return ""
}

func (m *Record) GetEventType() string {
	if m != nil {
		return m.EventType
	}
	return ""
}

func (m *Record) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *Record) GetMetadata() string {
	if m != nil {
		return m.Metadata
	}
	return ""
}

func init() {
	proto.RegisterType((*EventsRequest)(nil), "rangedbpb.EventsRequest")
	proto.RegisterType((*EventsByStreamRequest)(nil), "rangedbpb.EventsByStreamRequest")
	proto.RegisterType((*Record)(nil), "rangedbpb.Record")
}

func init() { proto.RegisterFile("rangedb.proto", fileDescriptor_40660d2deccf2909) }

var fileDescriptor_40660d2deccf2909 = []byte{
	// 349 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0x4f, 0x4f, 0xfa, 0x40,
	0x10, 0xfd, 0x95, 0x1f, 0x16, 0x3a, 0x06, 0x0d, 0x13, 0x8d, 0x0d, 0x31, 0x86, 0x10, 0x0f, 0x9c,
	0x88, 0xc1, 0x8b, 0x1e, 0x25, 0x35, 0xa6, 0x07, 0x39, 0x14, 0x12, 0xcf, 0xbb, 0x30, 0xa9, 0x4d,
	0xe8, 0x1f, 0x76, 0x17, 0x13, 0xbe, 0xa2, 0x77, 0xbf, 0x8f, 0xe9, 0x14, 0x10, 0x2a, 0x70, 0xf0,
	0xb6, 0xf3, 0xde, 0xce, 0x9b, 0x37, 0x99, 0x07, 0x0d, 0x25, 0x92, 0x90, 0xa6, 0xb2, 0x97, 0xa9,
	0xd4, 0xa4, 0xe8, 0xac, 0xca, 0x4c, 0x76, 0x7c, 0x68, 0x3c, 0x7f, 0x50, 0x62, 0x74, 0x40, 0xf3,
	0x05, 0x69, 0x83, 0x0f, 0x70, 0xa5, 0x8d, 0x50, 0x26, 0x4a, 0xc2, 0xb7, 0xc8, 0xbc, 0x33, 0x39,
	0x5c, 0xc4, 0x92, 0x94, 0x6b, 0xb5, 0xad, 0x6e, 0x35, 0x38, 0x44, 0x77, 0xe6, 0x70, 0x59, 0x48,
	0x0d, 0x96, 0x23, 0xa3, 0x48, 0xc4, 0x6b, 0xc9, 0x1b, 0x00, 0xcd, 0xc0, 0x50, 0xc4, 0xc4, 0x2a,
	0x4e, 0xb0, 0x85, 0x1c, 0x1b, 0x59, 0x39, 0x3e, 0xf2, 0xb3, 0x02, 0x76, 0x40, 0x93, 0x54, 0x4d,
	0xf1, 0x16, 0x1a, 0x4f, 0x61, 0xa8, 0x28, 0x14, 0x86, 0xc6, 0xcb, 0x6c, 0x3d, 0x67, 0x17, 0xc4,
	0x36, 0x9c, 0x6e, 0x00, 0xdf, 0x63, 0x79, 0x27, 0xd8, 0x86, 0xb0, 0x0f, 0x17, 0x2f, 0xb3, 0x54,
	0x8a, 0xd9, 0x28, 0x77, 0x9f, 0x4c, 0x68, 0xe5, 0xe4, 0x3f, 0x3b, 0xd9, 0xcb, 0xe5, 0x3d, 0xc5,
	0xc6, 0xa5, 0x9e, 0x6a, 0xd1, 0xb3, 0x8f, 0xc3, 0x2e, 0x9c, 0xfb, 0x89, 0x26, 0x65, 0xc6, 0x51,
	0x4c, 0xda, 0x88, 0x38, 0x73, 0x4f, 0xf8, 0x7b, 0x19, 0x46, 0x17, 0x6a, 0xbc, 0xb3, 0xef, 0xb9,
	0x36, 0xfb, 0x5d, 0x97, 0x78, 0x0d, 0x0e, 0x3f, 0x79, 0xdf, 0x1a, 0x73, 0x3f, 0x00, 0x22, 0x54,
	0x3d, 0x61, 0x84, 0x5b, 0x67, 0x82, 0xdf, 0xd8, 0x82, 0xfa, 0x2b, 0x19, 0x31, 0xcd, 0x71, 0x87,
	0xf1, 0x4d, 0xdd, 0xff, 0xb2, 0xa0, 0x16, 0xe4, 0xc1, 0xf0, 0x06, 0xf8, 0x08, 0x76, 0x71, 0x4b,
	0x74, 0x7b, 0x9b, 0xb0, 0xf4, 0x76, 0x92, 0xd2, 0x6a, 0x6e, 0x31, 0xc5, 0x11, 0x3a, 0xff, 0xee,
	0x2c, 0xf4, 0xe1, 0x6c, 0x37, 0x06, 0xd8, 0xfe, 0x25, 0x51, 0x4a, 0xc8, 0x21, 0x29, 0x0f, 0x9a,
	0xa3, 0x85, 0xd4, 0x13, 0x15, 0x49, 0x1a, 0xa7, 0x7f, 0x34, 0x24, 0x6d, 0x0e, 0xfd, 0xfd, 0x77,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x6c, 0xe6, 0xf6, 0x7b, 0x05, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RangeDBClient is the client API for RangeDB service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RangeDBClient interface {
	Events(ctx context.Context, in *EventsRequest, opts ...grpc.CallOption) (RangeDB_EventsClient, error)
	EventsByStream(ctx context.Context, in *EventsByStreamRequest, opts ...grpc.CallOption) (RangeDB_EventsByStreamClient, error)
	SubscribeToEvents(ctx context.Context, in *EventsRequest, opts ...grpc.CallOption) (RangeDB_SubscribeToEventsClient, error)
}

type rangeDBClient struct {
	cc grpc.ClientConnInterface
}

func NewRangeDBClient(cc grpc.ClientConnInterface) RangeDBClient {
	return &rangeDBClient{cc}
}

func (c *rangeDBClient) Events(ctx context.Context, in *EventsRequest, opts ...grpc.CallOption) (RangeDB_EventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RangeDB_serviceDesc.Streams[0], "/rangedbpb.RangeDB/Events", opts...)
	if err != nil {
		return nil, err
	}
	x := &rangeDBEventsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RangeDB_EventsClient interface {
	Recv() (*Record, error)
	grpc.ClientStream
}

type rangeDBEventsClient struct {
	grpc.ClientStream
}

func (x *rangeDBEventsClient) Recv() (*Record, error) {
	m := new(Record)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *rangeDBClient) EventsByStream(ctx context.Context, in *EventsByStreamRequest, opts ...grpc.CallOption) (RangeDB_EventsByStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RangeDB_serviceDesc.Streams[1], "/rangedbpb.RangeDB/EventsByStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &rangeDBEventsByStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RangeDB_EventsByStreamClient interface {
	Recv() (*Record, error)
	grpc.ClientStream
}

type rangeDBEventsByStreamClient struct {
	grpc.ClientStream
}

func (x *rangeDBEventsByStreamClient) Recv() (*Record, error) {
	m := new(Record)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *rangeDBClient) SubscribeToEvents(ctx context.Context, in *EventsRequest, opts ...grpc.CallOption) (RangeDB_SubscribeToEventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RangeDB_serviceDesc.Streams[2], "/rangedbpb.RangeDB/SubscribeToEvents", opts...)
	if err != nil {
		return nil, err
	}
	x := &rangeDBSubscribeToEventsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RangeDB_SubscribeToEventsClient interface {
	Recv() (*Record, error)
	grpc.ClientStream
}

type rangeDBSubscribeToEventsClient struct {
	grpc.ClientStream
}

func (x *rangeDBSubscribeToEventsClient) Recv() (*Record, error) {
	m := new(Record)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RangeDBServer is the server API for RangeDB service.
type RangeDBServer interface {
	Events(*EventsRequest, RangeDB_EventsServer) error
	EventsByStream(*EventsByStreamRequest, RangeDB_EventsByStreamServer) error
	SubscribeToEvents(*EventsRequest, RangeDB_SubscribeToEventsServer) error
}

// UnimplementedRangeDBServer can be embedded to have forward compatible implementations.
type UnimplementedRangeDBServer struct {
}

func (*UnimplementedRangeDBServer) Events(req *EventsRequest, srv RangeDB_EventsServer) error {
	return status.Errorf(codes.Unimplemented, "method Events not implemented")
}
func (*UnimplementedRangeDBServer) EventsByStream(req *EventsByStreamRequest, srv RangeDB_EventsByStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method EventsByStream not implemented")
}
func (*UnimplementedRangeDBServer) SubscribeToEvents(req *EventsRequest, srv RangeDB_SubscribeToEventsServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeToEvents not implemented")
}

func RegisterRangeDBServer(s *grpc.Server, srv RangeDBServer) {
	s.RegisterService(&_RangeDB_serviceDesc, srv)
}

func _RangeDB_Events_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EventsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RangeDBServer).Events(m, &rangeDBEventsServer{stream})
}

type RangeDB_EventsServer interface {
	Send(*Record) error
	grpc.ServerStream
}

type rangeDBEventsServer struct {
	grpc.ServerStream
}

func (x *rangeDBEventsServer) Send(m *Record) error {
	return x.ServerStream.SendMsg(m)
}

func _RangeDB_EventsByStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EventsByStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RangeDBServer).EventsByStream(m, &rangeDBEventsByStreamServer{stream})
}

type RangeDB_EventsByStreamServer interface {
	Send(*Record) error
	grpc.ServerStream
}

type rangeDBEventsByStreamServer struct {
	grpc.ServerStream
}

func (x *rangeDBEventsByStreamServer) Send(m *Record) error {
	return x.ServerStream.SendMsg(m)
}

func _RangeDB_SubscribeToEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EventsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RangeDBServer).SubscribeToEvents(m, &rangeDBSubscribeToEventsServer{stream})
}

type RangeDB_SubscribeToEventsServer interface {
	Send(*Record) error
	grpc.ServerStream
}

type rangeDBSubscribeToEventsServer struct {
	grpc.ServerStream
}

func (x *rangeDBSubscribeToEventsServer) Send(m *Record) error {
	return x.ServerStream.SendMsg(m)
}

var _RangeDB_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rangedbpb.RangeDB",
	HandlerType: (*RangeDBServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Events",
			Handler:       _RangeDB_Events_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "EventsByStream",
			Handler:       _RangeDB_EventsByStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeToEvents",
			Handler:       _RangeDB_SubscribeToEvents_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "rangedb.proto",
}
