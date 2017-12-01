package protobufdemo

import (
	pb "github.com/avshabanov/go-code/serialization/protobufdemo/.gen/proto"
	"github.com/golang/protobuf/ptypes/wrappers"
)

// Foo is a sample function that makes use of protobuf entites
func Foo(s string) string {
	s2 := wrappers.StringValue{Value: s}
	return s2.Value + "1"
}

// NewBobProfile creates a new profile for user bob
func NewBobProfile() *pb.UserProfile {
	return &pb.UserProfile{
		Id:    "1",
		Name:  &pb.FullName{First: "bob"},
		Age:   25,
		State: pb.ProfileState_ACTIVE,
	}
}
