<script>
  import { PUBLIC_API_URL } from "$env/static/public";
  import Modal from "$lib/Modal.svelte";

  let redeemModal;
  let redeemModalText;
  let showRedeemModal = false;

  export let id;
  export let name;
  export let desc;
  export let redemptions;
  export let expiry_date;

  function format_date(date) {
    const d = new Date(date);
    return d.toLocaleDateString();
  }

  async function redeem(id) {
    try {
      const res = await fetch(`${PUBLIC_API_URL}/api/${id}/redeem`, {
        method: "POST",
      });

      if (!res.ok) {
        // TODO handle error when fully redeemed
        if (res.status == 400) {
          showRedeemModal = true
          redeemModalText = "Coupon fully redeemed!"
          return;
        }

        throw Error(`${res.status} ${res.statusText}`);
      }
    } catch (err) {
      console.error(err);
      throw Error("Something went wrong");
    }

    showRedeemModal = true
    redeemModalText = "Coupon redeemed!"
    // refresh coupon on each redemption
    refresh(id);
  }

  async function refresh(id) {
    await fetch(`${PUBLIC_API_URL}/api/coupon/${id}`)
      .then((res) => {
        if (!res.ok) {
          return { coupons: [] };
        }

        return res.json();
      })
      .then((updated) => {
        if (updated.coupons.length < 1) {
          throw Error("Error fetching coupon after redemption");
        }
        var new_coupon = updated.coupons[0];
        redemptions = new_coupon.redemptions;
      })
      .catch((err) => {
        console.error(err);
      });
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
      <br />
      Expires: {format_date(expiry_date)}
    </small>
  </footer>
</article>


<Modal bind:showModal={showRedeemModal} bind:this={redeemModal}>
  <div class="modal">
    {redeemModalText}
    <button on:click={() => redeemModal.close()}>OK</button>
  </div>
</Modal>

<style>
  h4 {
    margin-bottom: 0;
  }

  article {
    --card-box-shadow: 0.0145rem 0.029rem 0.174rem rgba(27, 40, 50, 0.01698),
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

  .modal {
    display: grid;
    justify-content: center;
    grid-gap: 1rem;
  }
</style>
