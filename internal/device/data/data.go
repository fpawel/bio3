package data

import (
	"github.com/fpawel/bio3/internal/products"

)

type DeviceInfoProvider struct {
	db products.DB
}

func FormatProductType (i int) string {
	return ProductTypes[i].What
}

func GoodProduct(product products.Product) bool {
	poductType := ProductTypes[product.Party.ProductTypeIndex()]
	for _, work := range MainTestsRele {
		for _,rele := range poductType.Reles {
			v := TestReleResult(product, work, rele)
			if v == nil {
				return false
			}
			if releValue := work.ReleValue(rele); releValue != nil {
				if *v != *releValue {
					return false
				}
			}
		}
	}
	return true
}

func BadProduct(product products.Product) bool {
	poductType := ProductTypes[product.Party.ProductTypeIndex()]
	for _, work := range MainTestsRele {
		for _,rele := range poductType.Reles {
			if v := TestReleResult(product, work, rele); v != nil {
				if releValue := work.ReleValue(rele); releValue != nil {
					if *v != *releValue {
						return true
					}
				}
			}
		}
	}
	return false
}


func TestReleResult(p products.Product, test *Work, rele *Rele) (result *bool){

	b := p.Party.Tx.Value(p.Test(test.What).Path(), rele.Index())
	if len(b) == 0{
		return nil
	}
	var r bool = b[0] != 0
	return &r
}


func SetTestReleResult( p products.Product, test *Work, rele *Rele, v *bool) {
	var b [] byte
	if v != nil {
		b = []byte{0}
		if *v {
			b[0] = 1
		}
	}
	p.Party.Tx.SetValue( p.Test(test.What).Path(), rele.Index(), b)
}

func (x *Work) PartyLogPath(tx products.Tx) products.DBPath   {
	return tx.Party().Test(x.What)
}

func (x *Work) WriteLog(tx products.Tx, timeKey [] byte, level int, text string)  {
	tx.WriteLog(x.PartyLogPath(tx).Path(), timeKey, level, text)
}