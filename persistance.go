package exstate

import "time"

type GetString func() string
type GetInt func() int
type GetFloat func() float64
type GetBool func() bool
type Setter[T any] func(val T)

type SetAdd func(value interface{})
type SetContains func(value interface{}) bool
type SetRemove func(value interface{})
type SetGet func() []string

type Source interface {
	String(path ...string) (GetString, Setter[string])
	Int(path ...string) (GetInt, Setter[int])
	Float(path ...string) (GetFloat, Setter[float64])
	Bool(path ...string) (GetBool, Setter[bool])
	Set(path ...string) (SetAdd, SetGet, SetRemove, SetContains)

	Close()
}

type Cache interface {
	CacheString(expires time.Duration, path ...string) (GetString, Setter[string])
	CacheInt(expires time.Duration, path ...string) (GetInt, Setter[int])
	CacheFloat(expires time.Duration, path ...string) (GetFloat, Setter[float64])
	CacheBool(expires time.Duration, path ...string) (GetBool, Setter[bool])

	Close()
}
