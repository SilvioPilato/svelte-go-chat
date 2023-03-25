<script lang="ts">
    import type { LoggedUser } from "../types";

    export let users: LoggedUser[] = []
</script>

<div class="flex flex-col p-4 gap-2">
    <h1 class="text-center text-lg">{"Connected users"}</h1>
    <div class="flex flex-col gap-1">
        {#each users as user}
            <div class="text-sm text-right text-slate-50">
                {user.fullUsername}
            </div>
        {/each}
    </div>
</div>