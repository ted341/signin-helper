package main

import (
	"log"
	"encoding/json"
)

const MyID = "Ud44e84d24b542d82756d6734ab63c1dd"
const CompanyID = "Ce87ea41acb66d329a86c50addb5e54fc"

type Profile struct {
	Name         string
	PhoneNumber  string
}

var MySquad = map[string]*Profile{
	"33076": &Profile{"張翔中", "0911855341"},
	"33077": &Profile{"林志謙", "0972739267"},
	"33078": &Profile{"陳毅",   "0972877326"},
	"33079": &Profile{"程憲文", "0979250392"},
	"33080": &Profile{"洪義軒", "0985054102"},
	"33081": &Profile{"黃健奇", "0928737611"},
	"33082": &Profile{"夏裕明", "0988326617"},
	"33083": &Profile{"林文揚", "0978212417"},
	"33084": &Profile{"白立弘", "0939868676"},
	"33085": &Profile{"潘立騏", "0926369958"},
	"33086": &Profile{"邱聖淯", "0953860305"},
	"33087": &Profile{"陳柏翰", "0970680730"},
	"33088": &Profile{"陳力維", "0933343659"},
	"33089": &Profile{"陳威丞", "0983572070"},
	"33090": &Profile{"呂曜銘", "0938533015"},
}

type User struct {
	StudentID    string
	Message      string
	Unblocked  	 bool
	*Profile
}

func (u User) ToJSON() (s string) {
	b, _ := json.Marshal(u)
	s = string(b)
	return
}

func GetUser(id string) *User {

	v, f, _ := __redis.Get(id)
	if f == false {
		log.Printf("user not found: %s", id)
		return nil
	}
	
	u := &User{}
	e := json.Unmarshal([]byte(v), u)
	if e != nil {
		log.Printf("parse user failed: %s", id)
		return nil
	}

	return u
}