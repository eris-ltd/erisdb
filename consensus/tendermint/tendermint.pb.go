// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: tendermint.proto

/*
	Package tendermint is a generated protocol buffer package.

	It is generated from these files:
		tendermint.proto

	It has these top-level messages:
		NodeInfo
*/
package tendermint

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/hyperledger/burrow/crypto"

import github_com_hyperledger_burrow_crypto "github.com/hyperledger/burrow/crypto"
import github_com_hyperledger_burrow_binary "github.com/hyperledger/burrow/binary"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type NodeInfo struct {
	ID            github_com_hyperledger_burrow_crypto.Address  `protobuf:"bytes,1,opt,name=ID,proto3,customtype=github.com/hyperledger/burrow/crypto.Address" json:"ID"`
	ListenAddress string                                        `protobuf:"bytes,2,opt,name=ListenAddress,proto3" json:"ListenAddress,omitempty"`
	Network       string                                        `protobuf:"bytes,3,opt,name=Network,proto3" json:"Network,omitempty"`
	Version       string                                        `protobuf:"bytes,4,opt,name=Version,proto3" json:"Version,omitempty"`
	Channels      github_com_hyperledger_burrow_binary.HexBytes `protobuf:"bytes,5,opt,name=Channels,proto3,customtype=github.com/hyperledger/burrow/binary.HexBytes" json:"Channels"`
	Moniker       string                                        `protobuf:"bytes,6,opt,name=Moniker,proto3" json:"Moniker,omitempty"`
	Other         []string                                      `protobuf:"bytes,7,rep,name=Other" json:"Other,omitempty"`
}

func (m *NodeInfo) Reset()                    { *m = NodeInfo{} }
func (m *NodeInfo) String() string            { return proto.CompactTextString(m) }
func (*NodeInfo) ProtoMessage()               {}
func (*NodeInfo) Descriptor() ([]byte, []int) { return fileDescriptorTendermint, []int{0} }

func (m *NodeInfo) GetListenAddress() string {
	if m != nil {
		return m.ListenAddress
	}
	return ""
}

func (m *NodeInfo) GetNetwork() string {
	if m != nil {
		return m.Network
	}
	return ""
}

func (m *NodeInfo) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *NodeInfo) GetMoniker() string {
	if m != nil {
		return m.Moniker
	}
	return ""
}

func (m *NodeInfo) GetOther() []string {
	if m != nil {
		return m.Other
	}
	return nil
}

func (*NodeInfo) XXX_MessageName() string {
	return "tendermint.NodeInfo"
}
func init() {
	proto.RegisterType((*NodeInfo)(nil), "tendermint.NodeInfo")
	golang_proto.RegisterType((*NodeInfo)(nil), "tendermint.NodeInfo")
}
func (m *NodeInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NodeInfo) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintTendermint(dAtA, i, uint64(m.ID.Size()))
	n1, err := m.ID.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if len(m.ListenAddress) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintTendermint(dAtA, i, uint64(len(m.ListenAddress)))
		i += copy(dAtA[i:], m.ListenAddress)
	}
	if len(m.Network) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintTendermint(dAtA, i, uint64(len(m.Network)))
		i += copy(dAtA[i:], m.Network)
	}
	if len(m.Version) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintTendermint(dAtA, i, uint64(len(m.Version)))
		i += copy(dAtA[i:], m.Version)
	}
	dAtA[i] = 0x2a
	i++
	i = encodeVarintTendermint(dAtA, i, uint64(m.Channels.Size()))
	n2, err := m.Channels.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	if len(m.Moniker) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintTendermint(dAtA, i, uint64(len(m.Moniker)))
		i += copy(dAtA[i:], m.Moniker)
	}
	if len(m.Other) > 0 {
		for _, s := range m.Other {
			dAtA[i] = 0x3a
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	return i, nil
}

func encodeVarintTendermint(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *NodeInfo) Size() (n int) {
	var l int
	_ = l
	l = m.ID.Size()
	n += 1 + l + sovTendermint(uint64(l))
	l = len(m.ListenAddress)
	if l > 0 {
		n += 1 + l + sovTendermint(uint64(l))
	}
	l = len(m.Network)
	if l > 0 {
		n += 1 + l + sovTendermint(uint64(l))
	}
	l = len(m.Version)
	if l > 0 {
		n += 1 + l + sovTendermint(uint64(l))
	}
	l = m.Channels.Size()
	n += 1 + l + sovTendermint(uint64(l))
	l = len(m.Moniker)
	if l > 0 {
		n += 1 + l + sovTendermint(uint64(l))
	}
	if len(m.Other) > 0 {
		for _, s := range m.Other {
			l = len(s)
			n += 1 + l + sovTendermint(uint64(l))
		}
	}
	return n
}

func sovTendermint(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozTendermint(x uint64) (n int) {
	return sovTendermint(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *NodeInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTendermint
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: NodeInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NodeInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTendermint
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTendermint
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ListenAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTendermint
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTendermint
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ListenAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Network", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTendermint
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTendermint
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Network = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTendermint
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTendermint
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Version = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Channels", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTendermint
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTendermint
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Channels.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Moniker", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTendermint
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTendermint
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Moniker = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Other", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTendermint
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTendermint
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Other = append(m.Other, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTendermint(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTendermint
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
func skipTendermint(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTendermint
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
					return 0, ErrIntOverflowTendermint
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTendermint
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthTendermint
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowTendermint
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipTendermint(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthTendermint = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTendermint   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("tendermint.proto", fileDescriptorTendermint) }
func init() { golang_proto.RegisterFile("tendermint.proto", fileDescriptorTendermint) }

var fileDescriptorTendermint = []byte{
	// 321 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0xbd, 0x4e, 0x3a, 0x41,
	0x14, 0xc5, 0xff, 0xb3, 0xfc, 0xf9, 0x9a, 0x60, 0x62, 0x36, 0x16, 0x13, 0x8a, 0x85, 0x18, 0x0b,
	0x0a, 0x61, 0x13, 0x3f, 0x1e, 0x40, 0xa4, 0x80, 0x44, 0x31, 0x6e, 0x61, 0x61, 0xc7, 0xb2, 0x97,
	0xdd, 0x0d, 0x30, 0x97, 0xdc, 0x99, 0x0d, 0xee, 0x43, 0xf9, 0x0e, 0x96, 0x94, 0xd6, 0x16, 0xc4,
	0xc0, 0x8b, 0x18, 0x66, 0x57, 0xd1, 0x46, 0xbb, 0xf9, 0x9d, 0x33, 0x73, 0xee, 0xc9, 0x5c, 0x7e,
	0xa8, 0x41, 0x06, 0x40, 0xf3, 0x58, 0xea, 0xce, 0x82, 0x50, 0xa3, 0xcd, 0xf7, 0x4a, 0xbd, 0x1d,
	0xc6, 0x3a, 0x4a, 0xfc, 0xce, 0x18, 0xe7, 0x6e, 0x88, 0x21, 0xba, 0xe6, 0x8a, 0x9f, 0x4c, 0x0c,
	0x19, 0x30, 0xa7, 0xec, 0x69, 0xbd, 0x36, 0xa6, 0x74, 0xa1, 0x73, 0x3a, 0x7e, 0xb6, 0x78, 0x65,
	0x88, 0x01, 0x0c, 0xe4, 0x04, 0xed, 0x1e, 0xb7, 0x06, 0x3d, 0xc1, 0x9a, 0xac, 0x55, 0xeb, 0x5e,
	0xac, 0xd6, 0x8d, 0x7f, 0x6f, 0xeb, 0xc6, 0xe9, 0xb7, 0xf4, 0x28, 0x5d, 0x00, 0xcd, 0x20, 0x08,
	0x81, 0x5c, 0x3f, 0x21, 0xc2, 0xa5, 0x9b, 0x87, 0x5d, 0x05, 0x01, 0x81, 0x52, 0x9e, 0x35, 0xe8,
	0xd9, 0x27, 0xfc, 0xe0, 0x26, 0x56, 0x1a, 0x64, 0x2e, 0x0a, 0xab, 0xc9, 0x5a, 0x55, 0xef, 0xa7,
	0x68, 0x0b, 0x5e, 0x1e, 0x82, 0x5e, 0x22, 0x4d, 0x45, 0xc1, 0xf8, 0x9f, 0xb8, 0x73, 0x1e, 0x80,
	0x54, 0x8c, 0x52, 0xfc, 0xcf, 0x9c, 0x1c, 0xed, 0x7b, 0x5e, 0xb9, 0x8e, 0x46, 0x52, 0xc2, 0x4c,
	0x89, 0xa2, 0x69, 0x79, 0x99, 0xb7, 0x6c, 0xff, 0xde, 0xd2, 0x8f, 0xe5, 0x88, 0xd2, 0x4e, 0x1f,
	0x9e, 0xba, 0xa9, 0x06, 0xe5, 0x7d, 0xc5, 0xec, 0x86, 0xdd, 0xa2, 0x8c, 0xa7, 0x40, 0xa2, 0x94,
	0x0d, 0xcb, 0xd1, 0x3e, 0xe2, 0xc5, 0x3b, 0x1d, 0x01, 0x89, 0x72, 0xb3, 0xd0, 0xaa, 0x7a, 0x19,
	0x74, 0xfb, 0xab, 0x8d, 0xc3, 0x5e, 0x37, 0x0e, 0x7b, 0xdf, 0x38, 0xec, 0x65, 0xeb, 0xb0, 0xd5,
	0xd6, 0x61, 0x8f, 0x67, 0x7f, 0x7c, 0x12, 0x4a, 0x05, 0x52, 0x25, 0xca, 0xdd, 0xaf, 0xcd, 0x2f,
	0x99, 0x05, 0x9c, 0x7f, 0x04, 0x00, 0x00, 0xff, 0xff, 0x29, 0x50, 0xa8, 0x1b, 0xdd, 0x01, 0x00,
	0x00,
}