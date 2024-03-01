package data_store

type Note struct {
	Id               string `json:"id"`
	Body             string `json:"body"`
	TimeCreated      int64  `json:"timeCreated"`
	TimeLastModified int64  `json:"timeLastModified"`
}

type CreateNoteRequest struct {
	Body string `json:"body"`
}

type UpdateNoteRequest struct {
	Id   string `json:"id"`
	Body string `json:"body"`
}

type DataStore interface {
	DeleteNote(id string) error
	CreateNote(request *CreateNoteRequest) (*Note, error)
	UpdateNote(request *UpdateNoteRequest) (*Note, error)
	GetNotes() ([]*Note, error)
}

func New() DataStore {
	return newBleveDataStore()
}
