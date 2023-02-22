<script>
  import { invalidate } from "$app/navigation";

  export let id;
  export let name;
  export let desc;
  export let redemptions;
  export let expiry_date;

  async function redeem(id) {
    let abort = new AbortController()
    try {
      const response = await fetch(`http://localhost:8080/api/${id}/redeem`, {
        method: 'POST',
        signal: abort.signal
      });

      if (!response.ok) {
        throw Error(`${response.status} ${response.statusText}`);
      }
    } catch (err) {
      if (err.name == 'AbortError') {
        console.log('fetch aborted')
      } else if (err.message.includes('Syntax Error')) {
        throw Error('invalid query format');
      } else {
        console.error(err)
        throw Error('Something went wrong');
      }
    }
    // TODO handle error when fully redeemed

    // refresh fetch on each redemption
    invalidate('coupons:fetch');
  }
</script>

<article>
  <header>
    <h4>{name}</h4>
  </header>

  <p>{desc}</p>
  <button on:click={redeem(id)}>Redeem</button>

  <footer>
    <small>
      {redemptions} uses left
      <br>
      Expires: {expiry_date}
    </small>
  </footer>
</article>

<style>
h4 {
  margin-bottom: 0;
}

article {
  --card-box-shadow:
    0.0145rem 0.029rem 0.174rem rgba(27, 40, 50, 0.01698),
    0.0335rem 0.067rem 0.402rem rgba(27, 40, 50, 0.024),
    0.0625rem 0.125rem 0.75rem rgba(27, 40, 50, 0.03),
    0.1125rem 0.225rem 1.35rem rgba(27, 40, 50, 0.036),
    0.2085rem 0.417rem 2.502rem rgba(27, 40, 50, 0.04302),
    0.5rem 1rem 6rem rgba(27, 40, 50, 0.06),
    0 0 0 0.0625rem rgba(27, 40, 50, 0.015);
  margin: 2rem 0;
  margin-bottom: 0;
  padding: 2rem 1.8rem;
  box-shadow: var(--card-box-shadow);
  text-align: center;
}
article > header,
article > footer {
  margin-right: calc(1rem * -1);
  margin-left: calc(1rem * -1);
  padding: calc(2rem * 0.66) 1rem;
}
article > header {
  margin-top: calc(2rem * -1);
  margin-bottom: 2rem;
}
article > footer {
  margin-top: 2rem;
  margin-bottom: calc(2rem * -1);
}
</style>
