package main

import (
	"testing"
	"time"
)

func TestRedeemSuccess(t *testing.T) {
	c := &Coupon{
		Id:          1,
		Name:        "Adidas",
		Description: "20% off shoes",
		Redemptions: 5,
		ExpiryDate:  time.Now().Add(24 * time.Hour),
	}
	err := c.redeem()

	if c.Redemptions != 4 {
		t.Errorf("coupon was not redeemed")
	}

	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}
}

func TestExceedRedemptions(t *testing.T) {
	c := &Coupon{
		Id:          2,
		Name:        "Adidas",
		Description: "20% off shoes",
		Redemptions: 0,
		ExpiryDate:  time.Now().AddDate(0, 0, 7),
	}
	err := c.redeem()

	if c.Redemptions < 0 {
		t.Errorf("coupon redemptions < 0")
	}

	if err == nil {
		t.Errorf("expected err: coupon has been fully redeemed")
	}
}

func TestCouponExpired(t *testing.T) {
	c := &Coupon{
		Id:          3,
		Name:        "Adidas",
		Description: "20% off shoes",
		Redemptions: 5,
		ExpiryDate:  time.Now().AddDate(0, 0, -7),
	}
	err := c.redeem()

	if c.Redemptions != 5 {
		t.Errorf("expired coupon redeemed")
	}

	if err == nil {
		t.Errorf("expected err: coupon has expired")
	}
}
