import { PUBLIC_API_URL } from '$env/static/public';

export async function load({ fetch, depends }) {
    const res = await fetch(`${PUBLIC_API_URL}/api/coupons`);

    if (!res.ok) {
        return {
            coupons: [],
            statusText: 'No coupons found'
        }
    }

    depends('coupons:fetch');

    const json = await res.json();
    return {
        coupons: json.coupons,
        statusText: '',
    }
}
