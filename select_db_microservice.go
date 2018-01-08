package select_db_microservice

import (
  "os"
  "fmt"
  "net/http"
  "github.com/op/go-logging"
  "github.com/json-iterator/go"
  "github.com/julienschmidt/httprouter"
  "github.com/pifagor87/conect_db_microservice"
)

/* Constant error number. */
const errc1 int = 41
const errc2 int = 42
const errc3 int = 43
const errc4 int = 44
const errc5 int = 45
const errc6 int = 46
const errc7 int = 47
const errc8 int = 48

/* Log error. */
var log = logging.MustGetLogger("sql")

/* json iterator variable. */
var json = jsoniter.ConfigCompatibleWithStandardLibrary

/* Structe POST data. */
type JsonData struct {
  Tables Tables `json:"tables"`
  Filters Filter `json:"filters"`
  Fields JsonArray `json:"fields"`
  Params Params `json:"params"`
}

/* Structe POST data field tables. */
type Tables struct {
  Origin TablesOrigin `json:"origin"`
  Join TablesJoin `json:"join"`
}

/* Structe POST data field table origin. */
type TablesOrigin struct {
  Table string `json:"table"`
  Alias string `json:"alias"`
}

/* Structe POST data field table Join. */
type TablesJoin []struct {
  Table string `json:"table"`
  Name string `json:"name"`
  Alias string `json:"alias"`
  Left string `json:"left"`
  Right string `json:"right"`
}

/* Structe POST data field filter. */
type Filter struct {
  And FilterData `json:"and"`
  Or FilterData `json:"or"`
}

/* Structe POST data field filter value. */
type FilterData []struct {
  Column string `json:"column"`
  Val JsonArray `json:"val"`
  Operator string `json:"operator"`
}

/* Structe POST data field paranetrs. */
type Params struct {
  Order Order `json:"order"`
  Group JsonArray `json:"group"`
  Limit string `json:"limit"`
}

/* Structe POST data field orders. */
type Order struct {
  Fields JsonArray `json:"fields"`
  Sort JsonArray `json:"sort"`
}

/* Type array. */
type JsonArray []string

/* Generate first level array. */
func result() (r map[string]interface{}) {
  return make(map[string]interface{})
}

/* Generate second level array. */
func results() (rs map[string]map[string]interface{}) {
  return make(map[string]map[string]interface{})
}

/* Load DB content values. */
func SSDomain(AccessDbPatch string)  httprouter.Handle {
  return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    start := true
    r.ParseForm()
    data := r.FormValue("data")
    var jd JsonData
    err := json.Unmarshal([]byte(data), &jd)
    if err != nil {
      start = false
      er := loadDataMessage("Json data error!", errc1)
      w.Write([]byte(er))
    }
    table, er := loadTableValue(jd.Tables.Origin, errc2)
    if er != "" {
      start = false
      w.Write([]byte(er))
    }
    fields, er := loadFieldsValue(jd.Fields, errc3)
    if er != "" {
      start = false
      w.Write([]byte(er))
    }
    whereAnd, er := loadWhere(jd.Filters.And, "and", errc4)
    if er != "" {
      start = false
      w.Write([]byte(er))
    }
    whereOr, er := loadWhere(jd.Filters.Or, "or", errc5)
    if er != "" {
      start = false
      w.Write([]byte(er))
    }
    params, er := loadParamsValue(jd.Params, errc6)
    if er != "" {
      start = false
      w.Write([]byte(er))
    }
    if start == true {
      db := conect_db_microservice.ConectPosqgresqlDb(AccessDbPatch)
      query := "SELECT " + fields + table + " WHERE" + whereAnd
      if whereOr != "" {
        query += " or " + whereOr
      }
      query += params
      rows, err := db.Query(query)
      if err != nil {
        log.Critical(err)
        log.Fatal(err)
      }
      defer rows.Close()
      cols := rows.FieldDescriptions()
      for rows.Next() {
        rowValues, err := rows.Values()
        if err != nil {
          log.Fatal(err)
        }
        m := result()
        rs := results()
        for i, col := range rowValues {
          if cols[i].Name != "" {
            m[cols[i].Name] = col;
          }
        }
        rs["data"] = m
        result, err1 := json.Marshal(rs)
        if err1 != nil {
          er := loadDataMessage("Json result data error!", errc7)
          w.Write([]byte(er))
        }
        w.Write([]byte(string(result)))
      }
    }
  }
}

/* Settings for error. */
func loadError() {
  var format = logging.MustStringFormatter(
    `%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
  )
  backend1 := logging.NewLogBackend(os.Stderr, "", 0)
  backend2 := logging.NewLogBackend(os.Stderr, "", 0)
  backend2Formatter := logging.NewBackendFormatter(backend2, format)
  backend1Leveled := logging.AddModuleLevel(backend1)
  backend1Leveled.SetLevel(logging.ERROR, "")
  logging.SetBackend(backend1Leveled, backend2Formatter)
}

/* Generate sql with POST data field table. */
func loadTableValue(jd TablesOrigin, er int) (string, string) {
  if jd.Table == "" {
    return "", loadDataMessage("No area origin table!", er)
  } else if jd.Alias == "" {
    return "", loadDataMessage("No area origin table alias!", er)
  }
  table := " FROM " + jd.Table + " AS " + jd.Alias
  return table, ""
}

/* Generate sql with POST data field table data. */
func loadFieldsValue(jd JsonArray, er int) (string, string) {
  fields := ""
  for index, element := range jd {
    if index == 0 {
      fields += element
    } else {
      fields += ", " + element
    }
  }
  if fields == "" {
    return "", loadDataMessage("No date area fields!", er)
  }
  return fields, ""
}

/* Generate sql with POST data field param. */
func loadParamsValue(jd Params, er int) (string, string) {
  data, sort, group := "", "", ""
  if len(jd.Order.Fields) > 0 {
    for index, element := range jd.Order.Fields {
      if element != "" && jd.Order.Sort[index] != "" {
        if index == 0 {
          sort += " ORDER BY " + element + " " + jd.Order.Sort[index]
        } else {
          sort += ", " + element + " " + jd.Order.Sort[index]
        }
      } else if element != "" && jd.Order.Sort[index] == "" {
        return "", loadDataMessage("No valid value sort Fields!", er)
      }
    }
  }
  if len(jd.Group) > 0 {
    for indexG, elementG := range jd.Group {
      if indexG == 0 {
        group += " GROUP BY " + elementG
      } else {
        group += ", " + elementG
      }
    }
  }
  if group != "" || sort != "" {
    data += group + sort
  }
  if jd.Limit != "" {
    data += " LIMIT " + jd.Limit
  }
  return data, ""
}

/* Generate sql with POST data field filters. */
func loadWhere(jd FilterData, ident string, er int) (string, string) {
  data, sqlv, operator := "", "", ""
  for index, element := range jd {
    if element.Column == "" {
      return "", loadDataMessage("No value Column!", er)
    } else if len(element.Val) == 0 {
      return "", loadDataMessage("No value Val!", er)
    } else if element.Operator == "" {
      return "", loadDataMessage("No value Operator!", er)
    }
    if len(element.Val) > 0 {
      if len(element.Val) > 1 {
        data = "in ("
        for indexV, _ := range jd {
          if indexV == 0 {
            data += element.Val[indexV]
          } else {
            data += ", " + element.Val[indexV]
          }
        }
        data += ")"
      } else {
        data = element.Val[0]
      }
    }
    if element.Val[0] != "" && (element.Operator == "ILIKE" || element.Operator == "LIKE")  {
      operator = " " + element.Operator + " '%' || '"+ element.Val[0] + "' || '%'"
    } else if data != "" {
      operator = element.Operator + data
    }
    if index == 0 {
      sqlv += " " + element.Column + operator
    } else {
      sqlv += " " + ident + " " + element.Column + operator
    }
  }
  return sqlv, ""
}

/* Callback errors. */
func loadDataMessage(er string, ern int) (string) {
  loadError()
  err := fmt.Sprintf("Error: %d. %s", ern, er)
  log.Error(err)
  r := result()
  rs := results()
  r["text"] = err
  r["id"] = ern
  rs["error"] = r
  result, error := json.Marshal(rs)
  if error != nil {
    er := fmt.Sprintf("Error: %d. Json critical error!", errc8)
    log.Critical(er)
    log.Fatal(er)
  }
  return string(result)
}