// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: deployer.proto

package deployer_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type DeploymentFilter struct {
	Cluster              string   `protobuf:"bytes,1,opt,name=cluster,proto3" json:"cluster,omitempty"`
	DeployerName         string   `protobuf:"bytes,2,opt,name=deployerName,proto3" json:"deployerName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeploymentFilter) Reset()         { *m = DeploymentFilter{} }
func (m *DeploymentFilter) String() string { return proto.CompactTextString(m) }
func (*DeploymentFilter) ProtoMessage()    {}
func (*DeploymentFilter) Descriptor() ([]byte, []int) {
	return fileDescriptor_deployer_50b68556edb7e7e5, []int{0}
}
func (m *DeploymentFilter) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DeploymentFilter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DeploymentFilter.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *DeploymentFilter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeploymentFilter.Merge(dst, src)
}
func (m *DeploymentFilter) XXX_Size() int {
	return m.Size()
}
func (m *DeploymentFilter) XXX_DiscardUnknown() {
	xxx_messageInfo_DeploymentFilter.DiscardUnknown(m)
}

var xxx_messageInfo_DeploymentFilter proto.InternalMessageInfo

func (m *DeploymentFilter) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func (m *DeploymentFilter) GetDeployerName() string {
	if m != nil {
		return m.DeployerName
	}
	return ""
}

type Deployment struct {
	Id                   string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Namespace            string            `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Name                 string            `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	ImageName            string            `protobuf:"bytes,4,opt,name=imageName,proto3" json:"imageName,omitempty"`
	IngressHost          string            `protobuf:"bytes,5,opt,name=ingressHost,proto3" json:"ingressHost,omitempty"`
	Labels               map[string]string `protobuf:"bytes,6,rep,name=labels" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Deployment) Reset()         { *m = Deployment{} }
func (m *Deployment) String() string { return proto.CompactTextString(m) }
func (*Deployment) ProtoMessage()    {}
func (*Deployment) Descriptor() ([]byte, []int) {
	return fileDescriptor_deployer_50b68556edb7e7e5, []int{1}
}
func (m *Deployment) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Deployment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Deployment.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Deployment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Deployment.Merge(dst, src)
}
func (m *Deployment) XXX_Size() int {
	return m.Size()
}
func (m *Deployment) XXX_DiscardUnknown() {
	xxx_messageInfo_Deployment.DiscardUnknown(m)
}

var xxx_messageInfo_Deployment proto.InternalMessageInfo

func (m *Deployment) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Deployment) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *Deployment) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Deployment) GetImageName() string {
	if m != nil {
		return m.ImageName
	}
	return ""
}

func (m *Deployment) GetIngressHost() string {
	if m != nil {
		return m.IngressHost
	}
	return ""
}

func (m *Deployment) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

type GetDeploymentsResponse struct {
	Deployments          []*Deployment `protobuf:"bytes,1,rep,name=deployments" json:"deployments,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *GetDeploymentsResponse) Reset()         { *m = GetDeploymentsResponse{} }
func (m *GetDeploymentsResponse) String() string { return proto.CompactTextString(m) }
func (*GetDeploymentsResponse) ProtoMessage()    {}
func (*GetDeploymentsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_deployer_50b68556edb7e7e5, []int{2}
}
func (m *GetDeploymentsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GetDeploymentsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GetDeploymentsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *GetDeploymentsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDeploymentsResponse.Merge(dst, src)
}
func (m *GetDeploymentsResponse) XXX_Size() int {
	return m.Size()
}
func (m *GetDeploymentsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDeploymentsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetDeploymentsResponse proto.InternalMessageInfo

func (m *GetDeploymentsResponse) GetDeployments() []*Deployment {
	if m != nil {
		return m.Deployments
	}
	return nil
}

type UpdateDeploymentRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Status               string   `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateDeploymentRequest) Reset()         { *m = UpdateDeploymentRequest{} }
func (m *UpdateDeploymentRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateDeploymentRequest) ProtoMessage()    {}
func (*UpdateDeploymentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_deployer_50b68556edb7e7e5, []int{3}
}
func (m *UpdateDeploymentRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UpdateDeploymentRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UpdateDeploymentRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *UpdateDeploymentRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateDeploymentRequest.Merge(dst, src)
}
func (m *UpdateDeploymentRequest) XXX_Size() int {
	return m.Size()
}
func (m *UpdateDeploymentRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateDeploymentRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateDeploymentRequest proto.InternalMessageInfo

func (m *UpdateDeploymentRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UpdateDeploymentRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UpdateDeploymentRequest) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type UpdateDeploymentResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateDeploymentResponse) Reset()         { *m = UpdateDeploymentResponse{} }
func (m *UpdateDeploymentResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateDeploymentResponse) ProtoMessage()    {}
func (*UpdateDeploymentResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_deployer_50b68556edb7e7e5, []int{4}
}
func (m *UpdateDeploymentResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UpdateDeploymentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UpdateDeploymentResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *UpdateDeploymentResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateDeploymentResponse.Merge(dst, src)
}
func (m *UpdateDeploymentResponse) XXX_Size() int {
	return m.Size()
}
func (m *UpdateDeploymentResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateDeploymentResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateDeploymentResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*DeploymentFilter)(nil), "deployer_v1.DeploymentFilter")
	proto.RegisterType((*Deployment)(nil), "deployer_v1.Deployment")
	proto.RegisterMapType((map[string]string)(nil), "deployer_v1.Deployment.LabelsEntry")
	proto.RegisterType((*GetDeploymentsResponse)(nil), "deployer_v1.GetDeploymentsResponse")
	proto.RegisterType((*UpdateDeploymentRequest)(nil), "deployer_v1.UpdateDeploymentRequest")
	proto.RegisterType((*UpdateDeploymentResponse)(nil), "deployer_v1.UpdateDeploymentResponse")
}
func (m *DeploymentFilter) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DeploymentFilter) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Cluster) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintDeployer(dAtA, i, uint64(len(m.Cluster)))
		i += copy(dAtA[i:], m.Cluster)
	}
	if len(m.DeployerName) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintDeployer(dAtA, i, uint64(len(m.DeployerName)))
		i += copy(dAtA[i:], m.DeployerName)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *Deployment) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Deployment) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Id) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintDeployer(dAtA, i, uint64(len(m.Id)))
		i += copy(dAtA[i:], m.Id)
	}
	if len(m.Namespace) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintDeployer(dAtA, i, uint64(len(m.Namespace)))
		i += copy(dAtA[i:], m.Namespace)
	}
	if len(m.Name) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintDeployer(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.ImageName) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintDeployer(dAtA, i, uint64(len(m.ImageName)))
		i += copy(dAtA[i:], m.ImageName)
	}
	if len(m.IngressHost) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintDeployer(dAtA, i, uint64(len(m.IngressHost)))
		i += copy(dAtA[i:], m.IngressHost)
	}
	if len(m.Labels) > 0 {
		for k, _ := range m.Labels {
			dAtA[i] = 0x32
			i++
			v := m.Labels[k]
			mapSize := 1 + len(k) + sovDeployer(uint64(len(k))) + 1 + len(v) + sovDeployer(uint64(len(v)))
			i = encodeVarintDeployer(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintDeployer(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintDeployer(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *GetDeploymentsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetDeploymentsResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Deployments) > 0 {
		for _, msg := range m.Deployments {
			dAtA[i] = 0xa
			i++
			i = encodeVarintDeployer(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *UpdateDeploymentRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpdateDeploymentRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Id) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintDeployer(dAtA, i, uint64(len(m.Id)))
		i += copy(dAtA[i:], m.Id)
	}
	if len(m.Name) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintDeployer(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.Status) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintDeployer(dAtA, i, uint64(len(m.Status)))
		i += copy(dAtA[i:], m.Status)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *UpdateDeploymentResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpdateDeploymentResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintDeployer(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *DeploymentFilter) Size() (n int) {
	var l int
	_ = l
	l = len(m.Cluster)
	if l > 0 {
		n += 1 + l + sovDeployer(uint64(l))
	}
	l = len(m.DeployerName)
	if l > 0 {
		n += 1 + l + sovDeployer(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Deployment) Size() (n int) {
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovDeployer(uint64(l))
	}
	l = len(m.Namespace)
	if l > 0 {
		n += 1 + l + sovDeployer(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovDeployer(uint64(l))
	}
	l = len(m.ImageName)
	if l > 0 {
		n += 1 + l + sovDeployer(uint64(l))
	}
	l = len(m.IngressHost)
	if l > 0 {
		n += 1 + l + sovDeployer(uint64(l))
	}
	if len(m.Labels) > 0 {
		for k, v := range m.Labels {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovDeployer(uint64(len(k))) + 1 + len(v) + sovDeployer(uint64(len(v)))
			n += mapEntrySize + 1 + sovDeployer(uint64(mapEntrySize))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *GetDeploymentsResponse) Size() (n int) {
	var l int
	_ = l
	if len(m.Deployments) > 0 {
		for _, e := range m.Deployments {
			l = e.Size()
			n += 1 + l + sovDeployer(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *UpdateDeploymentRequest) Size() (n int) {
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovDeployer(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovDeployer(uint64(l))
	}
	l = len(m.Status)
	if l > 0 {
		n += 1 + l + sovDeployer(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *UpdateDeploymentResponse) Size() (n int) {
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovDeployer(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozDeployer(x uint64) (n int) {
	return sovDeployer(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DeploymentFilter) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDeployer
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
			return fmt.Errorf("proto: DeploymentFilter: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DeploymentFilter: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cluster", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeployer
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
				return ErrInvalidLengthDeployer
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Cluster = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeployerName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeployer
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
				return ErrInvalidLengthDeployer
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DeployerName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDeployer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDeployer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Deployment) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDeployer
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
			return fmt.Errorf("proto: Deployment: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Deployment: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeployer
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
				return ErrInvalidLengthDeployer
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Namespace", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeployer
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
				return ErrInvalidLengthDeployer
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Namespace = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeployer
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
				return ErrInvalidLengthDeployer
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ImageName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeployer
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
				return ErrInvalidLengthDeployer
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ImageName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IngressHost", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeployer
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
				return ErrInvalidLengthDeployer
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IngressHost = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Labels", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeployer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDeployer
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Labels == nil {
				m.Labels = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowDeployer
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowDeployer
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthDeployer
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowDeployer
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthDeployer
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipDeployer(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthDeployer
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Labels[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDeployer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDeployer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GetDeploymentsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDeployer
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
			return fmt.Errorf("proto: GetDeploymentsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetDeploymentsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deployments", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeployer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDeployer
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Deployments = append(m.Deployments, &Deployment{})
			if err := m.Deployments[len(m.Deployments)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDeployer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDeployer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *UpdateDeploymentRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDeployer
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
			return fmt.Errorf("proto: UpdateDeploymentRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UpdateDeploymentRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeployer
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
				return ErrInvalidLengthDeployer
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeployer
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
				return ErrInvalidLengthDeployer
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDeployer
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
				return ErrInvalidLengthDeployer
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Status = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDeployer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDeployer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *UpdateDeploymentResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDeployer
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
			return fmt.Errorf("proto: UpdateDeploymentResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UpdateDeploymentResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipDeployer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDeployer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipDeployer(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDeployer
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
					return 0, ErrIntOverflowDeployer
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
					return 0, ErrIntOverflowDeployer
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
				return 0, ErrInvalidLengthDeployer
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowDeployer
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
				next, err := skipDeployer(dAtA[start:])
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
	ErrInvalidLengthDeployer = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDeployer   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("deployer.proto", fileDescriptor_deployer_50b68556edb7e7e5) }

var fileDescriptor_deployer_50b68556edb7e7e5 = []byte{
	// 414 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0x4b, 0x8e, 0xd3, 0x40,
	0x10, 0x4d, 0x3b, 0x33, 0x86, 0x29, 0xa3, 0xc1, 0x94, 0xd0, 0x4c, 0xcb, 0x02, 0x2b, 0x6a, 0x40,
	0x9a, 0x95, 0x05, 0xc3, 0x86, 0x81, 0x1d, 0x1a, 0x3e, 0x8b, 0x11, 0x20, 0x47, 0x11, 0x4b, 0xd4,
	0x89, 0x4b, 0x91, 0x85, 0x7f, 0xb8, 0xdb, 0x91, 0x72, 0x01, 0xae, 0x00, 0x47, 0x62, 0xc9, 0x11,
	0x50, 0xb8, 0x08, 0x4a, 0xdb, 0xc6, 0x0e, 0xc1, 0x42, 0xec, 0xfa, 0xbd, 0xaa, 0x57, 0xfd, 0xea,
	0xb5, 0x0d, 0xc7, 0x11, 0x15, 0x49, 0xbe, 0xa6, 0x32, 0x28, 0xca, 0x5c, 0xe7, 0xe8, 0xb4, 0xf8,
	0xc3, 0xea, 0x91, 0x78, 0x07, 0xee, 0xa5, 0x81, 0x29, 0x65, 0xfa, 0x65, 0x9c, 0x68, 0x2a, 0x91,
	0xc3, 0xb5, 0x45, 0x52, 0x29, 0x4d, 0x25, 0x67, 0x13, 0x76, 0x76, 0x14, 0xb6, 0x10, 0x05, 0xdc,
	0x68, 0xc5, 0x6f, 0x64, 0x4a, 0xdc, 0x32, 0xe5, 0x1d, 0x4e, 0x7c, 0xb6, 0x00, 0xba, 0x91, 0x78,
	0x0c, 0x56, 0x1c, 0x35, 0x73, 0xac, 0x38, 0xc2, 0x3b, 0x70, 0x94, 0xc9, 0x94, 0x54, 0x21, 0x17,
	0xad, 0xbe, 0x23, 0x10, 0xe1, 0x60, 0x0b, 0xf8, 0xd8, 0x14, 0xcc, 0x79, 0xab, 0x88, 0x53, 0xb9,
	0x24, 0x73, 0xe3, 0x41, 0xad, 0xf8, 0x4d, 0xe0, 0x04, 0x9c, 0x38, 0x5b, 0x96, 0xa4, 0xd4, 0xeb,
	0x5c, 0x69, 0x7e, 0x68, 0xea, 0x7d, 0x0a, 0x9f, 0x81, 0x9d, 0xc8, 0x39, 0x25, 0x8a, 0xdb, 0x93,
	0xf1, 0x99, 0x73, 0x7e, 0x2f, 0xe8, 0x05, 0x10, 0x74, 0x56, 0x83, 0x2b, 0xd3, 0xf5, 0x22, 0xd3,
	0xe5, 0x3a, 0x6c, 0x24, 0xde, 0x05, 0x38, 0x3d, 0x1a, 0x5d, 0x18, 0x7f, 0xa4, 0x75, 0xb3, 0xce,
	0xf6, 0x88, 0xb7, 0xe1, 0x70, 0x25, 0x93, 0xaa, 0xdd, 0xa5, 0x06, 0x4f, 0xad, 0x27, 0x4c, 0x4c,
	0xe1, 0xe4, 0x15, 0xe9, 0x6e, 0xbe, 0x0a, 0x49, 0x15, 0x79, 0xa6, 0x08, 0x2f, 0xa0, 0x79, 0x03,
	0x43, 0x73, 0x66, 0x6c, 0x9d, 0x0e, 0xd8, 0x0a, 0xfb, 0xbd, 0x62, 0x06, 0xa7, 0xb3, 0x22, 0x92,
	0x9a, 0x7a, 0x0d, 0xf4, 0xa9, 0x22, 0xb5, 0x9f, 0x74, 0x9b, 0xa5, 0xd5, 0xcb, 0xf2, 0x04, 0x6c,
	0xa5, 0xa5, 0xae, 0x54, 0x93, 0x70, 0x83, 0x84, 0x07, 0x7c, 0x7f, 0x6c, 0xed, 0xf6, 0xfc, 0x8b,
	0x05, 0xd7, 0x2f, 0x1b, 0x6b, 0xf8, 0x16, 0x6e, 0x4d, 0x75, 0x49, 0x32, 0xed, 0xed, 0x85, 0x77,
	0x07, 0xac, 0xd7, 0xdf, 0x93, 0x37, 0xb4, 0x99, 0x18, 0x3d, 0x64, 0xf8, 0x1e, 0x6e, 0x5e, 0xc5,
	0x4a, 0xff, 0xc7, 0xb8, 0xdd, 0xf7, 0xfb, 0x7b, 0xc4, 0x62, 0x84, 0x12, 0xdc, 0x3f, 0x57, 0xc2,
	0xfb, 0x3b, 0xd2, 0x81, 0x20, 0xbd, 0x07, 0xff, 0xe8, 0x6a, 0xaf, 0x78, 0xee, 0x7e, 0xdb, 0xf8,
	0xec, 0xfb, 0xc6, 0x67, 0x3f, 0x36, 0x3e, 0xfb, 0xfa, 0xd3, 0x1f, 0xcd, 0x6d, 0xf3, 0x8b, 0x3d,
	0xfe, 0x15, 0x00, 0x00, 0xff, 0xff, 0x92, 0x5c, 0xd9, 0x78, 0x74, 0x03, 0x00, 0x00,
}
