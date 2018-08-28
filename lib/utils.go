package lib

import (
  "fmt"
  "regexp"
  "strings"
  "bytes"
)
// FilteringMode - return sql operator
func FilteringMode(src string, key int) string {
  switch src {
  case "EQ":
    return fmt.Sprintf("=$%d", key)
  case "NE":
    return fmt.Sprintf("!=$%d", key)
  case "GT":
    return fmt.Sprintf(">$%d", key)
  case "GE":
    return fmt.Sprintf(">=$%d", key)
  case "LT":
    return fmt.Sprintf("<$%d", key)
  case "LE":
    return fmt.Sprintf("<=$%d", key)
  case "NOT_NULL":
    return " IS NOT NULL"
  case "IS_NULL":
    return " IS NULL"
  default:
    return "!=0"
  }
}

// ReplaceNameToKey - преобразуем строку в ключ для CouchDB
// src := "Привет|name? \\ / {} ! `@#$%^&()-+=~ '<[ витя \"рога Копїта \"]>' = 1.2кг 1,03"
// stopWord := " кг | для | литр | і | и | нет"
// @response - привет_name_витя_рога_копїта
func ReplaceNameToKey(src string, stopWord string) (result string) {
  r, _ := regexp.Compile("[0-9,.\\]\\[<>|'?!`\\\\@#$%^&()\\-/{}+=~\"]+")
  res := r.ReplaceAllString(strings.ToLower(src), " ")
  r = regexp.MustCompile(fmt.Sprintf("(%s)+", stopWord))
  res = r.ReplaceAllString(res, " ")
  re := regexp.MustCompile("  +")
  replaced := re.ReplaceAll(bytes.TrimSpace([]byte(res)), []byte(" "))
  return strings.Replace(string(replaced), " ", "_", -1)
}
