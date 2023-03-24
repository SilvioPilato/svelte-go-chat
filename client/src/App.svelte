<script lang="ts">
    import { onMount } from "svelte";
    import Chat from "./lib/Chat.svelte";
    import Login from "./lib/Login.svelte";
    import type { WSMessage } from "./types";
    let loggedIn = false;
    let messages: string[] = [];
    let ws: WebSocket;
    let loggedUsers = [];
    const fetchLoggedInUsers = async () => {
        try {
            let res = await fetch("http://localhost:8080/users");
            loggedUsers = await res.json();
        } catch (error) {
            console.error(error)
        }
    }
    onMount(() => {
      fetchLoggedInUsers();
      ws = new WebSocket("ws://localhost:8080/ws");
      ws.addEventListener('message', (event) => {
        let parsed = JSON.parse(event.data);
        handleWSMessage(parsed as WSMessage)
      });
    });
    const handleWSMessage = (data: WSMessage) => {
      switch (data.topic) {
        case "login_success":
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
    const onChatSubmit = (message: string) => {
      if (!ws) return;
      let wsMessage = {
        topic: "chat",
        message: message
      }
      ws.send(JSON.stringify(wsMessage));
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
    <Chat onChatSubmit={onChatSubmit} messages={messages}/>
    {:else}
      <Login onLoginSubmit={onLoginSubmit} />
    {/if}
  </main>