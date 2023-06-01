import { PUBLIC_API_URL } from "$env/static/public";

const coupons = {
  TEST: {
    name: "Adidas",
    desc: "20% off shoes",
    redemptions: 5,
    expiry_date: "2023-12-31T00:00:00Z",
  },
  "15JEANS": {
    name: "Levis",
    desc: "15% off jeans",
    redemptions: 3,
    expiry_date: "2023-12-31T00:00:00Z",
  },
  "40OFF": {
    name: "Nike",
    desc: "40% off shirts",
    redemptions: 2,
    expiry_date: "2023-12-31T00:00:00Z",
  },
};

export async function AddCoupon(code) {
  // validate coupon code

  try {
    const res = await fetch(`${PUBLIC_API_URL}/api/coupons`, {
      method: "POST",
      body: JSON.stringify(coupons[code]),
    });
    if (!res.ok) {
      throw Error(`${res.status} ${res.statusText}`);
    }
  } catch (err) {
    console.error(err);
    return;
  }
  console.log("Coupon added!");
}
