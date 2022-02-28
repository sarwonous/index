package entity

import (
	"encoding/json"
	"time"

	"github.com/gosimple/slug"
)

type File struct {
	ID           int64     `json:"ID" db:"id"`
	ObjectID     string    `json:"object_id" db:"object_id"`
	Name         string    `json:"name" db:"name"`
	PathName     string    `json:"path_name" db:"path_name"`
	AbsolutePath string    `json:"absolute_path" db:"absolute_path"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	ModifiedAt   time.Time `json:"modified_at" db:"modified_at"`
	Size         int64     `json:"size" db:"size"`
	Type         string    `json:"type" db:"type"`
	Status       string    `json:"status" db:"status"`
	CatalogID    int64     `json:"catalog_id" db:"catalog_id"`
}

func (f *File) GetObjectID() {
	if f.ObjectID == "" {
		s := slug.Make(f.PathName)
		f.ObjectID = s
	}
}

func (f *File) ToJSON() string {
	j, err := json.Marshal(f)
	if err != nil {
		return ""
	}
	return string(j)
}

func NewFile(path string, name string) *File {
	c := &File{
		Name:     name,
		PathName: path,
		ObjectID: slug.Make(path),
	}
	return c
}
