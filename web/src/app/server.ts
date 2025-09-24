import { ref, computed, watch } from 'vue'
import { useWebSocket } from '@vueuse/core'

const WEB_SOCKET_URL = 'ws://localhost:8003/generate'

// Progress updates come from server messages; see watcher below for shape.

enum ServerConnectionStatus {
    Open = "OPEN",
    Connecting = "CONNECTING",
    Closed = "CLOSED"
}

type Step = 'starting' | 'scanning_epub' | 'chapter_creation' | 'building_metadata' | 'create_audiobook' | ''

export function useAudiobookConnection() {
    const { status, data, open, close } = useWebSocket(WEB_SOCKET_URL, {
        immediate: false,
        heartbeat: {
            message: 'ping',
            interval: 30000
        }
    })
    
    const progress = ref(0)
    const chapter = ref(0)
    const step = ref<Step>()
    const error = ref(null)
    const isComplete = ref(false)

    const isConnected = computed(() => status.value === ServerConnectionStatus.Open)

    // const percentage = computed(() => {
    //     return progress.value.toFixed(0) + '%'
    // })

    watch(data, (newData: string) => {
        if (!newData) return

        const message = JSON.parse(newData)

        switch (message.type) {
            case "progress":
                step.value = message.step
                if (typeof message.chapter === 'number') {
                    chapter.value = message.chapter
                }
                if (typeof message.percentage === 'number') {
                    progress.value = Math.round(message.percentage * 100)
                }
                break
            case "complete":
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
        chapter.value = 0
        progress.value = 0
        error.value = null
        isComplete.value = false
    }

    return {
        progress,
        chapter,
        step,
        error,
        isComplete,
        isConnected,
        start: startGeneration,
        stop: stopGeneration
    }
}