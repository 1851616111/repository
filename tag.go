package main

func (db *DB) setTag(t *Tag) error {
	_, err := db.InsertOne(t)
	return err
}
