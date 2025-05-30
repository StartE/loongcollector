// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: log_event.proto

package protocol

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type LogEvent struct {
	Timestamp  uint64              `protobuf:"varint,1,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	Contents   []*LogEvent_Content `protobuf:"bytes,2,rep,name=Contents,proto3" json:"Contents,omitempty"`
	Level      []byte              `protobuf:"bytes,3,opt,name=Level,proto3" json:"Level,omitempty"`
	FileOffset uint64              `protobuf:"varint,4,opt,name=FileOffset,proto3" json:"FileOffset,omitempty"`
	RawSize    uint64              `protobuf:"varint,5,opt,name=RawSize,proto3" json:"RawSize,omitempty"`
}

func (m *LogEvent) Reset()         { *m = LogEvent{} }
func (m *LogEvent) String() string { return proto.CompactTextString(m) }
func (*LogEvent) ProtoMessage()    {}
func (*LogEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_9ab154f10bf01af7, []int{0}
}
func (m *LogEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LogEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LogEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LogEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogEvent.Merge(m, src)
}
func (m *LogEvent) XXX_Size() int {
	return m.Size()
}
func (m *LogEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_LogEvent.DiscardUnknown(m)
}

var xxx_messageInfo_LogEvent proto.InternalMessageInfo

func (m *LogEvent) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *LogEvent) GetContents() []*LogEvent_Content {
	if m != nil {
		return m.Contents
	}
	return nil
}

func (m *LogEvent) GetLevel() []byte {
	if m != nil {
		return m.Level
	}
	return nil
}

func (m *LogEvent) GetFileOffset() uint64 {
	if m != nil {
		return m.FileOffset
	}
	return 0
}

func (m *LogEvent) GetRawSize() uint64 {
	if m != nil {
		return m.RawSize
	}
	return 0
}

type LogEvent_Content struct {
	Key   []byte `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (m *LogEvent_Content) Reset()         { *m = LogEvent_Content{} }
func (m *LogEvent_Content) String() string { return proto.CompactTextString(m) }
func (*LogEvent_Content) ProtoMessage()    {}
func (*LogEvent_Content) Descriptor() ([]byte, []int) {
	return fileDescriptor_9ab154f10bf01af7, []int{0, 0}
}
func (m *LogEvent_Content) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LogEvent_Content) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LogEvent_Content.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LogEvent_Content) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogEvent_Content.Merge(m, src)
}
func (m *LogEvent_Content) XXX_Size() int {
	return m.Size()
}
func (m *LogEvent_Content) XXX_DiscardUnknown() {
	xxx_messageInfo_LogEvent_Content.DiscardUnknown(m)
}

var xxx_messageInfo_LogEvent_Content proto.InternalMessageInfo

func (m *LogEvent_Content) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *LogEvent_Content) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*LogEvent)(nil), "logtail.models.LogEvent")
	proto.RegisterType((*LogEvent_Content)(nil), "logtail.models.LogEvent.Content")
}

func init() { proto.RegisterFile("log_event.proto", fileDescriptor_9ab154f10bf01af7) }

var fileDescriptor_9ab154f10bf01af7 = []byte{
	// 242 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcf, 0xc9, 0x4f, 0x8f,
	0x4f, 0x2d, 0x4b, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xcb, 0xc9, 0x4f,
	0x2f, 0x49, 0xcc, 0xcc, 0xd1, 0xcb, 0xcd, 0x4f, 0x49, 0xcd, 0x29, 0x56, 0x7a, 0xc9, 0xc8, 0xc5,
	0xe1, 0x93, 0x9f, 0xee, 0x0a, 0x52, 0x22, 0x24, 0xc3, 0xc5, 0x19, 0x92, 0x99, 0x9b, 0x5a, 0x5c,
	0x92, 0x98, 0x5b, 0x20, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x12, 0x84, 0x10, 0x10, 0xb2, 0xe1, 0xe2,
	0x70, 0xce, 0xcf, 0x2b, 0x49, 0xcd, 0x2b, 0x29, 0x96, 0x60, 0x52, 0x60, 0xd6, 0xe0, 0x36, 0x52,
	0xd0, 0x43, 0x35, 0x4d, 0x0f, 0x66, 0x92, 0x1e, 0x54, 0x61, 0x10, 0x5c, 0x87, 0x90, 0x08, 0x17,
	0xab, 0x4f, 0x6a, 0x59, 0x6a, 0x8e, 0x04, 0xb3, 0x02, 0xa3, 0x06, 0x4f, 0x10, 0x84, 0x23, 0x24,
	0xc7, 0xc5, 0xe5, 0x96, 0x99, 0x93, 0xea, 0x9f, 0x96, 0x56, 0x9c, 0x5a, 0x22, 0xc1, 0x02, 0xb6,
	0x12, 0x49, 0x44, 0x48, 0x82, 0x8b, 0x3d, 0x28, 0xb1, 0x3c, 0x38, 0xb3, 0x2a, 0x55, 0x82, 0x15,
	0x2c, 0x09, 0xe3, 0x4a, 0x19, 0x72, 0xb1, 0x43, 0xcd, 0x16, 0x12, 0xe0, 0x62, 0xf6, 0x4e, 0xad,
	0x04, 0x3b, 0x98, 0x27, 0x08, 0xc4, 0x04, 0x59, 0x16, 0x96, 0x98, 0x53, 0x9a, 0x2a, 0xc1, 0x04,
	0xb1, 0x0c, 0xcc, 0x71, 0x92, 0x38, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f,
	0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x24,
	0x36, 0x70, 0xe0, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xb9, 0xc8, 0x0b, 0xdc, 0x2f, 0x01,
	0x00, 0x00,
}

func (m *LogEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LogEvent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LogEvent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RawSize != 0 {
		i = encodeVarintLogEvent(dAtA, i, uint64(m.RawSize))
		i--
		dAtA[i] = 0x28
	}
	if m.FileOffset != 0 {
		i = encodeVarintLogEvent(dAtA, i, uint64(m.FileOffset))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Level) > 0 {
		i -= len(m.Level)
		copy(dAtA[i:], m.Level)
		i = encodeVarintLogEvent(dAtA, i, uint64(len(m.Level)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Contents) > 0 {
		for iNdEx := len(m.Contents) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Contents[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintLogEvent(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Timestamp != 0 {
		i = encodeVarintLogEvent(dAtA, i, uint64(m.Timestamp))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *LogEvent_Content) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LogEvent_Content) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LogEvent_Content) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Value) > 0 {
		i -= len(m.Value)
		copy(dAtA[i:], m.Value)
		i = encodeVarintLogEvent(dAtA, i, uint64(len(m.Value)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintLogEvent(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintLogEvent(dAtA []byte, offset int, v uint64) int {
	offset -= sovLogEvent(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *LogEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Timestamp != 0 {
		n += 1 + sovLogEvent(uint64(m.Timestamp))
	}
	if len(m.Contents) > 0 {
		for _, e := range m.Contents {
			l = e.Size()
			n += 1 + l + sovLogEvent(uint64(l))
		}
	}
	l = len(m.Level)
	if l > 0 {
		n += 1 + l + sovLogEvent(uint64(l))
	}
	if m.FileOffset != 0 {
		n += 1 + sovLogEvent(uint64(m.FileOffset))
	}
	if m.RawSize != 0 {
		n += 1 + sovLogEvent(uint64(m.RawSize))
	}
	return n
}

func (m *LogEvent_Content) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovLogEvent(uint64(l))
	}
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovLogEvent(uint64(l))
	}
	return n
}

func sovLogEvent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozLogEvent(x uint64) (n int) {
	return sovLogEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *LogEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLogEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: LogEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LogEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLogEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Contents", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLogEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthLogEvent
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLogEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Contents = append(m.Contents, &LogEvent_Content{})
			if err := m.Contents[len(m.Contents)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Level", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLogEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthLogEvent
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthLogEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Level = append(m.Level[:0], dAtA[iNdEx:postIndex]...)
			if m.Level == nil {
				m.Level = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FileOffset", wireType)
			}
			m.FileOffset = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLogEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FileOffset |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RawSize", wireType)
			}
			m.RawSize = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLogEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RawSize |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipLogEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLogEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *LogEvent_Content) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLogEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Content: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Content: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLogEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthLogEvent
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthLogEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = append(m.Key[:0], dAtA[iNdEx:postIndex]...)
			if m.Key == nil {
				m.Key = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLogEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthLogEvent
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthLogEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = append(m.Value[:0], dAtA[iNdEx:postIndex]...)
			if m.Value == nil {
				m.Value = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLogEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLogEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipLogEvent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLogEvent
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowLogEvent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowLogEvent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthLogEvent
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupLogEvent
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthLogEvent
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthLogEvent        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLogEvent          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupLogEvent = fmt.Errorf("proto: unexpected end of group")
)
