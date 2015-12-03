package main

import (
	"log"
	"testing"
)

func Test_buildTime(t *testing.T) {
	param := []string{"2013-01-02 15:04:05", "2014-01-02 15:04:05", "2015-05-26 15:04:05", "2015-10-28 15:04:05", "2015-11-28 15:04:05", "2015-11-30 15:04:05"}
	for i := 0; i < len(param); i++ {

		t.Log(param[i])
		t.Log(buildTime(param[i]))
	}
}

func Test_ifInLabel(t *testing.T) {
	label := map[string]interface{}{"sys": map[string]interface{}{"supply_style": "api"}}
	e := ifInLabel(label, "supply_style")
	log.Println(e.ErrToString())
}

func Test_chkPrice(t *testing.T) {

	price := []map[string]interface{}{
		map[string]interface{}{"times": 1000, "money": 5, "expire": DATAITEM_PRICE_EXPIRE},
		map[string]interface{}{"times": 10000, "money": 45, "expire": DATAITEM_PRICE_EXPIRE},
		map[string]interface{}{"times": 100000.00, "money": 400.00, "expire": DATAITEM_PRICE_EXPIRE},
	}
	log.Println(chkPrice(price, "api"))

	price2 := []map[string]interface{}{
		map[string]interface{}{"time": 1000, "unit": "h", "money": 5, "expire": DATAITEM_PRICE_EXPIRE},
		map[string]interface{}{"time": 10000, "unit": "h", "money": 45, "expire": DATAITEM_PRICE_EXPIRE},
		map[string]interface{}{"time": 100000.00, "unit": "d", "money": 400.00, "expire": DATAITEM_PRICE_EXPIRE},
	}
	log.Println(chkPrice(price2, "flow"))

}