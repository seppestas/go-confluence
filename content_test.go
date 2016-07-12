package confluence

import (
	"testing"
	"github.com/google/go-querystring/query"
	"time"
)

func TestURIBuilding(t *testing.T){
	now:=time.Now()
	q:=GetContentQuery{
		Limit:10,
		PostingDay:now,
		Expand:[]string{"one","two","three"},
	}

	values, err:=query.Values(q)
	if err!=nil{
		t.Fatalf("Can't build query err:%s", err)
	}

	if values.Get("limit")!="10" {t.Error("limit value is incorrect")}
	if values.Get("postingDay")!=now.Format("2006-01-02") {t.Errorf("postingDay value is incorrect '%s'",values.Get("postingDay"))}
	if values.Get("expand")!="one,two,three" {t.Errorf("expand value is incorrect '%s'",values.Get("expand"))}
}
