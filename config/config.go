package config

var OneDay = 24 * 60 * 60

const Issuer = "sazi"
const JWTSecret = "dz"

// DSN const DSN = "sazi:dz@tcp(127.0.0.1:3306)/douyin?parseTime=true&charset=utf8&loc=Local"
const DSN = "root:123456@tcp(192.168.17.129:3306)/douyin?parseTime=true&charset=utf8&loc=Local"
const RedisAddr = "192.168.17.129:6379"
const FeedLimit = 30
