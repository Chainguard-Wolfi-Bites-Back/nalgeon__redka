// Package rstring is a database-backed string repository.
// It provides methods to interact with strings in the database.
package rstring

import (
	"database/sql"
	"time"

	"github.com/nalgeon/redka/internal/core"
	"github.com/nalgeon/redka/internal/sqlx"
)

// DB is a database-backed string repository.
// A string is a slice of bytes associated with a key.
// Use the string repository to work with individual strings.
type DB struct {
	*sqlx.DB[*Tx]
}

// New connects to the string repository.
// Does not create the database schema.
func New(db *sql.DB) *DB {
	d := sqlx.New(db, NewTx)
	return &DB{d}
}

// Get returns the value of the key.
// If the key does not exist or is not a string, returns ErrNotFound.
func (d *DB) Get(key string) (core.Value, error) {
	tx := NewTx(d.SQL)
	return tx.Get(key)
}

// GetMany returns a map of values for given keys.
// Ignores keys that do not exist or not strings,
// and does not return them in the map.
func (d *DB) GetMany(keys ...string) (map[string]core.Value, error) {
	tx := NewTx(d.SQL)
	return tx.GetMany(keys...)
}

// GetSet returns the previous value of a key after setting it to a new value.
// Optionally sets the expiration time (if ttl > 0).
// If the key already exists, overwrites the value and ttl.
// If the key exists but is not a string, returns ErrKeyType.
// If the key does not exist, returns nil as the previous value.
func (d *DB) GetSet(key string, value any, ttl time.Duration) (core.Value, error) {
	var val core.Value
	err := d.Update(func(tx *Tx) error {
		var err error
		val, err = tx.GetSet(key, value, ttl)
		return err
	})
	return val, err
}

// Incr increments the integer key value by the specified amount.
// Returns the value after the increment.
// If the key does not exist, sets it to 0 before the increment.
// If the key value is not an integer, returns ErrValueType.
// If the key exists but is not a string, returns ErrKeyType.
func (d *DB) Incr(key string, delta int) (int, error) {
	var val int
	err := d.Update(func(tx *Tx) error {
		var err error
		val, err = tx.Incr(key, delta)
		return err
	})
	return val, err
}

// IncrFloat increments the float key value by the specified amount.
// Returns the value after the increment.
// If the key does not exist, sets it to 0 before the increment.
// If the key value is not an float, returns ErrValueType.
// If the key exists but is not a string, returns ErrKeyType.
func (d *DB) IncrFloat(key string, delta float64) (float64, error) {
	var val float64
	err := d.Update(func(tx *Tx) error {
		var err error
		val, err = tx.IncrFloat(key, delta)
		return err
	})
	return val, err
}

// Set sets the key value that will not expire.
// Overwrites the value if the key already exists.
// If the key exists but is not a string, returns ErrKeyType.
func (d *DB) Set(key string, value any) error {
	err := d.Update(func(tx *Tx) error {
		return tx.Set(key, value)
	})
	return err
}

// SetExists sets the key value if the key exists.
// Optionally sets the expiration time (if ttl > 0).
// Returns true if the key was set, false if the key does not exist.
// If the key exists but is not a string, returns ErrKeyType.
func (d *DB) SetExists(key string, value any, ttl time.Duration) (bool, error) {
	var ok bool
	err := d.Update(func(tx *Tx) error {
		var err error
		ok, err = tx.SetExists(key, value, ttl)
		return err
	})
	return ok, err
}

// SetExpires sets the key value with an optional expiration time (if ttl > 0).
// Overwrites the value and ttl if the key already exists.
// If the key exists but is not a string, returns ErrKeyType.
func (d *DB) SetExpires(key string, value any, ttl time.Duration) error {
	err := d.Update(func(tx *Tx) error {
		return tx.SetExpires(key, value, ttl)
	})
	return err
}

// SetMany sets the values of multiple keys.
// Overwrites values for keys that already exist and
// creates new keys/values for keys that do not exist.
// Removes the TTL for existing keys.
// If any of the keys exists but is not a string, returns ErrKeyType.
func (d *DB) SetMany(items map[string]any) error {
	err := d.Update(func(tx *Tx) error {
		return tx.SetMany(items)
	})
	return err
}

// SetManyNX sets the values of multiple keys, but only if none
// of them yet exist. Returns true if the keys were set,
// false if any of them already exist.
// If any of the keys exists but is not a string, returns ErrKeyType.
func (d *DB) SetManyNX(items map[string]any) (bool, error) {
	var ok bool
	err := d.Update(func(tx *Tx) error {
		var err error
		ok, err = tx.SetManyNX(items)
		return err
	})
	return ok, err
}

// SetNotExists sets the key value if the key does not exist.
// Optionally sets the expiration time (if ttl > 0).
// Returns true if the key was set, false if the key already exists.
// If the key exists but is not a string, returns ErrKeyType.
func (d *DB) SetNotExists(key string, value any, ttl time.Duration) (bool, error) {
	var ok bool
	err := d.Update(func(tx *Tx) error {
		var err error
		ok, err = tx.SetNotExists(key, value, ttl)
		return err
	})
	return ok, err
}
