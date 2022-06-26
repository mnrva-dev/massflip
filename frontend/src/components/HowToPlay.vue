<template>
<div id="HowToBG">
    <div id="close" @click="$emit('display', 'game')">X</div>
    <div id="HowTo">
        <p>
            The idea of <span class="logo">Massflip</span> is simple: place a bet on a coinflip. 
            Every player starts with the same amount of GP - the game's currency. The goal is to win
            as much GP as possible without going bankrupt.
            <br>
            The game's clock will tick 
            down to each coin flip, and in that time you can enter a bet into the Betting input, 
            and click on either heads or tails to place your bet. If you win, you get a payout based on
            how much GP was bet against you, and how much you contributed to the winning side's GP.
            <br>
            Compete with your friends to see how much GP you can win without going bankrupt!
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
    #HowTo {
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
    
    @media (max-aspect-ratio:13/16) {
        #HowTo {
            width: 80%;
        }
    }
</style>