package coder

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type Coder interface {
	EncodeSegment(bytes []byte) []byte
}

type CoderImpl struct {
	errorProbability   int // in %
	missingProbability int
	logger             *logrus.Entry
}

func NewCoderImpl(errorProbability, missingProbability int, logger *logrus.Logger) *CoderImpl {
	rand.Seed(time.Now().UnixNano())

	return &CoderImpl{
		errorProbability:   errorProbability,
		missingProbability: missingProbability,
		logger:             logger.WithField("component", "coder"),
	}
}

func (c *CoderImpl) splitByte(dataByte byte) [2]byte {
	return [2]byte{dataByte >> 4, dataByte & 0x0F}
}

func (c *CoderImpl) encode(data byte) byte {
	return data ^ data<<1 ^ data<<3
}

func (c *CoderImpl) encodeWithErrorProbability(data byte) byte {
	encodedByte := c.encode(data)

	randomNumber := rand.Intn(100) // from 0 to 99

	if randomNumber < c.errorProbability {
		randomBitNumber := rand.Intn(7) // from 0 to 6
		c.logger.Infof("message got error at %d bit", randomBitNumber)

		encodedByte = encodedByte ^ (1 << randomBitNumber)
		c.logger.Info("error bit index : ", randomBitNumber)
	}

	return encodedByte
}

// EncodeSegment encodes segment with given error probability or may lose segment with given missing probability
func (c *CoderImpl) EncodeSegment(bytes []byte) []byte {
	random := rand.Intn(99)
	if random < c.missingProbability {
		c.logger.Info("message missed")
		return []byte{} // message missed with 1 % probability
	}

	encodedBatch := make([]byte, 0, 2*len(bytes))
	for _, b := range bytes {
		dataBlocks := c.splitByte(b) // get 2 data blocks, 4 bits each

		for _, dataBlock := range dataBlocks {
			encodedDataBlock := c.encodeWithErrorProbability(dataBlock)

			encodedBatch = append(encodedBatch, encodedDataBlock)
		}
	}

	return encodedBatch
}
