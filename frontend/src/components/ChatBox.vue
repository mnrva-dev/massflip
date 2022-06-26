<template>
    <div id="ChatBoxContainer">
        <div id="ChatWindow">
        </div>
        <div id="ChatInputContainer">
            <div v-if="userStore().username != ''">
                <iframe name="cca" style="display:none;"></iframe>
                <form action="#" target="cca" autocomplete="off">
                    <input id="ChatInput" v-model="chat" type="text" :placeholder="ChatString" maxlength="336"/>
                    <button id="submit" @click="sendChat">Send</button>
                </form>
            </div>
            <p v-show="userStore().username == ''">You need to be logged in to chat</p>
        </div>
    </div>
</template>

<script setup>
import { userStore } from '../store.js'
import { Queue } from '../main.js'
import { WS,WSSend } from '../main.js'
import { ref,onMounted, computed } from 'vue'


const chatQueue = new Queue()
const ChatString = computed(() => {
    return `Chat as ${userStore().username}:`
})
const chat = ref('')
const ChatColors = {
    green: "limegreen",
    yellow: "gold",
    cyan: "cyan",
    red: "firebrick",
    pink: "fuchsia",
    violet: "violet",
    orange: "orange",
}

const _CHAT_MAX_HISTORY = 75;

onMounted(() => {
    let log = document.getElementById("ChatWindow")
    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }
    WS.addEventListener("message",  function (evt) {
        let wsMsg = JSON.parse(evt.data)
        if (wsMsg.type == "chat") {
            chatQueue.enqueue(wsMsg)
            if (chatQueue.length  >= _CHAT_MAX_HISTORY) {
                chatQueue.dequeue()
            }
            log.innerHTML = ""
            for (let message of Object.values(chatQueue.elements)) {
                var item = document.createElement("div")
                let fromUser = document.createElement("span")
                fromUser.style = `color: ${message.color};`
                fromUser.innerText = message.username
                item.appendChild(fromUser)
                let chatScore = document.createElement("span")
                chatScore.innerText = `(${message.points})`
                chatScore.style = `color: ${message.color};font-family: 'Helvetica';font-size: 12px;`
                item.appendChild(chatScore)
                let chatMsg = document.createTextNode(`: ${message.message}`)
                item.appendChild(chatMsg)
                appendLog(item)
            }
        }
    })
})

function sendChat() {
    if (chat.value.startsWith("/color")) {  
        let newColor = chat.value.split(" ")[1]
        if (newColor in ChatColors) {
            userStore().$patch({color:ChatColors[newColor]})
            chat.value = ''
            let req = new XMLHttpRequest()
            req.open("PUT", "/api/chatcolor")
            req.send(JSON.stringify({
                username: userStore().username,
                color: ChatColors[newColor],
            }))
            return
        } else if (newColor == 'list') {
            // put a chat message that shows all the colors
            chat.value = ''
            return
        }
    }
    let wsSend = {
        type: "chat",
        username: userStore().username,
        color: userStore().color,
        message: chat.value,
        points: userStore().points
    }
    WSSend(JSON.stringify(wsSend))
    chat.value = ''
}
</script>

<style scoped>
    #ChatBoxContainer {
        height: 70vh;
        width: 420px;
        max-width: 50%;
        margin: 2vh;
        border-radius: 30px;
        background-color: #000814;
        text-align: left;
        font-family: 'Concert One', cursive;
    }
    #ChatInputContainer p {
        text-align: center;
    }
    #ChatWindow {
        box-sizing: border-box;
        height: 92%;
        padding: 20px;
        margin: 1vh;
        margin-bottom: 5px;
        overflow-y: scroll;
        overflow-x: hidden;
        scrollbar-color: #001d3d;
        scrollbar-width: thin;
        font-size: 16px;
    }
    #ChatWindow::-webkit-scrollbar {
        background-color: #000814;
    }
    #ChatInput {
        height: 20px;
        width: 75%;
        margin-left: 15px;
        border-radius: 30px;
        border: 1px solid white;
        color:#000814;
        font-size: 12px;
        font-weight: bold;
        padding-left: 10px;
    }
    #submit {
        background-color: #000814;
        color: #f3f9f8;
        margin-left: 10px;
        outline: none;
        border: none;
        font-weight: bold;
    }
    #submit:hover {
        cursor: pointer;
        color:grey;
    }
    #ChatInput:focus {
        outline: none !important;
    }
    
    @media (max-aspect-ratio:13/16) {
        #ChatBoxContainer {
            width:80%;
            max-width: 80%;
            margin: auto;
            margin-bottom: 30px;
        }
        #ChatWindow {
            width:95%;
        }
        #ChatInput {
            width: 65%;
        }
        #ChatInputContainer {
            width: 95%;
            margin: auto;
        }
    }
</style>
