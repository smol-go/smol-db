package main

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

func createDB(filename string) SmolDb {
	return SmolDb{
		Filename: filename,
		KeyPairs: []KeyPair{},
	}
}
