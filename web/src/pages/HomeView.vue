<script setup lang="ts">
import { useAudiobookConnection } from '@/app/server'
import { Button, ProgressBar } from 'primevue'

const { start, stop, step, chapter, progress, error, isConnected, isComplete } = useAudiobookConnection()
</script>

<template>
<div id="home-view">
    <h4 class="muted">Overview</h4>
    <h3>Go Audiobook</h3>
    <div class="flex-center flex-col">
        <Button v-if="!isConnected && !isComplete" @click="start">Start Generation</Button>
        <Button v-if="isConnected || isComplete" @click="stop">Stop</Button>
    
        <div v-if="!isComplete && isConnected">
            <template v-if="step === ''">
                <h3>Starting up the engines! 🚀</h3>
            </template>
            <template v-else>
                <template v-if="step === 'starting'">
                    <h3>Setting up...</h3>
                </template>
                <template v-if="step === 'scanning_epub'">
                    <h3>Scanning EPUB...</h3>
                </template>
                <template v-if="step === 'chapter_creation'">
                    <h3>Processing Chapter: {{ chapter }}</h3>
                </template>
                <template v-if="step === 'building_metadata'">
                    <h3>Building metadata...</h3>
                </template>
                <template v-if="step === 'create_audiobook'">
                    <h3>Generating final audiobook...</h3>
                </template>
                <ProgressBar :value="progress" />
            </template>
        </div>
    
        <p v-if="error">Error: {{ error }}</p>
        <p v-if="isComplete">Audiobook Complete!</p>
    </div>
</div>
</template>

<style scoped lang="scss">

</style>