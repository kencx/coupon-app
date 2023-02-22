package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type envelope map[string]interface{}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
		next.ServeHTTP(w, r)
	})
}

func (a *App) getCouponsHandler(rw http.ResponseWriter, r *http.Request) {
	coupons, err := a.db.GetAllCoupons()
	if err != nil {
		errorResponse(rw, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	if len(coupons) < 1 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	writeResponse(rw, http.StatusOK, envelope{"coupons": coupons})
}

func (a *App) addCouponHandler(rw http.ResponseWriter, r *http.Request) {

	var coupon Coupon
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&coupon)
	if err != nil {
		errorResponse(rw, http.StatusBadRequest, fmt.Sprintf("err decoding request: %v", err))
		return
	}

	err = a.db.AddCoupon(coupon)
	if err != nil {
		errorResponse(rw, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (a *App) redeemCouponHandler(rw http.ResponseWriter, r *http.Request) {
	id := handleInt64("id", rw, r)

	coupon, err := a.db.GetCoupon(id)
	if err != nil {
		errorResponse(rw, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	if err = coupon.redeem(); err != nil {
		errorResponse(rw, http.StatusBadRequest, fmt.Sprint(err))
		return
	}

	if err = a.db.UpdateCoupon(*coupon); err != nil {
		errorResponse(rw, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (a *App) deleteCouponHandler(rw http.ResponseWriter, r *http.Request) {
	id := handleInt64("id", rw, r)

	if err := a.db.DeleteCoupon(id); err != nil {
		errorResponse(rw, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func handleInt64(key string, rw http.ResponseWriter, r *http.Request) int64 {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars[key])
	if err != nil {
		errorResponse(rw, http.StatusBadRequest, fmt.Sprintf("err decoding request: %v", err))
	}
	return int64(id)
}

func errorResponse(rw http.ResponseWriter, statusCode int, message string) {
	writeResponse(rw, statusCode, envelope{"error": message})
}

func writeResponse(rw http.ResponseWriter, statusCode int, v interface{}) {
	res, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		log.Printf("error writing response: %v", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	rw.Write(res)
}
