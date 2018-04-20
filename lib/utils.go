package lib

import (
  "fmt"
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
