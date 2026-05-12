package db

import "fmt"

func DefineGormDSN(host string, user string, password string, port string) string {
	return fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, port)
}
