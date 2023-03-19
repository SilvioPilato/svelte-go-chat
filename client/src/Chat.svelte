<script lang="ts">
  import { onMount } from "svelte";
    import Login from "./lib/Login.svelte";
  type WSMessage = {
    topic: "login" | "error" | "chat"
    message: string
    username?: string
  }
  let loggedIn = false;
  let messages: string[] = [];
  let ws: WebSocket;
  let currentMessage = "";
  onMount(() => {
    ws = new WebSocket("ws://localhost:8080/ws");
    ws.addEventListener('message', (event) => {
      let parsed = JSON.parse(event.data);
      handleWSMessage(parsed as WSMessage)
    });
  });
  const handleWSMessage = (data: WSMessage) => {
    switch (data.topic) {
      case "login":
        console.log("logged in")
        loggedIn = true;
        break;
      case "chat":
        let message = `${data.username}: ${data.message}`;
        messages = [...messages, message];
        break;
      case "error":
        console.error(data.message);
        break;
      default:
        console.info(data);
        break;
    }
  }
  const onChatSubmit = () => {
    if (!ws) return;
    let message = {
      topic: "chat",
      message: currentMessage
    }
    ws.send(JSON.stringify(message));
    currentMessage = "";
  }
  const onLoginSubmit = (username: string) => {
    if (!ws) return;
    let message = {
      topic: "login",
      message: username
    }
    ws.send(JSON.stringify(message));
  }
</script>

<main class="flex flex-col h-screen pt-4 bg-slate-200 justify-center items-center">
  {#if loggedIn}
  <div class="mx-auto h-4/5">  
    {#each messages as message}
      <p>{message}</p>
    {/each}
  </div>
  <form on:submit|preventDefault={onChatSubmit} class="flex justify-center">
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
  {:else}
    <Login onLoginSubmit={onLoginSubmit} />
  {/if}
</main>

<style>

</style>