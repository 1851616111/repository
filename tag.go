package main

func (db *DB) setTag(t *Tag) error {
	_, err := db.InsertOne(t)
	return err
}
func (db *DB) getTags(tags ...interface{}) ([]Tag, error) {
	l := []Tag{}

	if len(tags) == 0 {
		return l, nil
	}
	err := db.In("DATAITEM_ID", tags...).Find(&l)
	return l, err
}
