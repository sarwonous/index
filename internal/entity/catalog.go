package entity

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gosimple/slug"
)

type Catalog struct {
	ID         int64     `json:"ID" db:"id"`
	ObjectID   string    `json:"object_id" db:"object_id"`
	Name       string    `json:"name" db:"name"`
	PathName   string    `json:"path_name" db:"path_name"`
	Files      []*File   `json:"files"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
	Status     string    `json:"status" db:"status"`
}

func (c *Catalog) GetObjectID() error {
	s := slug.Make(c.Name)
	c.ObjectID = s
	return nil
}

func (c *Catalog) AddFile(f *File, save bool) error {
	c.Files = append(c.Files, f)
	return nil
}

func (c *Catalog) Add(filepath string, save bool) error {
	infofile, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	file := &File{
		Name:       infofile.Name(),
		PathName:   filepath,
		CreatedAt:  infofile.ModTime(),
		ModifiedAt: infofile.ModTime(),
	}
	c.Files = append(c.Files, file)
	return nil
}

func (c *Catalog) ToJSON() (string, error) {
	j, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func NewCatalog(path string, name string) *Catalog {
	c := &Catalog{
		Name:      name,
		PathName:  path,
		CreatedAt: time.Now(),
		ObjectID:  slug.Make(name),
	}
	return c
}
