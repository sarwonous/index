package storage

import (
	"context"
	"database/sql"
	"fmt"
	"gcat/internal/entity"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Base struct {
	*sqlx.DB
	Name string
	Path string
}

var db *Base

func (b *Base) Setup() error {
	ctx := context.Background()
	sqltx, err := b.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	createCatalogQuery, err := sqltx.Prepare(`
	CREATE TABLE IF NOT EXISTS "Catalog" (
		ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		object_id VARCHAR(500) UNIQUE NOT NULL,
		path_name VARCHAR(500) NOT NULL,
		created_at DATETIME,
		modified_at DATETIME,
		status VARCHAR(10)
	);`)
	if err != nil {
		fmt.Println("Setup: createCatalogQuery.Prepare", err)
		return err
	}
	_, err = createCatalogQuery.Exec()
	defer createCatalogQuery.Close()
	if err != nil {
		fmt.Println("Setup: createCatalogQuery.Exec", err)
		return err
	}
	createFileQuery, err := sqltx.Prepare(`
	CREATE TABLE IF NOT EXISTS "File" (
		ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		object_id VARCHAR(500) NOT NULL,
		path_name VARCHAR(500) NOT NULL,
		absolute_path VARCHAR(500),
		created_at DATETIME,
		modified_at DATETIME,
		type VARCHAR(100),
		size BIGINT,
		catalog_id INTEGER NOT NULL,
		status VARCHAR(10)
	);`)
	if err != nil {
		fmt.Println("Setup: createFileQuery.Prepare", err)
		return err
	}
	_, err = createFileQuery.Exec()
	defer createFileQuery.Close()
	if err != nil {
		fmt.Println("Setup: createFileQuery.Exec", err)
		return err
	}
	err = sqltx.Commit()
	if err != nil {
		fmt.Println("Setup: Commit", err)
		return err
	}
	return nil
}

func (b *Base) IsCatalogExist(c *entity.Catalog) (*entity.Catalog, error) {
	checkQuery, err := b.Prepare(`
		SELECT ID, object_id, name, path_name, created_at, modified_at
		FROM Catalog
		WHERE object_id = ?
	`)
	if err != nil {
		fmt.Println("IsCatalogExist: checkQuery.Prepare", err)
		return c, err
	}

	row := checkQuery.QueryRow(c.ObjectID)
	defer checkQuery.Close()
	err = row.Scan(&c.ID, &c.ObjectID, &c.Name, &c.PathName, &c.CreatedAt, &c.ModifiedAt)
	if err != nil {
		fmt.Println("IsCatalogExist: checkQuery.Scan", err)
		return c, nil
	}
	return c, nil
}

func (b *Base) CreateCatalog(c *entity.Catalog) (*entity.Catalog, error) {
	query, err := b.PrepareNamed(`
		INSERT INTO Catalog (
			object_id,
			name,
			path_name,
			created_at,
			modified_at
		)
		VALUES (
			:object_id,
			:name,
			:path_name,
			:created_at,
			:modified_at
		)
	`)
	if err != nil {
		fmt.Println("CreateCatalog: PrepareNamed", err)
		return nil, err
	}
	result, err := query.Exec(c)
	defer query.Close()
	if err != nil {
		fmt.Println("CreateCatalog: Exec", err)
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("CreateCatalog: result.LastInsertId", err)
		return nil, err
	}
	c.ID = id
	return c, nil
}

func (b *Base) SaveCatalog(c *entity.Catalog) (*entity.Catalog, error) {
	c, err := b.IsCatalogExist(c)
	if err != nil {
		fmt.Println("SaveCatalog: IsCatalogExist", err)
		return c, err
	}
	if c.ID == 0 {
		_, err := b.CreateCatalog(c)
		if err != nil {
			fmt.Println("SaveCatalog: CreateCatalog", err)
			return nil, err
		}
	}
	query, err := b.PrepareNamed(`
	UPDATE Catalog
	SET
		name = :name,
		path_name = :path_name,
		created_at = :created_at,
		modified_at = :modified_at
	WHERE
		object_id = :object_id
	`)
	if err != nil {
		fmt.Println("SaveCatalog: query.PrepareNamed", err)
		return nil, err
	}

	_, err = query.Exec(c)
	defer query.Close()

	if err != nil {
		fmt.Println("SaveCatalog: query.Exec", err)
		return nil, err
	}
	return c, nil
}

func (b *Base) IsFileExist(f *entity.File) (bool, error) {
	var count int
	checkQuery, err := b.Prepare(`
		SELECT count(*)
		FROM File
		WHERE object_id = ?
		AND catalog_id = ?
	`)
	if err != nil {
		return false, err
	}
	row := checkQuery.QueryRow(f.ObjectID, f.CatalogID)
	defer checkQuery.Close()
	err = row.Scan(&count)
	if err != nil {
		return false, nil
	}
	return count > 0, nil
}

func (b *Base) SaveFile(c *entity.Catalog, f *entity.File, l *sync.Mutex) (*entity.File, error) {
	l.Lock()
	var query *sqlx.NamedStmt
	var err error
	exist, err := b.IsFileExist(f)
	if err != nil {
		l.Unlock()
		fmt.Println("SaveFile: IsFileExist error", err)
		return f, err
	}
	if exist {
		query, err = b.PrepareNamed(`
			UPDATE File
			SET
				name = :name,
				path_name = :path_name,
				created_at = :created_at,
				modified_at = :modified_at,
				size = :size,
				type = :type,
				catalog_id = :catalog_id
			WHERE
				object_id = :object_id
		`)
		if err != nil {
			l.Unlock()
			fmt.Println("SaveFile: PrepareNamed Update error", err)
			return nil, err
		}
	} else {
		query, err = b.PrepareNamed(`
			INSERT INTO File (
				object_id,
				name,
				path_name,
				absolute_path,
				created_at,
				modified_at,
				type,
				size,
				catalog_id
			)
			VALUES (
				:object_id,
				:name,
				:path_name,
				:absolute_path,
				:created_at,
				:modified_at,
				:type,
				:size,
				:catalog_id
			)
		`)
		if err != nil {
			l.Unlock()
			fmt.Println("SaveFile: PrepareNamed Insert error", err)
			return nil, err
		}
	}

	_, err = query.Exec(f)
	defer query.Close()
	if err != nil {
		l.Unlock()
		fmt.Println("SaveFile: Exec error", err)
		return nil, err
	}
	l.Unlock()
	return f, nil
}

func (b *Base) Open() error {
	db, err := sqlx.Open("sqlite3", fmt.Sprintf("%s?%s", b.Path, "_locking_mode=exclusive"))
	if err != nil {
		return err
	}
	b.DB = db
	return nil
}

func NewDB(name string, path string) (*Base, error) {
	if db != nil && db.Name == name {
		return db, nil
	}
	b := &Base{
		Name: name,
		Path: path,
	}

	err := b.Open()
	if err != nil {
		return nil, err
	}

	return b, nil
}
