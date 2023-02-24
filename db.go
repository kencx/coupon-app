package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const (
	getAllKey = "getAllCoupons"
)

type DB struct {
	db    *sql.DB
	ctx   context.Context
	redis *redis.Client
}

func Open(dbDSN, cacheDSN string) (*DB, error) {
	database, err := sql.Open("postgres", dbDSN)
	if err != nil {
		return nil, fmt.Errorf("database failed to open: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cacheDSN,
		Password: "",
		DB:       0,
	})

	d := &DB{
		db:    database,
		ctx:   context.Background(),
		redis: rdb,
	}

	if err = database.Ping(); err != nil {
		return nil, fmt.Errorf("database failed to connect: %v", err)
	}

	err = rdb.Ping(d.ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("cache failed to connect: %v", err)
	}

	if err := d.initTable(); err != nil {
		return nil, fmt.Errorf("failed to initialize :%v", err)
	}
	return d, nil
}

func (d *DB) initTable() error {
	_, err := d.db.Exec(`CREATE TABLE IF NOT EXISTS coupons
		(id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        name TEXT,
        description TEXT,
        redemptions INT,
        expiry_date TIMESTAMP)`)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetCoupon(id int64) (*Coupon, error) {
	var coupon Coupon
	err := d.db.QueryRow("SELECT * FROM coupons WHERE id=$1", id).
		Scan(&coupon.Id, &coupon.Name, &coupon.Description, &coupon.Redemptions, &coupon.ExpiryDate)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("db: coupon %d does not exist", id)
	}

	if err != nil {
		return nil, fmt.Errorf("db: get coupon %d failed: %v", id, err)
	}
	return &coupon, nil
}

func (d *DB) GetAllCoupons() ([]*Coupon, error) {
	var coupons []*Coupon

	// query cache
	value, err := d.redis.Get(d.ctx, getAllKey).Result()
	if err != nil {
		// query db if cache miss
		if err == redis.Nil {
			log.Printf("key %q does not exist, querying database", getAllKey)

			coupons, err = d.getAllCouponsDB()
			if err != nil {
				return nil, err
			}

			// add to cache upon successful db query
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if err := enc.Encode(coupons); err != nil {
				return nil, fmt.Errorf("cache: failed to gob encode: %v", err)
			}

			err = d.redis.Set(d.ctx, getAllKey, buf.Bytes(), 0).Err()
			if err != nil {
				return nil, fmt.Errorf("cache: set data to cache failed: %v", err)
			}

		} else {
			return nil, fmt.Errorf("cache: get key %q failed: %v", getAllKey, err)
		}
	} else {
		buf := bytes.NewBuffer([]byte(value))
		dec := gob.NewDecoder(buf)
		if err := dec.Decode(&coupons); err != nil {
			return nil, fmt.Errorf("cache: failed to gob decode: %v", err)
		}
	}

	return coupons, nil
}

func (d *DB) getAllCouponsDB() ([]*Coupon, error) {
	var coupons []*Coupon

	rows, err := d.db.Query(`SELECT * FROM coupons`)
	if err != nil {
		return nil, fmt.Errorf("db: get coupons failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var coupon Coupon
		if err := rows.Scan(&coupon.Id, &coupon.Name, &coupon.Description, &coupon.Redemptions, &coupon.ExpiryDate); err != nil {
			return nil, fmt.Errorf("db: get coupons failed: %v", err)
		}
		coupons = append(coupons, &coupon)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("db: get coupons failed: %v", err)
	}
	return coupons, nil
}

func (d *DB) AddCoupon(c Coupon) error {
	stmt, err := d.db.Prepare(`INSERT INTO coupons
        (name, description, redemptions, expiry_date) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return fmt.Errorf("db: add coupon %d failed: %v", c.Id, err)
	}

	if _, err := stmt.Exec(c.Name, c.Description, c.Redemptions, c.ExpiryDate); err != nil {
		return fmt.Errorf("db: add coupon %d failed: %v", c.Id, err)
	}
	log.Printf("db: added new coupon: %v", c)

	if err := d.invalidateCache(getAllKey); err != nil {
		return err
	}
	return nil
}

func (d *DB) UpdateCoupon(c Coupon) error {
	stmt, err := d.db.Prepare("UPDATE coupons SET redemptions=$1 WHERE id=$2")
	if err != nil {
		return fmt.Errorf("db: update coupon %d failed: %v", c.Id, err)
	}

	if _, err := stmt.Exec(c.Redemptions, c.Id); err != nil {
		return fmt.Errorf("db: update coupon %d failed: %v", c.Id, err)
	}
	log.Printf("db: updated coupon: %d", c.Id)

	if err := d.invalidateCache(getAllKey); err != nil {
		return err
	}
	return nil
}

func (d *DB) DeleteCoupon(id int64) error {
	stmt, err := d.db.Prepare("DELETE FROM coupons WHERE id=$1")
	if err != nil {
		return fmt.Errorf("db: delete coupon %d failed: %v", id, err)
	}

	if _, err := stmt.Exec(id); err != nil {
		return fmt.Errorf("db: delete coupon %d failed: %v", id, err)
	}
	log.Printf("db: deleted coupon: %d", id)

	if err := d.invalidateCache(getAllKey); err != nil {
		return err
	}
	return nil
}

func (d *DB) invalidateCache(key string) error {
	err := d.redis.Del(d.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("cache: failed to invalidate cache: %v", err)
	}

	log.Printf("cache: key %q deleted", key)
	return nil
}
