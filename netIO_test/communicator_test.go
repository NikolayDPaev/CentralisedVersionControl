package netio_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

func TestSendVarInt(t *testing.T) {
	b := new(bytes.Buffer)
	c := netio.NewCommunicator(100, b, b)

	values := []int64{
		1,
		-1,
		0,
		125,
		555555,
		-3242,
	}

	for _, input := range values {
		err1 := c.SendVarInt(input)

		result, err2 := c.RecvVarInt()

		if input != result || err1 != nil || err2 != nil {
			t.Errorf("Send and recv var int failed: expected %d, result is %d", input, result)
		}
	}
}

func TestSendString(t *testing.T) {
	b := new(bytes.Buffer)
	c := netio.NewCommunicator(100, b, b)

	values := []string{
		"",
		"medium sized string",
		"bigger string with \n new lines",
	}

	for _, testCase := range values {
		err1 := c.SendString(testCase)

		resultLen, err2 := c.RecvVarInt()
		resultBuf := make([]byte, resultLen)
		b.Read(resultBuf)

		if int(resultLen) != len(testCase) || !bytes.Equal(resultBuf, []byte(testCase)) ||
			err1 != nil || err2 != nil {
			t.Errorf("Send string failed: expected %s, result is %s", testCase, string(resultBuf))
		}
	}
}

func TestRecvString(t *testing.T) {
	b := new(bytes.Buffer)
	c := netio.NewCommunicator(100, b, b)

	values := []string{
		"",
		"medium sized string",
		"bigger string with \n new lines",
	}

	for _, testCase := range values {
		err1 := c.SendVarInt(int64(len(testCase)))
		b.Write([]byte(testCase))

		result, err2 := c.RecvString()

		if len(testCase) != len(result) || testCase != result ||
			err1 != nil || err2 != nil {
			t.Errorf("Recv string failed: expected %s, result is %s", testCase, result)
		}
	}
}

func TestSendStringSlice(t *testing.T) {
	b := new(bytes.Buffer)
	c := netio.NewCommunicator(100, b, b)

	values := [][]string{
		{},
		{"", "string", "different string"},
	}

	for _, testCase := range values {
		err1 := c.SendStringSlice(testCase)

		sliceLen, err := c.RecvVarInt()
		resultSlice := make([]string, sliceLen)
		for i := 0; i < int(sliceLen) && err == nil; i++ {
			resultSlice[i], err = c.RecvString()
		}

		if int(sliceLen) != len(testCase) || !reflect.DeepEqual(testCase, resultSlice) ||
			err != nil || err1 != nil {
			t.Errorf("Send string slice failed: expected %s, result is %s", testCase, resultSlice)
		}
	}
}

func TestRecvStringSlice(t *testing.T) {
	b := new(bytes.Buffer)
	c := netio.NewCommunicator(100, b, b)

	values := [][]string{
		{},
		{"", "string", "different string"},
	}

	for _, testCase := range values {
		err := c.SendVarInt(int64(len(testCase)))
		for i := 0; i < len(testCase) && err == nil; i++ {
			err = c.SendString(testCase[i])
		}

		slice, err1 := c.RecvStringSlice()
		if len(slice) != len(testCase) || !reflect.DeepEqual(testCase, slice) ||
			err != nil || err1 != nil {
			t.Errorf("Send string slice failed: expected %s, result is %s", testCase, slice)
		}
	}
}

func TestSendFileData(t *testing.T) {
	b := new(bytes.Buffer)
	c := netio.NewCommunicator(5, b, b)

	values := [][]byte{
		[]byte(""),
		[]byte("something larger than 5 bytes"),
	}

	for _, testCase := range values {
		input := new(bytes.Buffer)
		input.Write(testCase)

		err := c.SendFileData(input, int64(len(testCase)))

		recvdLen, err1 := c.RecvVarInt()
		bytes := make([]byte, recvdLen)
		b.Read(bytes)

		if int(recvdLen) != len(testCase) || !reflect.DeepEqual(testCase, bytes) ||
			err != nil || err1 != nil {
			t.Errorf("Send file data failed: expected %s, result is %s", testCase, string(bytes))
		}
	}
}

func TestRecvFileData(t *testing.T) {
	b := new(bytes.Buffer)
	c := netio.NewCommunicator(5, b, b)

	values := [][]byte{
		[]byte(""),
		[]byte("something larger than 5 bytes"),
	}

	for _, testCase := range values {
		err1 := c.SendVarInt(int64(len(testCase)))
		b.Write(testCase)

		output := new(bytes.Buffer)
		err := c.RecvFileData(output)

		bytes := make([]byte, len(testCase))
		output.Read(bytes)
		if !reflect.DeepEqual(testCase, bytes) ||
			err != nil || err1 != nil {
			t.Errorf("Send file data failed: expected \"%s\", result is \"%s\"", testCase, string(bytes))
		}
	}
}
