package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

func Open(dsn string) (*DB, error) {
	database, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("database failed to open: %v", err)
	}

	if err = database.Ping(); err != nil {
		return nil, fmt.Errorf("database failed to connect: %v", err)
	}

	d := &DB{db: database}
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
	return nil
}
