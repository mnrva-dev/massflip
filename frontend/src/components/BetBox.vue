<template>
    <div id="GameWindow">
        <div id="Timer">
            <div id="timeUntilFlip">{{ clock }}</div>
            <div id="until">{{ until }}</div>
        </div>
        <div id="BetGraph">
            <div class="spin">
                <div class="tails-pie pie animate pie-base" :style="tailsStyle"></div>
                <div class="heads-pie pie animate pie-overlap" :style="headsStyle"></div>
            </div>
                <div class="data-overlap" id="data">
                    <div id="headsInfo" class="percent heads" :class="{ unfocused: User.bet == 'tails' }">{{ headsPercent }}%</div>
                    <div id="headsPool" class="pool heads" :class="{ unfocused: User.bet == 'tails' }"> {{ headsPool }}gp</div>
                    <div id="tailsInfo" class="percent tails" :class="{ unfocused: User.bet == 'heads' }">{{ tailsPercent }}%</div>
                    <div id="tailsPool" class="pool tails" :class="{ unfocused: User.bet == 'heads' }">{{ tailsPool }}gp</div>
                </div>
        </div>
        <div id="lower" v-show="userStore().username != ''">
            <BetInput />
            <div id="gpCount">
                Your GP: 
                <div id="gp">
                    {{ userStore().points }}
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { userStore } from '../store.js'
import BetInput from './BetInput.vue'
import { WS } from '../main.js'
import { onMounted,ref,reactive } from 'vue'

const User = userStore()

const clock = ref('0:00s')
const until = ref('until next flip')
const headsPercent = ref(0)
const tailsPercent = ref(0)
const headsPool = ref(0)
const tailsPool = ref(0)

const headsStyle = reactive({
    "--p":50,
    "--c":"#efbe08", //unfocused #7a661b
    "--b":"30px"
})
const tailsStyle = reactive({
    "--p":100,
    "--c":"#dc2f91", //unfocused #6b214b
    "--b":"30px"
})

User.$subscribe(() => {
    if (User.bet == 'heads') {
        tailsStyle['--c'] = '#6b214b'
    } else if (User.bet == 'tails') {
        headsStyle['--c'] = '#7a661b'
    } else {
        headsStyle['--c'] = '#efbe08'
        tailsStyle['--c'] = '#dc2f91'
    }
})

onMounted(() => {
    function updatePool(wsMsg) {
        let hP = Math.floor((wsMsg.headspool / (wsMsg.headspool + wsMsg.tailspool)) * 100)
        if (isNaN(hP)) {
            headsPercent.value = 0
            headsStyle['--p'] = 50
            
        } else {
            headsPercent.value = hP
            headsStyle['--p'] = headsPercent.value
        }
        
        if (isNaN(hP) || wsMsg.tailspool == 0) {
            tailsPercent.value = 0
        } else {
            tailsPercent.value = (100 - hP)
        }
        headsPool.value = wsMsg.headspool
        tailsPool.value = wsMsg.tailspool
    }
    WS.addEventListener("message",  function (evt) {
        let MSGs = evt.data.split('\n')
            MSGs.forEach((i) => {
                let wsMsg = JSON.parse(i)
            if (wsMsg.type == "pool") {
                updatePool(wsMsg)
            } else if (wsMsg.type == "win") {
                userStore().addPoints(wsMsg.value)
            } else if (wsMsg.type == "tick") {
                let time = wsMsg.clock
                let timeString = (Math.floor(time/60)).toString() + ":" + ((time%60)>9?"":"0") + (time%60).toString() + "s"
                clock.value = timeString
                until.value = 'until next flip'
                updatePool(wsMsg)
            } else if (wsMsg.type == "flip") {
                clock.value = wsMsg.value
                until.value = ''
                userStore().setBet("")
            } else if (wsMsg.type == "state") {
                if (wsMsg.value == true) {
                    userStore().bet = wsMsg.bet
                }
            }
        })
        
    })
})

</script>

<style scoped>
#GameWindow {
    display: flex;
    align-items: center;
    flex-direction: column;
    margin: auto;
    justify-content: space-around;
    height: 90%;
}
#BetGraph {
    height: 350px;
    width: 350px;
    background-color: #000814;
    border-radius: 175px;
    user-select: none;
    margin: none;
    padding: 0;
}
#lower {
    height: 200px;
}
@media (max-aspect-ratio:13/16) {
    #GameWindow {
        flex-direction: row;
    }
}

@property --p{
  syntax: '<number>';
  inherits: true;
  initial-value: 0;
}
.pie {
    --p:50;      /* the percentage */
    --b:30px;    /* the thickness */
    --c:#444;  /* the color */
    --w:350px;   /* the size*/

    width:var(--w);
    aspect-ratio:1/1;
    position:fixed;
    display:inline-grid;
    place-content:center;
    font-size:25px;
    font-weight:bold;
    font-family:sans-serif;
    margin: 0px;
    transition: all, .5s;
}
.pie:before,
.pie:after {
  content:"";
  position:absolute;
  border-radius:50%;
}
.heads-pie:before {
  inset:0;
  background:
    conic-gradient(var(--c) calc(var(--p)*1%),#0000 0);
  -webkit-mask:radial-gradient(farthest-side,#0000 calc(99% - var(--b)),#000 calc(100% - var(--b)));
          mask:radial-gradient(farthest-side,#0000 calc(99% - var(--b)),#000 calc(100% - var(--b)));
}
.tails-pie:before {
  inset:0;
  background:
    conic-gradient(var(--c) calc(var(--p)*1%),#0000 0);
  -webkit-mask:radial-gradient(farthest-side,#0000 calc(99% - var(--b)),#000 calc(100% - var(--b)));
          mask:radial-gradient(farthest-side,#0000 calc(99% - var(--b)),#000 calc(100% - var(--b)));
}
.animate {
  animation:p 1s .5s both;
}
@-moz-keyframes p{
  from{--p:0}
}
@keyframes p{
  from{--p:0}
}
.pie-base {
    position: absolute;
}
.pie-overlap {
    position:relative;
}
.data-overlap {
    position:relative;
    bottom:304px;
    left: 50px;
    background-color: #001d3d;
    width: 250px;
    height: 250px;
    border-radius: 125px;
}
#data-bg {
    height: 150px;
    background-color: #001d3d;
}
.percent {
    font-family: 'Fredoka One', cursive;
    font-size: 48px;
    margin: 0px;
}
.pool {
    font-family: 'Fredoka One', cursive;
    font-size: 15px;
    margin-top: 5px;
}
.heads {
    color:#efbe08;
}
.tails {
    color:#dc2f91;
}
.heads.unfocused {
    color:#7a661b;
}
.tails.unfocused {
    color:#6b214b;
}
.pool.heads {
    padding-bottom: 2em;
}
.heads.percent {
    padding-top: .5em;
}
#Timer {
    font-family: 'Fredoka One', cursive;
    font-size: 56px;
}
#until {
    font-size: 18px;
    padding-bottom: 2em;
    height: 20px;
}
#gpCount {
    font-family: 'Concert One', cursive;
    font-size: 24px;
}
#gp {
    font-size: 36px;
    background-image: linear-gradient(0deg, #efbe08, #dc2f91);
    background-size: 100%;
    background-clip: text;
    color: transparent;
}
.spin {
    animation: rotation 20s infinite linear;
}
@keyframes rotation {
    from {
        transform: rotate(0deg);
    }
    to {
        transform: rotate(359deg);
    }
}
</style>