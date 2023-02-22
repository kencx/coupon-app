export async function load({ fetch, depends }) {
    const res = await fetch(`http://localhost:8080/api/coupons`);

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
