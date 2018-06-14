package boltstore

import "errors"

var (
	// ErrWriteControlPacket ControlPacket's 'Write' function error
	ErrWriteControlPacket = errors.New("failed to convert control packet to byte array via io.writer")

	// ErrReadControlPacket failed to convert control packet's byte array to io.reader
	ErrReadControlPacket = errors.New("failed to convert control packet's byte array to io.reader")

	// ErrCreateBucket failed to create boltdb bucket
	ErrCreateBucket = errors.New("failed to create a boltdb bucket")

	// ErrGetBucket faild to get bucket
	ErrGetBucket = errors.New("failed to get the bucket")

	// ErrDelBucket failed to delete the default bucket
	ErrDelBucket = errors.New("failed to delete the default bucket, maybe empty?")

	// ErrPutValue failed to put key value pair
	ErrPutValue = errors.New("failed to put a key value pair")

	// ErrGetValue failed to get value from the bucket
	ErrGetValue = errors.New("failed to get vaule from the bucket")

	// ErrDelValue failed to delete a value from the bucket
	ErrDelValue = errors.New("failed to delete vaule from the bucket")

	// ErrOpenDB failed to open boltdb
	ErrOpenDB = errors.New("failed to open a boltdb")

	// ErrUseDB try to use the bolt store, but not open
	ErrUseDB = errors.New("try to use the boltdb, but not open")
)
