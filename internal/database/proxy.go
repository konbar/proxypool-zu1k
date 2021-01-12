package database

import (
	"github.com/konbar/proxypool/pkg/proxy"
	"gorm.io/gorm"
)

type Proxy struct {
	gorm.Model
	proxy.Base
	Link       string
	Identifier string `gorm:"unique"`
}

func InitTables() {
	if DB == nil {
		err := connect()
		if err != nil {
			return
		}
	}
	err := DB.AutoMigrate(&Proxy{})
	if err != nil {
		panic(err)
	}
}

const roundSize = 100

func SaveProxyList(pl proxy.ProxyList) {
	if DB == nil {
		return
	}

	// TODO 批量插入因为是生成一个sql，如果插入失败（重复）就全都没了
	//size := pl.Len()
	//round := (size + roundSize - 1) / roundSize
	//
	//for r := 0; r < round; r++ {
	//	proxies := make([]Proxy, 0, roundSize)
	//	for i, j := r*roundSize, (r+1)*roundSize-1; i < j && i < size; i++ {
	//		p := pl[i]
	//		proxies = append(proxies, Proxy{
	//			Base:       *p.BaseInfo(),
	//			Link:       p.Link(),
	//			Identifier: p.Identifier(),
	//		})
	//	}
	//	DB.Create(&proxies)
	//}

	for _, p := range pl {
		DB.Create(&Proxy{
			Base:       *p.BaseInfo(),
			Link:       p.Link(),
			Identifier: p.Identifier(),
		})
	}
}

func GetAllProxies() (proxies proxy.ProxyList) {
	proxies = make(proxy.ProxyList, 0)
	if DB == nil {
		return
	}

	proxiesDB := make([]Proxy, 0)
	DB.Select("link").Find(&proxiesDB)

	for _, proxyDB := range proxiesDB {
		if proxiesDB != nil {
			p, err := proxy.ParseProxyFromLink(proxyDB.Link)
			if err == nil && p != nil {
				proxies = append(proxies, p)
			}
		}
	}
	return
}
