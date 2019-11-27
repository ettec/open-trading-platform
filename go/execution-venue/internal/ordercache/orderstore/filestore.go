package orderstore

import (
	"fmt"
	"github.com/coronationstreet/open-trading-platform/execution-venue/pb"
	"github.com/golang/protobuf/proto"
	"os"
)




type FileStore struct {
	file *os.File
}

func NewFileStore(path string) (*FileStore, error) {
	result := FileStore{}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666 )
	if err != nil {
		return nil, fmt.Errorf("unable to create file store: %w", err)
	}

	result.file = file

	return &result, nil
}

func (fs *FileStore) Close() {
	fs.file.Close()
}

func (fs *FileStore) Write(order *pb.Order) error {
	bytes, err := proto.Marshal(order)
	if err != nil {
		return fmt.Errorf("unable to convert order %v to bytes: %w", order, err)
	}
	_, err = fs.file.Write(bytes)
	if err != nil {
		return fmt.Errorf("unable to write order %v bytes: %w", order, err)
	}
	return nil
}