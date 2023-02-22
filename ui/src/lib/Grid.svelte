<script>
  import Coupon from "$lib/Coupon.svelte";

  export let data;
</script>

<div class="grid">
  {#if data.coupons.length < 1}
    <p>{data.statusText}</p>
  {:else}
    {#each data.coupons as coupon}
      <div><Coupon {...coupon}/></div>
    {/each}
  {/if}
</div>

<style>
.grid {
  /* auto scaling grid with max num of columns and min width*/
  /* https://css-tricks.com/an-auto-filling-css-grid-with-max-columns/ */
  --grid-layout-gap: 1rem;
  --grid-column-count: 4;
  --grid-item--min-width: 160px;

  --gap-count: calc(var(--grid-column-count) - 1);
  --total-gap-width: calc(var(--gap-count) * var(--grid-layout-gap));
  --grid-item--max-width: calc((100% - var(--total-gap-width)) / var(--grid-column-count));

  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(max(var(--grid-item--min-width), var(--grid-item--max-width)), 1fr));
  grid-gap: var(--grid-layout-gap);
  margin-bottom: 4rem;
}

/*
  TODO mobile width diff view
  horizontal card view
*/
</style>
