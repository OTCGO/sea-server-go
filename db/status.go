package db

import "fmt"

func InsertStatus(status *Status) error {
	if status == nil {
		return fmt.Errorf("status is nil")
	}
	if status.Name == "" {
		
		return fmt.Errorf("status.name is empty")
	}

	_, err := db.engine.Insert(status)
	if err != nil {
		return err
	}
	return nil
}

func MustUpdateStatus(status Status) error {
	up, err := db.engine.Cols("update_height").Update(&status, &Status{Name: status.Name})
	if up == 0 {
		return fmt.Errorf("not found")
	}
	return err
}

func FindAllStatus() ([]*Status, error) {
	var s []*Status
	err := db.engine.Find(&s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func GetStatus(name string) (*Status, error) {
	var status Status
	ok, err := db.engine.Where("name = ?", name).Get(&status)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return &status, nil
}

func InsertBlock(block *Block) (bool, error) {
	if block == nil {
		return false, fmt.Errorf("block is nil")
	}

	exists, err := db.engine.Exist(&Block{Height: block.Height})
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	effected, err := db.engine.Insert(block)
	if err != nil {
		return false, err
	}

	return effected == 1, nil
}

func GetBlock(height int) (*Block, error) {
	var block Block
	ok, err := db.engine.Where("height = ?", height).Get(&block)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return &block, nil
}