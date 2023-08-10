package dao

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

type CacheDao interface {
	CreateWithExpire(data any, expire float64) error
	Query(data any) error
	Close()
}

type CacheDaoRedis struct {
	RedisClient redis.Conn
}

// factory pattern
func NewCacheDaoRedis(conn redis.Conn) CacheDao {
	if conn == nil {
		return nil
	}

	return &CacheDaoRedis{conn}
}

func (cli *CacheDaoRedis) Close() {
	_ = cli.RedisClient.Close()
}

func (cli *CacheDaoRedis) CreateWithExpire(data any, expire float64) error {
	m := cli.toMap(data)
	key, val := "", ""

	for k, v := range m {
		key = k
		val = v
	}

	// SET key value EX seconds
	_, err := cli.RedisClient.Do("SET", key, val, "EX", expire)
	return err
}

func (cli *CacheDaoRedis) toMap(data any) map[string]string {
	result := make(map[string]string)

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	var keyStr *string
	var valStr *string

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		fieldValue := value.Field(i).String()

		tag := field.Tag.Get("cache")
		if tag == "key" {
			keyStr = &fieldValue
			if valStr != nil {
				result[*keyStr] = *valStr
			} else {
				result[*keyStr] = ""
			}
		} else if tag == "val" || tag == "value" {
			valStr = &fieldValue
			if keyStr != nil {
				result[*keyStr] = fieldValue
			} else {
				// do noting, wait found key
			}
		}
	}

	return result
}

func (cli *CacheDaoRedis) setValue(pair interface{}, v string) error {
	valuePtr := reflect.ValueOf(pair)
	if valuePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("data must be a pointer to a struct")
	}
	value := valuePtr.Elem()

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		fieldValue := value.Field(i)

		tag := field.Tag.Get("cache")

		if tag == "value" || tag == "val" {
			if fieldValue.CanSet() {
				switch fieldValue.Kind() {
				case reflect.String:
					fieldValue.SetString(v)
				case reflect.Int:
					// 假设value是一个整数字符串
					intValue, err := strconv.Atoi(v)
					if err != nil {
						return err
					}
					fieldValue.SetInt(int64(intValue))
				}
			} else {
				return fmt.Errorf("field %s is not settable", field.Name)
			}
		}
	}

	return nil
}

func (cli *CacheDaoRedis) Query(pair any) error {
	if pair == nil {
		return errors.New("Must post a non-nil imageCaptcha obj")
	}

	m := cli.toMap(pair)
	key := ""

	for k := range m {
		key = k
	}

	/**
	return of redis.Conn.Do("GET" ...):
		Not found:		nil, nil
		found:			[]byte, nil
		occur a error:	nil, error
	*/
	code, err := redis.String(cli.RedisClient.Do("GET", key))
	if err != nil {
		if cli.IsNotExist(err) {
			return NotExistError
		}
		return err
	}

	err = cli.setValue(pair, code)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (cli *CacheDaoRedis) IsNotExist(err error) bool {
	return errors.Is(err, redis.ErrNil)
}
