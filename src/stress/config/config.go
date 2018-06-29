package  config

import (
	"time"
	"math/rand"
	"stress/structs"
    "io/ioutil"
    "os"
	"fmt"
	"encoding/json"
	"strings"
	"strconv"
	"bytes"
)

func DoBody(conf structs.RequestData) (string, structs.SoapData) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	body := conf.Body
	soap := structs.SoapData{[]structs.LogData{}}

	for i := range conf.Replaces {
		if (conf.Replaces[i].Type == "array"){
			if(conf.Replaces[i].Sort == "rand"){
				v := r.Intn(len(conf.Replaces[i].Values))
				j := strconv.Itoa(conf.Replaces[i].Values[v])
				body = strings.Replace(body, conf.Replaces[i].Mark, j, 1)
				soap.Rpl = append(soap.Rpl, structs.LogData{ conf.Replaces[i].Mark, j })
			}
		}else if (conf.Replaces[i].Type == "date"){
			today := time.Now()
			newDate := today.Add(time.Duration(conf.Replaces[i].Value * (1000 * 1000 * 1000)))
			date :=  newDate.Format(conf.Replaces[i].Format)
			body = strings.Replace(body, conf.Replaces[i].Mark, date, 1)
			soap.Rpl = append(soap.Rpl, structs.LogData{ conf.Replaces[i].Mark, date } )
		}
	}
    return body, soap
}


func DoQueryString(conf structs.RequestData) (string) {
    var buffer bytes.Buffer
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	url := conf.Url
	buffer.WriteString(url)

	if (len(conf.Query) > 0){
		buffer.WriteString("?")
	}

	for i := range conf.Query {
		v := r.Intn(len(conf.Query[i].Values))
		j := conf.Query[i].Values[v]
		if( j != "" ){
			if (i >= 1){
				buffer.WriteString("&")
			}
			buffer.WriteString(conf.Query[i].Name)
			buffer.WriteString("=")
			buffer.WriteString(j)
		}
	}
    return buffer.String()
}



func GetFile(file string) structs.RequestData {
    raw, err := ioutil.ReadFile(file)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    var c structs.RequestData
    json.Unmarshal(raw, &c)
    return c
}
