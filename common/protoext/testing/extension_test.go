package testing

import (
	"github.com/stretchr/testify/assert"
	"github.com/v2fly/v2ray-core/v4/common/protoext"
	"google.golang.org/protobuf/reflect/protoreflect"
	"testing"
)

func TestMessageOpt(t *testing.T) {
	msg := TestingMessage{}
	opt, err := protoext.GetMessageOptions(msg.ProtoReflect().Descriptor())
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"demo", "demo2"}, opt.Type)
}

func TestFieldOpt(t *testing.T) {
	msg := TestingMessage{
		TestField: "Test",
	}
	msgreflect := msg.ProtoReflect()
	msgreflect.Range(func(descriptor protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		opt, err := protoext.GetFieldOptions(descriptor)
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"test", "test2"}, opt.AllowedValues)
		return true
	})
}
