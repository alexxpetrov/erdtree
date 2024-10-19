// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: internal/api/erdtree.proto

package dbv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Operation int32

const (
	Operation_OPERATION_UNSPECIFIED Operation = 0
	Operation_OPERATION_SET         Operation = 1
	Operation_OPERATION_DELETE      Operation = 2
)

// Enum value maps for Operation.
var (
	Operation_name = map[int32]string{
		0: "OPERATION_UNSPECIFIED",
		1: "OPERATION_SET",
		2: "OPERATION_DELETE",
	}
	Operation_value = map[string]int32{
		"OPERATION_UNSPECIFIED": 0,
		"OPERATION_SET":         1,
		"OPERATION_DELETE":      2,
	}
)

func (x Operation) Enum() *Operation {
	p := new(Operation)
	*p = x
	return p
}

func (x Operation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Operation) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_api_erdtree_proto_enumTypes[0].Descriptor()
}

func (Operation) Type() protoreflect.EnumType {
	return &file_internal_api_erdtree_proto_enumTypes[0]
}

func (x Operation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Operation.Descriptor instead.
func (Operation) EnumDescriptor() ([]byte, []int) {
	return file_internal_api_erdtree_proto_rawDescGZIP(), []int{0}
}

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_erdtree_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_erdtree_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_internal_api_erdtree_proto_rawDescGZIP(), []int{0}
}

func (x *GetRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type GetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Exists bool   `protobuf:"varint,1,opt,name=exists,proto3" json:"exists,omitempty"`
	Value  []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_erdtree_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_erdtree_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_erdtree_proto_rawDescGZIP(), []int{1}
}

func (x *GetResponse) GetExists() bool {
	if x != nil {
		return x.Exists
	}
	return false
}

func (x *GetResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type SetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"` // uint64 ttl = 3; // Time-to-live in seconds (optional)
}

func (x *SetRequest) Reset() {
	*x = SetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_erdtree_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetRequest) ProtoMessage() {}

func (x *SetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_erdtree_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetRequest.ProtoReflect.Descriptor instead.
func (*SetRequest) Descriptor() ([]byte, []int) {
	return file_internal_api_erdtree_proto_rawDescGZIP(), []int{2}
}

func (x *SetRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SetRequest) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type SetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *SetResponse) Reset() {
	*x = SetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_erdtree_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetResponse) ProtoMessage() {}

func (x *SetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_erdtree_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetResponse.ProtoReflect.Descriptor instead.
func (*SetResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_erdtree_proto_rawDescGZIP(), []int{3}
}

func (x *SetResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_erdtree_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_erdtree_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_internal_api_erdtree_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type DeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *DeleteResponse) Reset() {
	*x = DeleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_erdtree_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteResponse) ProtoMessage() {}

func (x *DeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_erdtree_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteResponse.ProtoReflect.Descriptor instead.
func (*DeleteResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_erdtree_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type ReplicationSyncRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LastSyncTimestamp uint64 `protobuf:"varint,1,opt,name=last_sync_timestamp,json=lastSyncTimestamp,proto3" json:"last_sync_timestamp,omitempty"`
}

func (x *ReplicationSyncRequest) Reset() {
	*x = ReplicationSyncRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_erdtree_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplicationSyncRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplicationSyncRequest) ProtoMessage() {}

func (x *ReplicationSyncRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_erdtree_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplicationSyncRequest.ProtoReflect.Descriptor instead.
func (*ReplicationSyncRequest) Descriptor() ([]byte, []int) {
	return file_internal_api_erdtree_proto_rawDescGZIP(), []int{6}
}

func (x *ReplicationSyncRequest) GetLastSyncTimestamp() uint64 {
	if x != nil {
		return x.LastSyncTimestamp
	}
	return 0
}

type LogEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp int64     `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Operation Operation `protobuf:"varint,2,opt,name=operation,proto3,enum=erdtree.v1.Operation" json:"operation,omitempty"`
	Key       string    `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Value     []byte    `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	ExpiresAt int64     `protobuf:"varint,5,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
}

func (x *LogEntry) Reset() {
	*x = LogEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_erdtree_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogEntry) ProtoMessage() {}

func (x *LogEntry) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_erdtree_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogEntry.ProtoReflect.Descriptor instead.
func (*LogEntry) Descriptor() ([]byte, []int) {
	return file_internal_api_erdtree_proto_rawDescGZIP(), []int{7}
}

func (x *LogEntry) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *LogEntry) GetOperation() Operation {
	if x != nil {
		return x.Operation
	}
	return Operation_OPERATION_UNSPECIFIED
}

func (x *LogEntry) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *LogEntry) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *LogEntry) GetExpiresAt() int64 {
	if x != nil {
		return x.ExpiresAt
	}
	return 0
}

type ReplicationSyncResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries          []*LogEntry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
	CurrentTimestamp uint64      `protobuf:"varint,2,opt,name=current_timestamp,json=currentTimestamp,proto3" json:"current_timestamp,omitempty"`
}

func (x *ReplicationSyncResponse) Reset() {
	*x = ReplicationSyncResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_api_erdtree_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplicationSyncResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplicationSyncResponse) ProtoMessage() {}

func (x *ReplicationSyncResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_erdtree_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplicationSyncResponse.ProtoReflect.Descriptor instead.
func (*ReplicationSyncResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_erdtree_proto_rawDescGZIP(), []int{8}
}

func (x *ReplicationSyncResponse) GetEntries() []*LogEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

func (x *ReplicationSyncResponse) GetCurrentTimestamp() uint64 {
	if x != nil {
		return x.CurrentTimestamp
	}
	return 0
}

var File_internal_api_erdtree_proto protoreflect.FileDescriptor

var file_internal_api_erdtree_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65,
	0x72, 0x64, 0x74, 0x72, 0x65, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x65, 0x72,
	0x64, 0x74, 0x72, 0x65, 0x65, 0x2e, 0x76, 0x31, 0x22, 0x1e, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x3b, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x78, 0x69, 0x73, 0x74,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x65, 0x78, 0x69, 0x73, 0x74, 0x73, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x34, 0x0a, 0x0a, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x27, 0x0a, 0x0b, 0x53,
	0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x22, 0x21, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x2a, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x22, 0x48, 0x0a, 0x16, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a,
	0x13, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x11, 0x6c, 0x61, 0x73, 0x74,
	0x53, 0x79, 0x6e, 0x63, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0xa4, 0x01,
	0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x33, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x65, 0x72,
	0x64, 0x74, 0x72, 0x65, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73,
	0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72,
	0x65, 0x73, 0x41, 0x74, 0x22, 0x76, 0x0a, 0x17, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2e, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x65, 0x72, 0x64, 0x74, 0x72, 0x65, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f,
	0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12,
	0x2b, 0x0a, 0x11, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2a, 0x4f, 0x0a, 0x09,
	0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x15, 0x4f, 0x50, 0x45,
	0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x53, 0x45, 0x54, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x4f, 0x50, 0x45, 0x52, 0x41,
	0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x02, 0x32, 0xa3, 0x02,
	0x0a, 0x0c, 0x45, 0x72, 0x64, 0x74, 0x72, 0x65, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x38,
	0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x16, 0x2e, 0x65, 0x72, 0x64, 0x74, 0x72, 0x65, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e,
	0x65, 0x72, 0x64, 0x74, 0x72, 0x65, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x03, 0x53, 0x65, 0x74, 0x12,
	0x16, 0x2e, 0x65, 0x72, 0x64, 0x74, 0x72, 0x65, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x65, 0x72, 0x64, 0x74, 0x72, 0x65,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x41, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x19, 0x2e, 0x65,
	0x72, 0x64, 0x74, 0x72, 0x65, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x65, 0x72, 0x64, 0x74, 0x72, 0x65,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5c, 0x0a, 0x0f, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x79, 0x6e, 0x63, 0x12, 0x22, 0x2e, 0x65, 0x72, 0x64, 0x74, 0x72,
	0x65, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x65,
	0x72, 0x64, 0x74, 0x72, 0x65, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x42, 0x5a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6f, 0x6c, 0x65, 0x6b, 0x73, 0x69, 0x69, 0x70, 0x2d, 0x61, 0x69, 0x6f, 0x6c, 0x61,
	0x2f, 0x65, 0x72, 0x64, 0x74, 0x72, 0x65, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61,
	0x70, 0x69, 0x3b, 0x64, 0x62, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_api_erdtree_proto_rawDescOnce sync.Once
	file_internal_api_erdtree_proto_rawDescData = file_internal_api_erdtree_proto_rawDesc
)

func file_internal_api_erdtree_proto_rawDescGZIP() []byte {
	file_internal_api_erdtree_proto_rawDescOnce.Do(func() {
		file_internal_api_erdtree_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_api_erdtree_proto_rawDescData)
	})
	return file_internal_api_erdtree_proto_rawDescData
}

var file_internal_api_erdtree_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_api_erdtree_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_internal_api_erdtree_proto_goTypes = []any{
	(Operation)(0),                  // 0: erdtree.v1.Operation
	(*GetRequest)(nil),              // 1: erdtree.v1.GetRequest
	(*GetResponse)(nil),             // 2: erdtree.v1.GetResponse
	(*SetRequest)(nil),              // 3: erdtree.v1.SetRequest
	(*SetResponse)(nil),             // 4: erdtree.v1.SetResponse
	(*DeleteRequest)(nil),           // 5: erdtree.v1.DeleteRequest
	(*DeleteResponse)(nil),          // 6: erdtree.v1.DeleteResponse
	(*ReplicationSyncRequest)(nil),  // 7: erdtree.v1.ReplicationSyncRequest
	(*LogEntry)(nil),                // 8: erdtree.v1.LogEntry
	(*ReplicationSyncResponse)(nil), // 9: erdtree.v1.ReplicationSyncResponse
}
var file_internal_api_erdtree_proto_depIdxs = []int32{
	0, // 0: erdtree.v1.LogEntry.operation:type_name -> erdtree.v1.Operation
	8, // 1: erdtree.v1.ReplicationSyncResponse.entries:type_name -> erdtree.v1.LogEntry
	1, // 2: erdtree.v1.ErdtreeStore.Get:input_type -> erdtree.v1.GetRequest
	3, // 3: erdtree.v1.ErdtreeStore.Set:input_type -> erdtree.v1.SetRequest
	5, // 4: erdtree.v1.ErdtreeStore.Delete:input_type -> erdtree.v1.DeleteRequest
	7, // 5: erdtree.v1.ErdtreeStore.ReplicationSync:input_type -> erdtree.v1.ReplicationSyncRequest
	2, // 6: erdtree.v1.ErdtreeStore.Get:output_type -> erdtree.v1.GetResponse
	4, // 7: erdtree.v1.ErdtreeStore.Set:output_type -> erdtree.v1.SetResponse
	6, // 8: erdtree.v1.ErdtreeStore.Delete:output_type -> erdtree.v1.DeleteResponse
	9, // 9: erdtree.v1.ErdtreeStore.ReplicationSync:output_type -> erdtree.v1.ReplicationSyncResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_internal_api_erdtree_proto_init() }
func file_internal_api_erdtree_proto_init() {
	if File_internal_api_erdtree_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_api_erdtree_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_api_erdtree_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_api_erdtree_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_api_erdtree_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*SetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_api_erdtree_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_api_erdtree_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_api_erdtree_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*ReplicationSyncRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_api_erdtree_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*LogEntry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_api_erdtree_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*ReplicationSyncResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_api_erdtree_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_api_erdtree_proto_goTypes,
		DependencyIndexes: file_internal_api_erdtree_proto_depIdxs,
		EnumInfos:         file_internal_api_erdtree_proto_enumTypes,
		MessageInfos:      file_internal_api_erdtree_proto_msgTypes,
	}.Build()
	File_internal_api_erdtree_proto = out.File
	file_internal_api_erdtree_proto_rawDesc = nil
	file_internal_api_erdtree_proto_goTypes = nil
	file_internal_api_erdtree_proto_depIdxs = nil
}
