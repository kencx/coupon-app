<script>
  import Modal from "$lib/Modal.svelte";
  import { AddCoupon } from "$lib/addCoupon.js";
  import { invalidate } from "$app/navigation";

  let formModal;
  let validateModal;
  let inputCode;
  let addCouponForm;

  export function open() {
    formModal.open();
  }

  function handleSubmit() {
    AddCoupon(inputCode);

    formModal.close();
    addCouponForm.reset();

    validateModal.open();
    invalidate("coupons:fetch");
  }
</script>

<Modal bind:this={formModal}>
  <article>
    <form on:submit|preventDefault={handleSubmit} bind:this={addCouponForm}>
      <h3>Add Coupon</h3>
      <input
        class="modal-input"
        type="text"
        placeholder="Coupon code"
        autocomplete="off"
        required
        bind:value={inputCode}
      />
      <div class="modal-buttons">
        <button on:click={() => formModal.close()}>Cancel</button>
        <button type="submit">Submit</button>
      </div>
    </form>
  </article>
</Modal>

<Modal bind:this={validateModal}>
  <div class="modal">Coupon added!</div>
  <button on:click={() => validateModal.close()}>OK</button>
</Modal>

<style>
  .modal-input {
    width: 100%;
  }

  .modal-buttons {
    margin: 1rem 0;
  }
</style>
