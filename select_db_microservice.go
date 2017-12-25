package select_db_microservice

import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "github.com/julienschmidt/httprouter"
  "github.com/pifagor87/conect_db_microservice"
  "github.com/pifagor87/authorization_microservice"
)

const AccessDbPatch string = "connect_db.json"
const Port string = "2378"

type DataStructure struct {
  id  uint64
  url string
}

type JsonData struct {
  Tables Tables `json:"tables"`
  Filters Filter `json:"filters"`
  Fields JsonArray `json:"fields"`
  Params Params `json:"params"`
}

type Tables struct {
  Origin TablesOrigin `json:"origin"`
  Join TablesJoin `json:"join"`
}

type TablesOrigin struct {
  Table string `json:"table"`
  Alias string `json:"alias"`
}

type TablesJoin []struct {
  Table string `json:"table"`
  Name string `json:"name"`
  Alias string `json:"alias"`
  Left string `json:"left"`
  Right string `json:"right"`
}

type Filter struct {
  And FilterData `json:"and"`
  Or FilterData `json:"or"`
}

type FilterData []struct {
  Column string `json:"column"`
  Val JsonArray `json:"val"`
  Operator string `json:"operator"`
}

type Params struct {
  ORDER_BY JsonArray `json:"ORDER BY"`
  DESC string `json:"DESC"`
  LIMIT string `json:"LIMIT"`
  Param string `json:"param"`
}

type JsonArray []string

func SearchLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  start := true
  r.ParseForm()
  data := r.FormValue("data")
  var jd JsonData
  err := json.Unmarshal([]byte(data), &jd)
  if err != nil {
    start = false
    er := LoadDataMessage("Json data error!", 41)
    w.Write([]byte(er))
  }
  table, er := LoadTableValue(jd.Tables.Origin, 42)
  if er != "" {
    start = false
    w.Write([]byte(er))
  }
  fields, er := LoadFieldsValue(jd.Fields, 43)
  if er != "" {
    start = false
    w.Write([]byte(er))
  }
  whereAnd, er := LoadWhere(jd.Filters.And, "and", 44)
  if er != "" {
    start = false
    w.Write([]byte(er))
  }
  fmt.Println(whereAnd)
  if start == true {
    fmt.Println(fields)
    db := conect_db_microservice.ConectPosqgresqlDb(AccessDbPatch)
    query := "SELECT " + fields
    query += " FROM " + table
    query += " WHERE "
    query += whereAnd + " LIMIT 10"
    rows, err := db.Query(query)
    if err != nil {
      log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
      dt := new(DataStructure)
      if err := rows.Scan(&dt.id, &dt.url); err != nil {
        log.Fatal(err)
      }
      fmt.Printf("%d is %s\n", dt.id, dt.url)
    }
  }
}

func LoadTableValue(jd TablesOrigin, er int) (string, string) {
  if jd.Table == "" {
    return "", LoadDataMessage("No area origin table!", er)
  } else if jd.Alias == "" {
    return "", LoadDataMessage("No area origin table alias!", er)
  }
  table := jd.Table + " AS " + jd.Alias
  return table, ""
}

func LoadFieldsValue(jd JsonArray, er int) (string, string) {
  fields := ""
  for index, element := range jd {
    if index == 0 {
      fields += element
    } else {
      fields += ", " + element
    }
  }
  if fields == "" {
    return "", LoadDataMessage("No date area fields!", er)
  }
  return fields, ""
}

func LoadWhere(jd FilterData, ident string, er int) (string, string) {
  data, sql, operator := "", "", ""
  for index, element := range jd {
    if element.Column == "" {
      return "", LoadDataMessage("No value Column!", er)
    } else if len(element.Val) == 0 {
      return "", LoadDataMessage("No value Val!", er)
    } else if element.Operator == "" {
      return "", LoadDataMessage("No value Operator!", er)
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
      sql += element.Column + operator
    } else {
      sql += " " + ident + " " + element.Column + operator
    }
  }
  return sql, ""
}

func LoadDataMessage(er string, ern int) (string) {
  return fmt.Sprintf("Error: %d. %s", ern, er)
}

func main() {
  router := httprouter.New()
  router.POST("/protected-select/", authorization_microservice.ProtectedEndpoint(SearchLocation, AccessUserPatch))
  log.Fatal(http.ListenAndServe(":" + Port, router))
}