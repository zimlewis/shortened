package repository

import (
	"context"
	"encoding/binary"
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
)

type BadgerRepository struct {
	db     *badger.DB
}


func NewBadger(options *badger.Options) (*BadgerRepository, error) {
	db, err := badger.Open(*options)
	return &BadgerRepository{
		db: db,
	}, err
}

func (repo *BadgerRepository) Close() {
	repo.db.Close()
}

func (repo *BadgerRepository) AddShortenedLink(ctx context.Context, shortened string, full string) error {
	db := repo.db

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set(fmt.Appendf(nil, "%s:full", shortened), []byte(full)) 
		return err
	})

	if err != nil {
		return fmt.Errorf("Cannot update badger: %w", err)
	}

	return nil
}

func (repo *BadgerRepository) GetShortenedResult(ctx context.Context, shortened string) (string, error) {
	db := repo.db
	var fullLink []byte

	err := db.View(func(txn *badger.Txn) error {
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


func (repo *BadgerRepository) IncreaseLinkClick(ctx context.Context, shortened string) (int, error) {
	db := repo.db

	var totalClick uint32 = 0
	err := db.Update(func(txn *badger.Txn) error {
		bs := make([]byte, 4)
		var newCount uint32 = 0
		key := fmt.Appendf(nil, "%s:count", shortened)

		item, err := txn.Get(key)
		if err != nil && err != badger.ErrKeyNotFound {
			return fmt.Errorf("Error retrieving item: %w", err)
		}
		switch err {
		case badger.ErrKeyNotFound:
			newCount = 1
		case nil:
			count, err := item.ValueCopy(nil)
			if err != nil { return fmt.Errorf("Cannot copy value: %w", err) }
			newCount = binary.LittleEndian.Uint32(count) + 1
		}

		binary.LittleEndian.PutUint32(bs, newCount)

		err = txn.Set(key, bs)
		if err != nil {
			return fmt.Errorf("Cannot set the click value: %w", err)
		}

		totalClick = newCount
		return nil
	})
	if err != nil {
		return 0, err
	}

	return int(totalClick), nil
}

func (repo *BadgerRepository) GetClickedCount(ctx context.Context, shortened string) (uint32, error) {
	db := repo.db

	var totalClick uint32 = 0
	err := db.Update(func(txn *badger.Txn) error {
		key := fmt.Appendf(nil, "%s:count", shortened)
		buf := make([]byte, 4)

		item, err := txn.Get(key)
		if err != nil {
			return fmt.Errorf("Error retrieving item: %w", err)
		}

		buf, err = item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("Error copying value: %w", err)
		}

		totalClick = binary.LittleEndian.Uint32(buf)

		return nil
	})
	if err != nil {
		return 0, err
	}

	return totalClick, nil
}

