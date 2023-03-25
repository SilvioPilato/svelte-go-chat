<script lang="ts">
    export let onLoginSubmit: (username: string) => void;
    let currentUsername: string;

    let onSubmit = () => {
        onLoginSubmit(currentUsername);
        currentUsername = "";
    }
</script>

<div class="w-64 h-96 bg-slate-50 flex flex-col items-center justify-evenly font-mono rounded-xl drop-shadow-md">
    <h1 class="text-3xl">Chat login</h1>
    <form on:submit|preventDefault={onSubmit} class="h-32 flex flex-col items-center justify-evenly">
        <input
          class="h-8 w-56 border pl-2 rounded-sm"
          type="text" 
          name="login-text" 
          id="login-text"
          spellcheck=false
          bind:value={currentUsername}
        />
        <input
          class="w-56 bg-fuchsia-400 h-8 rounded-sm hover:cursor-pointer"
          type="submit"
          name="message-button"
          id="message-button"
        />
      </form>
</div>
