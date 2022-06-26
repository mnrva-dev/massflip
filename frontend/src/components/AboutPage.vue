<template>
<div id="AboutBG">
    <div id="close" @click="$emit('display', 'game')">X</div>
    <div id="About">
        <p>
            <span class="logo">Massflip</span> was created by MNRVA as a way to learn about frontend
            frameworks, websockets, APIs, and databases. The goal for <span class="logo">Massflip</span>
            is to eventually become an open source coin-flip-game web server, with integration in twitch
            streams, discord bots, etc., and also to be a project to put on my resume :)
            <br>
            <br>
            If you enjoy this site, share it with your friends! This game thrives on the community chatting
            and betting against each other, so more players are always welcome.
        </p>
    </div>
</div>
</template>

<script setup>
import { defineEmits } from 'vue'
defineEmits(['display'])
</script>

<style scoped>
    #close {
        color: #ccccd1;
        font-size: 26px;
        position: absolute;
        top: 30px;
        right: 50px;
        font-family: 'Fredoka One';
    }
    #close:hover {
        cursor: pointer;
        color: #888;
    }
    #About {
        width: 40%;
        padding: 30px;
        margin: auto;
        margin-top: 5em;
        border-radius: 30px;
        background-color: rgba(0, 30, 61, 1);
    }
    p {
        font-family: 'Concert One', cursive;
        font-size: 18px;
        color: #f3f9f8;
    }
    .logo {
        font-family: 'Fredoka One', cursive;
        font-size: 20px;
        background-image: linear-gradient(0deg, #efbe08, #dc2f91);
        background-size: 100%;
        background-clip: text;
        color: transparent;
    }
</style>