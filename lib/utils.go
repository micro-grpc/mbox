package lib

import (
  "fmt"
  "math/rand"
  "regexp"
  "strings"
  "bytes"
  "math"
  "time"
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

// GetTotal - возвращает округленное число до большего
func GetTotal(cnt int64, max int64) int64 {
  total := int64(math.Ceil(float64(cnt) / float64(max)))
  return total
}

// GetTodayShow - возвращает текущую дату как строка 20180624
func GetTodayShow() string {
  dt := time.Now()
  dm := fmt.Sprintf("%d", int(dt.Month()))
  dd := fmt.Sprintf("%d", dt.Day())
  if int(dt.Month()) < 10 {
    dm = fmt.Sprintf("0%s", dm)
  }
  if dt.Day() < 10 {
    dd = fmt.Sprintf("0%s", dd)
  }
  return fmt.Sprintf("%d%s%s", dt.Year(), dm, dd)
}

// GetFullNumber - Возвращает номер заполняя в начале нулями определенной длины
// card - изнвчальное число
// limit - длина символов в результате
// format 1234567 -> 0001234567
func GetFullNumber(card string, limit int) (res string) {
  cl := len(card)
  if cl >= limit {
    res = card
  } else if cl < limit {
    lenCard := limit - cl
    for i := 0; i < lenCard; i++ {
      res = fmt.Sprintf("%s0", res)
    }
    res += card
  }
  return res
}

// IsExits - проверяем есть ли такой аргумент
//  args := []string{"sub", "arg1", "arg2"}
//  fmt.Println("is exits:", IsExits("arg2", args))
func IsExits(name string, args []string) bool {
  for _, arg := range args {
    if arg == name {
      return true
    }
  }
  return false
}

// Contains - ищем в масиве строк нужную строку
func Contains(a []string, x string) bool {
  for _, n := range a {
    if x == n {
      return true
    }
  }
  return false
}

// RandSeq - generate random string
func RandSeq(n int) string {
  var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
  rand.Seed(time.Now().UnixNano())
  b := make([]rune, n)
  for i := range b {
    b[i] = letters[rand.Intn(len(letters))]
  }
  return string(b)
}

// RandNumberSeq generate random number string
func RandNumberSeq(n int) string {
  var letters = []rune("1234567890")
  rand.Seed(time.Now().UnixNano())
  b := make([]rune, n)
  for i := range b {
    b[i] = letters[rand.Intn(len(letters))]
  }
  return string(b)
}
