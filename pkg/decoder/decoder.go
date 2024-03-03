package decoder

import "github.com/sirupsen/logrus"

type Decoder interface {
	DecodeSegment(bytes []byte) []byte
}

type DecoderImpl struct {
	logger *logrus.Entry
}

var errorSyndrom = map[byte]int{
	0x01: 0,
	0x02: 1,
	0x04: 2,
	0x03: 3,
	0x06: 4,
	0x07: 5,
	0x05: 6,
}

func NewDecoderImpl(logger *logrus.Logger) *DecoderImpl {
	return &DecoderImpl{
		logger: logger.WithField("component", "decoder"),
	}
}

func (d *DecoderImpl) divide(dividend byte) (byte, byte) {
	divisor := byte(0x0B)
	quotient := byte(0)

	for i := 4; i >= 0; i-- {
		if dividend&(1<<(3+i)) != 0 {
			dividend ^= divisor << i
			quotient = quotient<<1 + 1
		} else {
			quotient = quotient << 1
		}
	}
	remainder := dividend

	return quotient, remainder
}

func (d *DecoderImpl) decodeByte(data byte) byte {
	quotient, remainder := d.divide(data)

	if remainder != 0 { // error
		d.logger.Info("error found")

		errorBitIndex := errorSyndrom[remainder]
		d.logger.Info("error bit index : ", errorBitIndex)

		data = data ^ (1 << errorBitIndex)

		quotient, remainder = d.divide(data)
	}

	return quotient
}

func (d *DecoderImpl) decodeTwoBytes(bytes [2]byte) byte {
	var result byte
	for i, b := range bytes {
		decodedByte := d.decodeByte(b)

		if i == 0 {
			result += decodedByte << 4
		} else {
			result += decodedByte
		}
	}

	return result
}

func (d *DecoderImpl) DecodeSegment(bytes []byte) []byte {
	decodedSegment := make([]byte, 0, len(bytes)/2)
	for i := 0; i < len(bytes); i += 2 {
		decodedSegment = append(decodedSegment, d.decodeTwoBytes([2]byte{bytes[i], bytes[i+1]}))
	}

	return decodedSegment
}
