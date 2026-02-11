package eparser

import (
	"fmt"
	"testing"
)

func TestSql(t *testing.T) {

	sql, parameters, err := ParseSql(`
select *  from yxhis${aaa.substring(0,4)}..tbmzfymx${aaa.substring(0,6)} where 1=1 
      ?{cbrh , and cmzh=#{cbrh} }
       ?{ typeof cbrid !== 'undefined' && cbrid!=null  && cbrid!='' ,   and cbrid=#{cbrid} }
      ?{ ighzl!=null&&ighzl.length>0 , and ighzl in( #{ighzl} )}
        `, map[string]interface{}{"aaa": "20250506", "cbrh": "1234", "cbrid": "", "ighzl": []int{1, 2, 3}})

	fmt.Println(sql, parameters, err)

}

func TestNameb(t *testing.T) {
	sql, parameters, err := ParseSql(`
select *   from yxhis${aaa.substring(0,4)}..tbmzfymx${aaa.substring(0,6)} where 1=1 
      ?{cbrh , and cmzh=#{cbrh} }
       ?{  !cbrid    || cbrid =='' ,   and cbrid="" }
      ?{ ighzl!=null&&ighzl.length>0 , and ighzl in( #{ighzl} )}
        `, map[string]interface{}{"aaa": "20250506", "cbrh": "1234", "cbrid": "", "ighzl": []int{1, 2, 3}})

	fmt.Println(sql, parameters, err)
}

func TestNameSql(t *testing.T) {
	sql, parameters, err := ParseSql("select PDFName from ${tableName} with(nolock) where CZYH=#{cbrh} and CBH=#{cbh}", map[string]any{
		"tableName": "tableName",
		"cbh":       "cbh",
		"cbrh":      "cbrh",
	})

	fmt.Println(sql, parameters, err)
}
