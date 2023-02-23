import { PUBLIC_API_URL } from "$env/static/public";

export async function load({ fetch, depends }) {
  try {
    const res = await fetch(`${PUBLIC_API_URL}/api/coupons`);
    depends("coupons:fetch");

    if (!res.ok) {
      return { coupons: [] };
    }
    return await res.json();
  } catch (err) {
    console.error(err);
    throw Error("Something went wrong");
  }
}
