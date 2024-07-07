package main

import "fmt"

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

	return fmt.Errorf("a key value pair with this key already exists")
}
