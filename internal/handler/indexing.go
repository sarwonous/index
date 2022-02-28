package handler

import (
	"encoding/json"
	"fmt"
	"gcat/internal/entity"
	"gcat/internal/logger"
	"gcat/internal/storage"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type Status struct {
	Total   int `json:"total"`
	Current int `json:"current"`
}

var bar *progressbar.ProgressBar

// load new catalog from disk

func WalkDirHere(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		// file
		bar.Add(1)
		// time.Sleep(1 * time.Millisecond)
	}
	return nil
}

func Load(path string, name string) (*entity.Catalog, error) {
	lock := sync.Mutex{}
	// initialize database ./catalog.db or from flag
	stdlog := logger.New()
	db, err := storage.NewDB("catalog", "./catalog.db")

	if db.Ping() != nil {
		defer db.Close()
	}

	if err != nil {
		stdlog.Error(err)
		panic(err)
	}

	// if no table exists create one
	if err = db.Setup(); err != nil {
		stdlog.Error(err)
		panic(err)
	}

	catalog := entity.NewCatalog(path, name)

	catalog, err = db.SaveCatalog(catalog)
	if catalog.ID == 0 {
		stdlog.Error(err)
		panic(err)
	}
	stdlog.Info("Loading catalog", name)
	bar = progressbar.Default(-1, "indexing...")
	status := &Status{
		Total:   0,
		Current: 0,
	}
	err = filepath.WalkDir(path, func(absolute_path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			bar.Add(1)
			finfo, err := d.Info()
			if err != nil {
				return err
			}
			rpath := strings.Replace(absolute_path, path, "", -1)
			rpath, err = filepath.Rel("/", fmt.Sprintf("%s/%s", rpath, finfo.Name()))

			if err != nil {
				return err
			}

			file := entity.NewFile(rpath, finfo.Name())
			file.AbsolutePath = path
			file.CreatedAt = finfo.ModTime()
			file.ModifiedAt = finfo.ModTime()
			file.Size = finfo.Size()
			file.CatalogID = catalog.ID
			// file.Type = finfo.Sys().(os.FileInfo).Mode().String()
			catalog.AddFile(file, false)
			status.Current = status.Current + 1
			d, err := json.Marshal(status)
			if err != nil {
				err := os.WriteFile("./tmp.json", d, 0644)
				if err != nil {
					logger.Info("error", nil)
				}
			}

			stdlog.WithField("data", file).Info("indexing")
		}
		return nil
	})
	if err != nil {
		stdlog.Error(err)
		panic(err)
	}

	for _, file := range catalog.Files {
		file.CatalogID = catalog.ID
		go db.SaveFile(catalog, file, &lock)
	}
	return nil, nil
}

func Reload(name string) (*entity.Catalog, error) {
	// initialize database ./catalog.db or from flag
	// find catalog by name
	// if no catalog found return error
	// if catalog found
	// calculating number of files show current file reading
	// saving to disk show progress bar
	// show success message
	return nil, nil
}

func Delete(name string) error {
	// initialize database ./catalog.db or from flag
	// find catalog by name
	// if no catalog found return error
	// if catalog found
	// delete catalog
	return nil
}
