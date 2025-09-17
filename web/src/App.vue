<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useWebSocket } from "@vueuse/core"

const { status, data, open, close } = useWebSocket('ws://localhost:8003/generate', {
    immediate: false,
    heartbeat: {
        message: 'ping',
        interval: 30000
    }
})

const progress = ref<{ percentage: number }>()
const error = ref(null)
const isComplete = ref(false)

const isConnected = computed(() => status.value === 'OPEN')

const percentage = computed(() => {
    if (!progress.value?.percentage) {
        return '0%'
    }
    return (progress.value.percentage * 100).toFixed(0) + '%'
})

watch(data, newData => {
    if (!newData) return

    const message = JSON.parse(newData)

    switch (message.type) {
        case "progress":
            progress.value = message
            break
        case "completed":
            isComplete.value = true
            break
        case "error":
            error.value = message.message
            break
    }
})

function startGeneration () {
    if (!isConnected.value) {
        open()
    }
}

function stopGeneration () {
    close()
    progress.value = undefined
    error.value = null
    isComplete.value = false
}
</script>

<template>
<div>
    <h1>Go Audiobook</h1>

    <p>Status: {{ status }}</p>
    <button v-if="!isConnected && !isComplete" @click="startGeneration">Start Generation</button>
    <button v-if="isConnected || isComplete" @click="stopGeneration">Stop</button>

    <div v-if="!isComplete && progress && progress.chapter > 0">
        <h3>{{ progress.message }}</h3>
        <p>Chapter: {{ progress.chapter }}</p>
        <p>Percentage: {{ percentage }}</p>
    </div>

    <p v-if="error" style="color: red;">Error: {{ error }}</p>
    <p v-if="isComplete" style="color: green;">Audiobook Complete!</p>
</div>
</template>
