<script setup lang="ts">
import { computed } from 'vue'
import { useDark, useToggle } from '@vueuse/core'
import { useRouter, RouterView } from 'vue-router'
import { Button } from 'primevue'

const router = useRouter()
const isDarkMode = useDark()

const toggleDarkMode = useToggle(isDarkMode)

const menu = computed(() => [
    {
        icon: "home",
        ariaLabel: "Home",
        handler: () => router.push('/')
    },
    {
        icon: "history",
        ariaLabel: "History",
        handler: () => router.push('/history')
    },
    {
        icon: "cog",
        ariaLabel: "Settings",
        handler: () => router.push('/settings')
    },
    {
        icon: isDarkMode.value ? 'sun' : 'moon',
        ariaLabel: isDarkMode.value ? 'Dark mode' : 'Light mode',
        handler: () => toggleDarkMode()
    }
])
</script>

<template>
<div class="dashboard">
    <aside>
        <nav>
            <Button 
                v-for="item in menu"
                :icon="'pi pi-'+item.icon"
                :aria-label="item.ariaLabel"
                @click="item.handler">
            </Button>
        </nav>
    </aside>
    <main class="page">
        <RouterView />
    </main>
</div>
</template>

<style lang="scss" scoped>
main, aside {
    @include p(2);
}

aside {
    height: 100%;

    nav {
        @include flex-col;
        @include p(4);
        height: 100%;
        background: var(--p-slate-800);
        border-radius: 8px;

        button {
            @include mb(4);
        }
    }
}

main {
    width: 100%;
}

.dashboard {
    @include flex;
    @include full;
}

.page {
    @include full;
}

.menubar {
    @include mb(4);
}
</style>