<template>
  <div id="Content">
    <PageHeader @display="displayHandle"/>
    <div id="Container">
      <GameWindow />
      <LoginForm @display="displayHandle" class="dullBG" v-show="display.content == 'login'"/>
      <AboutPage @display="displayHandle" class="dullBG" v-show="display.content == 'about'"/>
      <HowToPlay @display="displayHandle" class="dullBG" v-show="display.content == 'howTo'"/>
    </div>
  </div>
</template>

<script setup>
import PageHeader from './PageHeader.vue'
import LoginForm from './LoginForm.vue'
import GameWindow from './GameWindow.vue'
import AboutPage from './AboutPage.vue'
import HowToPlay from './HowToPlay.vue'
import { reactive } from 'vue'

const display = reactive({ content: 'game' })

function displayHandle(c) {
  if (display.content == c) {
    display.content = 'game'
  } else {
    display.content = c
  }
}

</script>

<style scoped>
#Container {
    width: 1000px;
    max-width: 100%;
    height: 75vh;
    border-radius: 30px;
    margin: auto;
    color: #f3f9f8;
    background-color: #001d3d;
    position: relative;
}
.dullBG {
    height: 100%;
    width: 100%;
    position: absolute;
    top: 0px;
    left: 0px;
    border-radius: 30px;
    background-color: rgba(0, 8, 20,0.7);
}
    
@media (max-aspect-ratio:13/16) {
    #Container {
        width: 95%;
        height: auto;
    }
}
</style>