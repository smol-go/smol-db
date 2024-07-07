package main

import (
	"compress/zlib"
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

type KeyPair struct {
	Key  string
	Pair interface{}
}

type SmolDb struct {
	Filename          string
	KeyPairs          []KeyPair
	Compress          bool
	ZlibCompressLevel int
}

func newKeyPair(key string, pair interface{}) KeyPair {
	return KeyPair{
		Key:  key,
		Pair: pair,
	}
}

func (smoldb *SmolDb) keyExists(key string) int {
	for idx, kp := range smoldb.KeyPairs {
		if kp.Key == key {
			return idx
		}
	}

	return -1
}

func createDB(filename string) SmolDb {
	return SmolDb{
		Filename: filename,
		KeyPairs: []KeyPair{},
	}
}

func (smoldb *SmolDb) clear() {
	smoldb.KeyPairs = []KeyPair{}
}

func (smoldb *SmolDb) compress() {
	smoldb.Compress = true
}

func (smolDb *SmolDb) add(key string, pair interface{}) error {
	if smolDb.keyExists(key) != -1 {
		return fmt.Errorf("key with this name already exists")
	}

	smolDb.KeyPairs = append(smolDb.KeyPairs, newKeyPair(key, pair))

	return nil
}

func (smoldb *SmolDb) get(key string) (interface{}, error) {
	if idx := smoldb.keyExists(key); idx != -1 {
		return smoldb.KeyPairs[idx], nil
	}

	return nil, fmt.Errorf("key with name %s not found", key)
}

func (smoldb *SmolDb) set(key string, value interface{}) error {
	if idx := smoldb.keyExists(key); idx != -1 {
		smoldb.KeyPairs[idx].Pair = value
		return nil
	}

	return fmt.Errorf("key value pair with this key already exists")
}

func (smoldb *SmolDb) delete(key string) error {
	if idx := smoldb.keyExists(key); idx != 1 {
		smoldb.KeyPairs = append(smoldb.KeyPairs[:idx], smoldb.KeyPairs[:idx+1]...)

		return nil
	}

	return fmt.Errorf("key value pair with this key does not exists")
}

func (smoldb *SmolDb) representate() map[string]interface{} {
	var repr map[string]interface{} = map[string]interface{}{}

	for _, kp := range smoldb.KeyPairs {
		repr[kp.Key] = kp.Pair
	}

	return repr
}

func (smoldb *SmolDb) groupByDatatype(datatype interface{}) map[string]interface{} {
	data := map[string]interface{}{}

	for _, kp := range smoldb.KeyPairs {
		if fmt.Sprintf("%T", kp.Pair) == fmt.Sprintf("%T", datatype) {
			data[kp.Key] = kp.Pair
		}
	}

	return data
}

func (smoldb *SmolDb) save() error {
	buf, err := os.OpenFile(smoldb.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	if smoldb.Compress {
		// If compression is enabled, create a new zlib writer with the specified compression level
		var out io.Writer

		out, err := zlib.NewWriterLevel(buf, smoldb.ZlibCompressLevel)

		if err != nil {
			return err
		}

		// Creates a new gob encoder that writes to the zlib writer
		enc := gob.NewEncoder(out)
		err = enc.Encode(smoldb)

		if err != nil {
			return err
		}

		// Checks if the zlib writer (out) implements the io.Closer interface
		if c, ok := out.(io.Closer); ok {
			err = c.Close()

			if err != nil {
				return err
			}

			err = buf.Close()

			if err != nil {
				return err
			}
		}

		return nil
	}

	// If compression is not enabled, create a gob encoder that writes to the file
	enc := gob.NewEncoder(buf)
	enc.Encode(smoldb)

	return nil
}
