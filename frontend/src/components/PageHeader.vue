<template>
    <div id="HeaderContainer">
        <div id="header">Massflip<span id="beta">alpha</span></div>
        <nav>
            <p @click="$emit('display', 'howTo')">How to play</p>
            <p @click="$emit('display', 'about')">About</p>
            <p v-if="userStore().username == ''" @click="$emit('display', 'login')" id="showLoginNav">Login</p>
            <p v-if="userStore().username != ''" @click="logout()">Logout</p>
        </nav>
    </div>
</template>

<script setup>
import { userStore } from '../store.js'
import { defineEmits } from 'vue'
defineEmits(['display'])
function logout() {
    let req = new XMLHttpRequest()
    req.open("POST", "/api/logout")
    req.send(JSON.stringify({
        username: userStore().username
    }))
    req.onreadystatechange = () => {
        if (req.readyState == XMLHttpRequest.DONE) {
            document.cookie = "session=; Max-Age=-99999999"
            location.reload()
        }
    }
}
</script>

<style scoped>
    #HeaderContainer {
        color: #f3f9f8;
        display: flex;
        height: 125px;
        flex-direction: column;
        align-items: center;
        user-select: none;
        font-family: 'Fredoka One', cursive;
        font-size: 48px;
        margin: 0;
        margin-top: 10px;
    }
    #header {
        display: flex;
    }
    #beta {
        font-size: 24px;
        background-image: linear-gradient(0deg, #efbe08 60%, #dc2f91);
        background-size: 100%;
        background-clip: text;
        color: transparent;
    }
    nav {
        display: flex;
        align-items: center;
        font-family: 'Concert One', cursive;
        width: 400px;
        justify-content: space-between;
        font-size: 18px;
    }
    nav p:hover {
        cursor: pointer;
        color: #bbb;
    }
    @media (max-aspect-ratio:13/16) {
        nav {
            max-width: 80%;

        }
    }
</style>