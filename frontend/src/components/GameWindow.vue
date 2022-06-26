<template>
    <div id="GameWindow">
        <BetBox />
        <ChatBox />
    </div>
</template>

<script setup>
    import BetBox from './BetBox.vue'
    import ChatBox from './ChatBox.vue'
</script>

<style scoped>
    #GameWindow {
        display: flex;
        justify-content: space-around;
    }
    @media (max-aspect-ratio:13/16) {
        #GameWindow {
            flex-direction: column;
        }
    }
</style>