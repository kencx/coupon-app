package main

import (
	"fmt"
	"log"
	"time"
)

type Coupon struct {
	Id          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"desc" db:"description"`
	Redemptions int       `json:"redemptions" db:"redemptions"`
	ExpiryDate  time.Time `json:"expiry_date" db:"expiry_date"`
}

func (c *Coupon) String() string {
	return fmt.Sprintf(`[id=%d redemptions=%d expirydate=%s]`, c.Id, c.Redemptions, c.ExpiryDate)
}

func (c *Coupon) redeem() error {
	if c.expired() {
		return fmt.Errorf("coupon %d has expired on %s", c.Id, c.ExpiryDate)
	}

	if c.Redemptions-1 < 0 {
		return fmt.Errorf("coupon %d has been fully redeemed", c.Id)
	}

	c.Redemptions -= 1
	log.Printf("coupon %d redeemed", c.Id)
	return nil
}

func (c *Coupon) expired() bool {
	return c.ExpiryDate.Before(time.Now())
}
