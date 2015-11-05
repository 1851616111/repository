package main

func (db *DB) getDataitemUsageByIds(dataItemIds ...interface{}) (map[int64]DataItemUsage, error) {
	m := map[int64]DataItemUsage{}
	err := db.In(" DATAITEM_ID ", dataItemIds).Find(&m)
	return m, err
}
