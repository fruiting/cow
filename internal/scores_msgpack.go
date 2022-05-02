package internal

import (
	"errors"
	"fmt"
	"gopkg.in/vmihailenco/msgpack.v2"
)

func (s *Score) DecodeMsgpack(d *msgpack.Decoder) error {
	var err error
	var length int
	if length, err = d.DecodeArrayLen(); err != nil {
		return errors.New("error on decode array len")
	}
	if length != 4 {
		return fmt.Errorf("array len doesn't match: %d", length)
	}
	if s.GameId, err = d.DecodeString(); err != nil {
		return err
	}
	if s.Name, err = d.DecodeString(); err != nil {
		return err
	}
	if s.Score, err = d.DecodeInt8(); err != nil {
		return err
	}
	if s.ExpiresAt, err = d.DecodeInt64(); err != nil {
		return err
	}

	return nil
}
