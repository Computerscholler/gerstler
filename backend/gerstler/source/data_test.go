package source_test

import (
	"testing"

	"github.com/computerscholler/gerstler/filesystem"
	"github.com/computerscholler/gerstler/source"
	"github.com/stretchr/testify/assert"
)

func TestAddEntry(t *testing.T) {
	txt := "Hi world, I am Adrian and he wants to learn Golang."
	src := source.DummyFile{Content: txt, Title: "Hello"}
	source.AddEntry(src)
	_, err := source.ReadRecord(source.Id("0"))
	assert.NoError(t, err)
}

func TestTranformSourceToData(t *testing.T) {
	txt := "Hi world, I am Adrian and he wants to learn Golang."
	src := source.DummyFile{Content: txt, Title: "Hello"}
	data := src.TransformToData()
	assert.Equal(t, "Hello", data.Title)
	assert.Equal(t, txt, data.Content)
}

func TestSaveData(t *testing.T) {
	data := source.DummyFile{Title: "Hello", Content: "Hola mundo"}
	source.AddEntry(data)
	file, close := filesystem.CreateTempFile(t, "")
	defer close()
	source.WriteRecordsToFile(file)

	uid := source.Id("0")
	_, err := source.ReadRecord(uid)
	assert.NoError(t, err)
	file.Seek(0, 0)
	_, err = source.ReadRecordFromFile(file, uid)
	assert.NoError(t, err)
}

func TestLoadDb(t *testing.T) {
	source.ResetMemoryDb()

	data := source.DummyFile{Title: "Hello", Content: "Hola mundo"}
	source.AddEntry(data)

	recFile, closeR := filesystem.CreateTempFile(t, "")
	defer closeR()
	idxFile, closeI := filesystem.CreateTempFile(t, "")
	defer closeI()
	source.GenerateIndexer()
	err := source.SaveDb(recFile, idxFile)
	assert.NoError(t, err)
	source.ResetMemoryDb()

	recFile.Seek(0, 0)
	idxFile.Seek(0, 0)
	err = source.LoadDb(recFile, idxFile)
	assert.NoError(t, err)

	records, _ := source.Search("mundo")
	assert.Equal(t, 1, len(records))
}
