// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.12.4
// source: saldo.proto

package pb

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Saldo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SaldoId        int32                `protobuf:"varint,1,opt,name=saldo_id,json=saldoId,proto3" json:"saldo_id,omitempty"`
	UserId         int32                `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	TotalBalance   int32                `protobuf:"varint,3,opt,name=total_balance,json=totalBalance,proto3" json:"total_balance,omitempty"`
	WithdrawTime   *timestamp.Timestamp `protobuf:"bytes,4,opt,name=withdraw_time,json=withdrawTime,proto3" json:"withdraw_time,omitempty"`
	WithdrawAmount int32                `protobuf:"varint,5,opt,name=withdraw_amount,json=withdrawAmount,proto3" json:"withdraw_amount,omitempty"`
	CreatedAt      *timestamp.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt      *timestamp.Timestamp `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Saldo) Reset() {
	*x = Saldo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saldo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Saldo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Saldo) ProtoMessage() {}

func (x *Saldo) ProtoReflect() protoreflect.Message {
	mi := &file_saldo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Saldo.ProtoReflect.Descriptor instead.
func (*Saldo) Descriptor() ([]byte, []int) {
	return file_saldo_proto_rawDescGZIP(), []int{0}
}

func (x *Saldo) GetSaldoId() int32 {
	if x != nil {
		return x.SaldoId
	}
	return 0
}

func (x *Saldo) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Saldo) GetTotalBalance() int32 {
	if x != nil {
		return x.TotalBalance
	}
	return 0
}

func (x *Saldo) GetWithdrawTime() *timestamp.Timestamp {
	if x != nil {
		return x.WithdrawTime
	}
	return nil
}

func (x *Saldo) GetWithdrawAmount() int32 {
	if x != nil {
		return x.WithdrawAmount
	}
	return 0
}

func (x *Saldo) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Saldo) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type SaldoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *SaldoRequest) Reset() {
	*x = SaldoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saldo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaldoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaldoRequest) ProtoMessage() {}

func (x *SaldoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_saldo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaldoRequest.ProtoReflect.Descriptor instead.
func (*SaldoRequest) Descriptor() ([]byte, []int) {
	return file_saldo_proto_rawDescGZIP(), []int{1}
}

func (x *SaldoRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type SaldoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Saldo *Saldo `protobuf:"bytes,1,opt,name=saldo,proto3" json:"saldo,omitempty"`
}

func (x *SaldoResponse) Reset() {
	*x = SaldoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saldo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaldoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaldoResponse) ProtoMessage() {}

func (x *SaldoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_saldo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaldoResponse.ProtoReflect.Descriptor instead.
func (*SaldoResponse) Descriptor() ([]byte, []int) {
	return file_saldo_proto_rawDescGZIP(), []int{2}
}

func (x *SaldoResponse) GetSaldo() *Saldo {
	if x != nil {
		return x.Saldo
	}
	return nil
}

type SaldoResponses struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Saldos []*Saldo `protobuf:"bytes,1,rep,name=saldos,proto3" json:"saldos,omitempty"`
}

func (x *SaldoResponses) Reset() {
	*x = SaldoResponses{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saldo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaldoResponses) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaldoResponses) ProtoMessage() {}

func (x *SaldoResponses) ProtoReflect() protoreflect.Message {
	mi := &file_saldo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaldoResponses.ProtoReflect.Descriptor instead.
func (*SaldoResponses) Descriptor() ([]byte, []int) {
	return file_saldo_proto_rawDescGZIP(), []int{3}
}

func (x *SaldoResponses) GetSaldos() []*Saldo {
	if x != nil {
		return x.Saldos
	}
	return nil
}

type CreateSaldoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId       int32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	TotalBalance int32 `protobuf:"varint,2,opt,name=total_balance,json=totalBalance,proto3" json:"total_balance,omitempty"`
}

func (x *CreateSaldoRequest) Reset() {
	*x = CreateSaldoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saldo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSaldoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSaldoRequest) ProtoMessage() {}

func (x *CreateSaldoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_saldo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSaldoRequest.ProtoReflect.Descriptor instead.
func (*CreateSaldoRequest) Descriptor() ([]byte, []int) {
	return file_saldo_proto_rawDescGZIP(), []int{4}
}

func (x *CreateSaldoRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *CreateSaldoRequest) GetTotalBalance() int32 {
	if x != nil {
		return x.TotalBalance
	}
	return 0
}

type UpdateSaldoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SaldoId        int32                `protobuf:"varint,1,opt,name=saldo_id,json=saldoId,proto3" json:"saldo_id,omitempty"`
	UserId         int32                `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	TotalBalance   int32                `protobuf:"varint,3,opt,name=total_balance,json=totalBalance,proto3" json:"total_balance,omitempty"`
	WithdrawAmount int32                `protobuf:"varint,4,opt,name=withdraw_amount,json=withdrawAmount,proto3" json:"withdraw_amount,omitempty"`
	WithdrawTime   *timestamp.Timestamp `protobuf:"bytes,5,opt,name=withdraw_time,json=withdrawTime,proto3" json:"withdraw_time,omitempty"`
}

func (x *UpdateSaldoRequest) Reset() {
	*x = UpdateSaldoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saldo_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateSaldoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateSaldoRequest) ProtoMessage() {}

func (x *UpdateSaldoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_saldo_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateSaldoRequest.ProtoReflect.Descriptor instead.
func (*UpdateSaldoRequest) Descriptor() ([]byte, []int) {
	return file_saldo_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateSaldoRequest) GetSaldoId() int32 {
	if x != nil {
		return x.SaldoId
	}
	return 0
}

func (x *UpdateSaldoRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UpdateSaldoRequest) GetTotalBalance() int32 {
	if x != nil {
		return x.TotalBalance
	}
	return 0
}

func (x *UpdateSaldoRequest) GetWithdrawAmount() int32 {
	if x != nil {
		return x.WithdrawAmount
	}
	return 0
}

func (x *UpdateSaldoRequest) GetWithdrawTime() *timestamp.Timestamp {
	if x != nil {
		return x.WithdrawTime
	}
	return nil
}

type DeleteSaldoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *DeleteSaldoResponse) Reset() {
	*x = DeleteSaldoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_saldo_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteSaldoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSaldoResponse) ProtoMessage() {}

func (x *DeleteSaldoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_saldo_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteSaldoResponse.ProtoReflect.Descriptor instead.
func (*DeleteSaldoResponse) Descriptor() ([]byte, []int) {
	return file_saldo_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteSaldoResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_saldo_proto protoreflect.FileDescriptor

var file_saldo_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x73, 0x61, 0x6c, 0x64, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70,
	0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xc0, 0x02, 0x0a, 0x05, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x61, 0x6c,
	0x64, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x73, 0x61, 0x6c,
	0x64, 0x6f, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x23, 0x0a,
	0x0d, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x42, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x12, 0x3f, 0x0a, 0x0d, 0x77, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x77, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x77, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x5f,
	0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x77, 0x69,
	0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x39, 0x0a, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x22, 0x1e, 0x0a, 0x0c, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02,
	0x69, 0x64, 0x22, 0x30, 0x0a, 0x0d, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x05, 0x73, 0x61, 0x6c, 0x64, 0x6f, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x05, 0x73,
	0x61, 0x6c, 0x64, 0x6f, 0x22, 0x33, 0x0a, 0x0e, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x12, 0x21, 0x0a, 0x06, 0x73, 0x61, 0x6c, 0x64, 0x6f, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x61, 0x6c, 0x64,
	0x6f, 0x52, 0x06, 0x73, 0x61, 0x6c, 0x64, 0x6f, 0x73, 0x22, 0x52, 0x0a, 0x12, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x22, 0xd7, 0x01,
	0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x61, 0x6c, 0x64, 0x6f, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x73, 0x61, 0x6c, 0x64, 0x6f, 0x49, 0x64, 0x12,
	0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x27, 0x0a,
	0x0f, 0x77, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x77, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77,
	0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x3f, 0x0a, 0x0d, 0x77, 0x69, 0x74, 0x68, 0x64, 0x72,
	0x61, 0x77, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x77, 0x69, 0x74, 0x68, 0x64,
	0x72, 0x61, 0x77, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x2f, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x32, 0x98, 0x03, 0x0a, 0x0c, 0x53, 0x61, 0x6c,
	0x64, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x37, 0x0a, 0x09, 0x47, 0x65, 0x74,
	0x53, 0x61, 0x6c, 0x64, 0x6f, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x12,
	0x2e, 0x70, 0x62, 0x2e, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x73, 0x12, 0x2f, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x12, 0x10,
	0x2e, 0x70, 0x62, 0x2e, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x42,
	0x79, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x61, 0x6c, 0x64,
	0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x61,
	0x6c, 0x64, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x12, 0x37, 0x0a, 0x10,
	0x47, 0x65, 0x74, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53,
	0x61, 0x6c, 0x64, 0x6f, 0x12, 0x16, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70,
	0x62, 0x2e, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x38, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x12, 0x16,
	0x2e, 0x70, 0x62, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x61, 0x6c, 0x64,
	0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0b, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x12, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x61,
	0x6c, 0x64, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x62, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x61, 0x6c, 0x64, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x2b, 0x5a, 0x29, 0x4d, 0x61, 0x6d, 0x61, 0x6e, 0x67, 0x52, 0x75, 0x73,
	0x74, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_saldo_proto_rawDescOnce sync.Once
	file_saldo_proto_rawDescData = file_saldo_proto_rawDesc
)

func file_saldo_proto_rawDescGZIP() []byte {
	file_saldo_proto_rawDescOnce.Do(func() {
		file_saldo_proto_rawDescData = protoimpl.X.CompressGZIP(file_saldo_proto_rawDescData)
	})
	return file_saldo_proto_rawDescData
}

var file_saldo_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_saldo_proto_goTypes = []interface{}{
	(*Saldo)(nil),               // 0: pb.Saldo
	(*SaldoRequest)(nil),        // 1: pb.SaldoRequest
	(*SaldoResponse)(nil),       // 2: pb.SaldoResponse
	(*SaldoResponses)(nil),      // 3: pb.SaldoResponses
	(*CreateSaldoRequest)(nil),  // 4: pb.CreateSaldoRequest
	(*UpdateSaldoRequest)(nil),  // 5: pb.UpdateSaldoRequest
	(*DeleteSaldoResponse)(nil), // 6: pb.DeleteSaldoResponse
	(*timestamp.Timestamp)(nil), // 7: google.protobuf.Timestamp
	(*empty.Empty)(nil),         // 8: google.protobuf.Empty
}
var file_saldo_proto_depIdxs = []int32{
	7,  // 0: pb.Saldo.withdraw_time:type_name -> google.protobuf.Timestamp
	7,  // 1: pb.Saldo.created_at:type_name -> google.protobuf.Timestamp
	7,  // 2: pb.Saldo.updated_at:type_name -> google.protobuf.Timestamp
	0,  // 3: pb.SaldoResponse.saldo:type_name -> pb.Saldo
	0,  // 4: pb.SaldoResponses.saldos:type_name -> pb.Saldo
	7,  // 5: pb.UpdateSaldoRequest.withdraw_time:type_name -> google.protobuf.Timestamp
	8,  // 6: pb.SaldoService.GetSaldos:input_type -> google.protobuf.Empty
	1,  // 7: pb.SaldoService.GetSaldo:input_type -> pb.SaldoRequest
	1,  // 8: pb.SaldoService.GetSaldoByUsers:input_type -> pb.SaldoRequest
	1,  // 9: pb.SaldoService.GetSaldoByUserId:input_type -> pb.SaldoRequest
	4,  // 10: pb.SaldoService.CreateSaldo:input_type -> pb.CreateSaldoRequest
	5,  // 11: pb.SaldoService.UpdateSaldo:input_type -> pb.UpdateSaldoRequest
	1,  // 12: pb.SaldoService.DeleteSaldo:input_type -> pb.SaldoRequest
	3,  // 13: pb.SaldoService.GetSaldos:output_type -> pb.SaldoResponses
	2,  // 14: pb.SaldoService.GetSaldo:output_type -> pb.SaldoResponse
	3,  // 15: pb.SaldoService.GetSaldoByUsers:output_type -> pb.SaldoResponses
	2,  // 16: pb.SaldoService.GetSaldoByUserId:output_type -> pb.SaldoResponse
	2,  // 17: pb.SaldoService.CreateSaldo:output_type -> pb.SaldoResponse
	2,  // 18: pb.SaldoService.UpdateSaldo:output_type -> pb.SaldoResponse
	6,  // 19: pb.SaldoService.DeleteSaldo:output_type -> pb.DeleteSaldoResponse
	13, // [13:20] is the sub-list for method output_type
	6,  // [6:13] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_saldo_proto_init() }
func file_saldo_proto_init() {
	if File_saldo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_saldo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Saldo); i {
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
		file_saldo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaldoRequest); i {
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
		file_saldo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaldoResponse); i {
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
		file_saldo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaldoResponses); i {
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
		file_saldo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSaldoRequest); i {
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
		file_saldo_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateSaldoRequest); i {
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
		file_saldo_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteSaldoResponse); i {
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
			RawDescriptor: file_saldo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_saldo_proto_goTypes,
		DependencyIndexes: file_saldo_proto_depIdxs,
		MessageInfos:      file_saldo_proto_msgTypes,
	}.Build()
	File_saldo_proto = out.File
	file_saldo_proto_rawDesc = nil
	file_saldo_proto_goTypes = nil
	file_saldo_proto_depIdxs = nil
}
