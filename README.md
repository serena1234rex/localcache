## Install

```
$ go get github.com/HoneyLiuJiaYi/localcache
```

## Example

### Manually set a key-value pair.

```go
func TestSimple(t *testing.T) {
	cache := localcache.Create().
		Tp(localcache.SIMPLE).
		Build()

	cache.Set("key", "ok")

	value, err := cache.Get("key")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(value)
}
```

```
Get: ok
```

### Manually set a key-value pair, with an expiration time.

```go
func TestExpire(t *testing.T) {
	cache := localcache.Create().
		Tp(localcache.SIMPLE).
		SetDuration(time.Millisecond * 10).
		Build()

	cache.Set("boy", "yes")

	time.Sleep(time.Duration(time.Millisecond * 2))

	value, _ := cache.Get("boy")
	fmt.Println(value)
}
```

### Manually set a key-value pair, with a flight register.

```go
func TestFlight(t *testing.T) {
	r := localcache.CreateRegister()

	cache := localcache.Create().
		Tp(localcache.SIMPLE).
		OpenFlight(&r).
		Build()

	cache.Set("a", "aa")

	cache.Get("a")
	cache.Get("b")

	fmt.Println(r.MissCount())
}
```

# Author
**Jiayu Liu**

* <https://github.com/HoneyLiuJiaYi>