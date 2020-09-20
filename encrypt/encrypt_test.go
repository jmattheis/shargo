package encrypt

import (
	"github.com/jmattheis/shargo/proto"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_success(t *testing.T) {
	start := &proto.Packet{Control: proto.ControlHello, Payload: []byte("hello world")}
	enc, err := Encrypt(Sha256("secure"), start)
	require.NoError(t, err)
	actual, err := Decrypt(Sha256("secure"), enc)
	require.NoError(t, err)
	require.Equal(t, start, actual)
}

func Test_Badpw(t *testing.T) {
	start := &proto.Packet{Control: proto.ControlHello, Payload: []byte("hello world")}
	enc, err := Encrypt(Sha256("secure"), start)
	require.NoError(t, err)
	_, err = Decrypt(Sha256("other"), enc)
	require.NotNil(t, err)
}
