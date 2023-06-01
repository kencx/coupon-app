<script>
  import AddCoupon from "$lib/AddCoupon.svelte";

  let addCouponModal;
</script>

<main class="container">
  <nav>
    <ul>
      <li><strong>Coupons</strong></li>
    </ul>
    <ul>
      <li>
        <button on:click={() => addCouponModal.open()}>+ Add</button>
      </li>
    </ul>
  </nav>

  <slot />
</main>

<AddCoupon bind:this={addCouponModal} />

<style>
  .container {
    width: 100%;
    margin-right: auto;
    margin-left: auto;
    padding-right: 1rem;
    padding-left: 1rem;
  }

  @media (min-width: 576px) {
    .container {
      max-width: 510px;
      padding-right: 0;
      padding-left: 0;
    }
  }
  @media (min-width: 768px) {
    .container {
      max-width: 700px;
    }
  }
  @media (min-width: 992px) {
    .container {
      max-width: 920px;
    }
  }
  @media (min-width: 1200px) {
    .container {
      max-width: 1130px;
    }
  }

  :where(nav li)::before {
    float: left;
    content: "";
  }

  nav,
  nav ul {
    display: flex;
  }

  nav {
    justify-content: space-between;
    margin-bottom: 2.5rem;
  }
  nav ul {
    align-items: center;
    margin-bottom: 0;
    padding: 0;
    list-style: none;
  }
  nav li {
    display: inline-block;
    margin: 0;
    padding: 1rem 0.5rem;
  }
  button {
    border-radius: 6px;
  }
</style>
