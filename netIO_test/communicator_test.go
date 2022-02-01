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

		result, err2 := c.ReceiveVarInt()

		if input != result || err1 != nil || err2 != nil {
			t.Errorf("Send and receive var int failed: expected %d, result is %d", input, result)
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

		resultLen, err2 := c.ReceiveVarInt()
		resultBuf := make([]byte, resultLen)
		b.Read(resultBuf)

		if int(resultLen) != len(testCase) || !bytes.Equal(resultBuf, []byte(testCase)) ||
			err1 != nil || err2 != nil {
			t.Errorf("Send string failed: expected %s, result is %s", testCase, string(resultBuf))
		}
	}
}

func TestReceiveString(t *testing.T) {
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

		result, err2 := c.ReceiveString()

		if len(testCase) != len(result) || testCase != result ||
			err1 != nil || err2 != nil {
			t.Errorf("Receive string failed: expected %s, result is %s", testCase, result)
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

		sliceLen, err := c.ReceiveVarInt()
		resultSlice := make([]string, sliceLen)
		for i := 0; i < int(sliceLen) && err == nil; i++ {
			resultSlice[i], err = c.ReceiveString()
		}

		if int(sliceLen) != len(testCase) || !reflect.DeepEqual(testCase, resultSlice) ||
			err != nil || err1 != nil {
			t.Errorf("Send string slice failed: expected %s, result is %s", testCase, resultSlice)
		}
	}
}

func TestReceiveStringSlice(t *testing.T) {
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

		slice, err1 := c.ReceiveStringSlice()
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

		receivedLen, err1 := c.ReceiveVarInt()
		bytes := make([]byte, receivedLen)
		b.Read(bytes)

		if int(receivedLen) != len(testCase) || !reflect.DeepEqual(testCase, bytes) ||
			err != nil || err1 != nil {
			t.Errorf("Send file data failed: expected %s, result is %s", testCase, string(bytes))
		}
	}
}

func TestReceiveFileData(t *testing.T) {
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
		err := c.ReceiveFileData(output)

		bytes := make([]byte, len(testCase))
		output.Read(bytes)
		if !reflect.DeepEqual(testCase, bytes) ||
			err != nil || err1 != nil {
			t.Errorf("Send file data failed: expected \"%s\", result is \"%s\"", testCase, string(bytes))
		}
	}
}
