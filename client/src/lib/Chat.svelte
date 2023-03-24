<script lang="ts">
  export let messages = [];
  export let onChatSubmit: (message: string) => void;
  const onSubmit = () => {
    onChatSubmit(currentMessage)
  }
  let currentMessage = "";
</script>

<div class="mx-auto h-4/5">  
  {#each messages as message}
    <p>{message}</p>
  {/each}
</div>
<form on:submit|preventDefault={onSubmit} class="flex justify-center">
  <input
    class="h-8 w-96 border"
    type="text" 
    name="message-text" 
    id="message-text"
    bind:value={currentMessage}
  />
  <input
    class="w-32"
    type="submit"
    name="message-button"
    id="message-button"
  >
</form>