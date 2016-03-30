package main

type delete struct {
	Rep   repProxy      `json:"repository"`
	Items itemProxyList `json:"dataitems"`
}

type repProxy struct {
	Repository_name string `json:"repository_name, omitempty"`
	repository
}

type itemProxy struct {
	Repository_name string `json:"repository_name, omitempty"`
	Dataitem_name   string `json:"dataitem_name, omitempty"`
	dataItem
}

type itemProxyList []itemProxy

func newRepProxy(rep repository) repProxy {
	return repProxy{rep.Repository_name, rep}
}

func newitemsProxy(items []dataItem) itemProxyList {
	itemProxies := itemProxyList{}
	for _, item := range items {
		itemProxies = append(itemProxies, itemProxy{
			item.Repository_name,
			item.Dataitem_name,
			item,
		})
	}
	return itemProxies
}
