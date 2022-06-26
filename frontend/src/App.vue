<template>
  <div id="Page">
    <MassFlip />
    <PageFooter />
  </div>
</template>

<script setup>
import MassFlip from './components/MassFlip.vue'
import PageFooter from './components/PageFooter.vue'
import { userStore } from './store.js'
import { WSSend } from './main.js'
import { onMounted } from 'vue'

onMounted(() => {
  let cookieStr = document.cookie
  function cookiesToObj(str) {
    str = str.split('; ');
    var result = {};
    for (var i = 0; i < str.length; i++) {
        var cur = str[i].split('=');
        result[cur[0]] = cur[1];
    }
    return result;
  }
  let jar = cookiesToObj(cookieStr)
  if ("session" in jar) {
    let id = jar["session"]
    // handle logged in user
    let req = new XMLHttpRequest
    req.open("POST", "/api/login/bysession")
    req.send(JSON.stringify({
      session: id
    }))
    req.onreadystatechange = () => {
        if (req.readyState == XMLHttpRequest.DONE) {
            let usr = JSON.parse(req.responseText)
            if ("error" in usr) {
              document.cookie = "session=; Max-Age=-99999999"
              console.log(usr["error"])
              return
            }
            userStore().updateUser(usr)
            let msg = JSON.stringify({
              type: "bind",
              username: usr.username
            })
            WSSend(msg)
        }
    }
  }
})
</script>

<style>
#Page {
  height: 100%;
  width: 100%;
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  display: flex;
  flex-direction: column;
}
body {
  background-color: #000814;
}
</style>
