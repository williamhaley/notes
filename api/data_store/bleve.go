package data_store

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/google/uuid"
)

const INDEX_DB_FILE_PATH = "/data/notes.bleve"

type BleveDataStore struct {
	index bleve.Index
}

func newBleveDataStore() DataStore {
	mapping := bleve.NewIndexMapping()

	var index bleve.Index

	if _, err := os.Stat(INDEX_DB_FILE_PATH); errors.Is(err, os.ErrNotExist) {
		index, err = bleve.New(INDEX_DB_FILE_PATH, mapping)
		if err != nil {
			panic(err)
		}
	} else {
		index, err = bleve.Open(INDEX_DB_FILE_PATH)
		if err != nil {
			panic(err)
		}
	}

	return &BleveDataStore{
		index: index,
	}
}

func (b *BleveDataStore) DeleteNote(id string) error {
	return b.index.Delete(id)
}

func (b *BleveDataStore) CreateNote(request *CreateNoteRequest) (*Note, error) {
	now := time.Now().Unix()
	id := uuid.NewString()

	note := &Note{
		Id:               id,
		TimeCreated:      now,
		TimeLastModified: now,
		Body:             request.Body,
	}

	if err := b.index.Index(id, note); err != nil {
		return nil, err
	}
	return note, nil
}

func (b *BleveDataStore) UpdateNote(request *UpdateNoteRequest) (*Note, error) {
	searchRequest := bleve.NewSearchRequest(bleve.NewDocIDQuery([]string{request.Id}))
	searchRequest.Fields = []string{"*"}
	searchRequest.Size = 2
	searchResult, _ := b.index.Search(searchRequest)

	now := time.Now().Unix()

	if searchResult.Total != 1 {
		return nil, fmt.Errorf("note with id %s not found (found at least %d records)", request.Id, searchResult.Total)
	}

	var timeCreated int64 = 0
	timeCreatedValue := searchResult.Hits[0].Fields["timeCreated"]
	if timeCreatedValue != nil {
		timeCreated = int64(timeCreatedValue.(float64))
	}

	note := &Note{
		Id:               request.Id,
		Body:             request.Body,
		TimeCreated:      timeCreated,
		TimeLastModified: now,
	}

	if err := b.index.Index(request.Id, note); err != nil {
		return nil, err
	}
	return note, nil
}

func (b *BleveDataStore) GetNotes() ([]*Note, error) {
	searchRequest := bleve.NewSearchRequest(bleve.NewMatchAllQuery())
	searchRequest.Fields = []string{"*"}
	// TODO paginate
	searchRequest.Size = 100
	searchRequest.SortBy([]string{"-timeLastModified"})
	searchResult, _ := b.index.Search(searchRequest)

	notes := make([]*Note, len(searchResult.Hits))

	for index, hit := range searchResult.Hits {
		notes[index] = unserializeNote(hit)
	}

	return notes, nil
}

func unserializeNote(hit *search.DocumentMatch) *Note {
	var body string
	if hit.Fields["body"] != nil {
		body = hit.Fields["body"].(string)
	}

	var timeCreated float64
	if hit.Fields["timeCreated"] != nil {
		timeCreated = hit.Fields["timeCreated"].(float64)
	}

	var timeLastModified float64
	if hit.Fields["timeLastModified"] != nil {
		timeLastModified = hit.Fields["timeLastModified"].(float64)
	}

	return &Note{
		Id:               hit.ID,
		Body:             body,
		TimeCreated:      int64(timeCreated),
		TimeLastModified: int64(timeLastModified),
	}
}
