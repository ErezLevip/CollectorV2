package DbHandler

import (
	"Collector/Entities"
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
	"strconv"
)

const mysqlDriver = "mysql"

type DbHandler interface {
	LogRequest(requests chan Entities.CollectorRequest) error
}

type MySqlDbHandler struct {
	ConStr string
}

func (m MySqlDbHandler) Make(conStr string) DbHandler {
	return &MySqlDbHandler{
		ConStr: conStr,
	}
}

type ParamKeyValue struct {
	Key string
	Type  reflect.Type
	Value interface{}
}

const DateFormatters = "'%m/%d/%Y'"

func queryParametersCreator(params []ParamKeyValue, values bool) string {
	var b bytes.Buffer //query params buffer

	for _, kvp := range params {
		var value string
		if (kvp.Type == reflect.TypeOf(time.Time{})) { //insert time aas parsed date
			t := kvp.Value.(time.Time)
			value += fmt.Sprintf("STR_TO_DATE('%d/%d/%d',%s),", t.Month(), t.Day(), t.Year(), DateFormatters)
		} else {
			value = "'" + kvp.Value.(string) + "',"
		}

		if (!values) {
			b.WriteString(kvp.Key + "=" + value + " ") // append each param to its value
		} else {
			b.WriteString(value)
		}
	}
	return strings.Trim(strings.Trim(b.String(), " "), ",")
}

// can be turned into a bulk insert if needed using a buffered channel.
func (m *MySqlDbHandler) LogRequest(requests chan Entities.CollectorRequest) error {

	db, err := sql.Open(mysqlDriver, m.ConStr) // open connection to mysql database

	go func() {

		defer db.Close()
		var queryBuffer bytes.Buffer
		queryBuffer.WriteString("INSERT INTO `events` VALUES ")
		itemsCount := 0
		for request := range requests {
			log.Println("handle request", request)

			// collect all the parameters in a map
			params := make([]ParamKeyValue,8)
			params[0] = ParamKeyValue{Value: "0", Type: reflect.TypeOf(request.Gender), Key: "idevents"}
			params[1] = ParamKeyValue{Value: request.Gender, Type: reflect.TypeOf(request.Gender), Key: "gender"}
			params[2] = ParamKeyValue{Value: strconv.FormatFloat(request.D1, 'f', 2, 64), Type: reflect.TypeOf(request.D1), Key: "d1"}
			params[3] = ParamKeyValue{Value: strconv.FormatFloat(request.P1, 'f', 2, 64), Type: reflect.TypeOf(request.P1), Key:"p1"}
			params[4] = ParamKeyValue{Value: request.App_name, Type: reflect.TypeOf(request.App_name), Key:"app_name"}
			params[5] = ParamKeyValue{Value: request.Mime_type, Type: reflect.TypeOf(request.Mime_type), Key:"mime_type"}
			params[6] = ParamKeyValue{Value: strconv.FormatFloat(request.N1, 'f', 2, 64), Type: reflect.TypeOf(request.N1), Key:"n1"}
			params[7] = ParamKeyValue{Value: request.Cc, Type: reflect.TypeOf(request.Cc), Key:"cc"}

			// create sql key value string
			queryParams := queryParametersCreator(params, true)
			queryParams = strings.Trim(queryParams, " ")
			queryParams = fmt.Sprintf("(%s),", queryParams)
			queryBuffer.WriteString(queryParams)

			itemsCount++

			if (itemsCount == cap(requests)) {
				query := strings.Trim(queryBuffer.String(), ",")
				log.Println("execute", query)
				_, err = db.Exec(query) // execute
				if err != nil {
					log.Println(err.Error())
				}
				itemsCount = 0
				queryBuffer.Reset()
				queryBuffer.WriteString("INSERT INTO events VALUES ")
			}
		}

		log.Println("insert is done")
	}()
	return err
}

