<template>
    <div id="BetInput">
        <input id="betAmount" :class="{ badInput: hasError}" :disabled="userStore().bet != ''" v-model="bet" type="number" placeholder="Place a bet:"/>
        <div id="buttons">
            <button id="betHeads" :class="{ disable: userStore().bet != ''}" @click="submitBet('heads')" type="button">Heads</button>
            <button id="betTails" :class="{ disable: userStore().bet != ''}" @click="submitBet('tails')" type="button">Tails</button>
        </div>
        <div id="err"> {{ error }}</div>
    </div>
</template>

<script setup>
import { userStore } from '../store.js'
import { WSSend } from '../main.js'
import { ref } from 'vue'

const bet = ref()
const hasError = ref(false)
const error = ref('')

function submitBet(HorT){
    if (userStore().bet != '') {
        return
    }
    if (bet.value <= 0 || isNaN(bet.value)) {
        hasError.value = true
        error.value = "Error: bet must be greater than 0"
        return
    } else if (bet.value % 1 != 0) {
        hasError.value = true
        error.value = "Error: bet must be a whole number"
        return
    }
    if (userStore().points - bet.value < 0) {
        hasError.value = true
        error.value = "Error: you cannot bet more gp than you have"
        return
    }
    error.value = ""
    hasError.value = false
    let wsMsg = {
        type: "bet",
        username: userStore().username.toLowerCase(),
        bet: HorT,
        amount: bet.value,
    }
    WSSend(JSON.stringify(wsMsg));
    userStore().subtractPoints(wsMsg.amount)
    bet.value = ""
    userStore().setBet(HorT)
}
</script>

<style scoped>
#BetInput {
    display: flex;
    align-items: center;
    flex-direction: column;
    padding-top: 2em;
    max-width: 150px;
}
input {
    height: 25px;
    width: 100%;
    border-radius: 30px;
    border: 1px solid white;
    color:#000814;
    font-size: 12px;
    font-weight: bold;
    padding-left: 10px;
}
input:focus {
    outline: none;
}
#buttons {
    display: flex;
    justify-content: space-around;
    gap: 10px;
}
button {
    background: none;
    border: none;
    color: white;
    padding: 10px;
    cursor: pointer;
    font-family: 'Concert One', cursive;
    font-size: 22px;
}
#betHeads {
    color: #efbe08;
}
#betTails {
    color: #dc2f91;
}
#betHeads:hover {
    color: #7a661b;
}
#betTails:hover {
    color:#6b214b;
}
#err {
    color:red;
}
#betHeads.disable {
    color: #7a661b;
    cursor: default;
}
#betTails.disable {
    color:#6b214b;
    cursor: default;
}
</style>