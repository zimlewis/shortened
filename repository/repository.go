package repository

import (
	"context"
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
)

type Repository struct {
	options badger.Options
}


func New(options badger.Options) (*Repository) {
	return &Repository{
		options: options,
	}
}

func (repo *Repository) AddShortenedLink(ctx context.Context, shortened string, full string) error {
	db, err := badger.Open(repo.options)
	if err != nil {
		return fmt.Errorf("Cannot open database %w", err)
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set(fmt.Appendf(nil, "%s:full", shortened), []byte(full)) 
		return err
	})

	if err != nil {
		return fmt.Errorf("Cannot update badger: %w", err)
	}

	return nil
}

func (repo *Repository) GetShortenedResult(ctx context.Context, shortened string) (string, error) {
	db, err := badger.Open(repo.options)
	if err != nil {
		return "", fmt.Errorf("Cannot open database %w", err)
	}
	defer db.Close()

	var fullLink []byte

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(fmt.Appendf(nil, "%s:full", shortened))
		if err == badger.ErrKeyNotFound {
			return err
		}

		fullLink, err = item.ValueCopy(nil)

		return err
	})
	
	if err != nil {
		return "", err
	}

	return string(fullLink), nil
}




