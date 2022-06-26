import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import VueGtag from "vue-gtag";

// switch these in production
export const WS = new WebSocket("wss://" + "massflip.mnrva.dev" + "/ws")
//export const WS = new WebSocket("ws://" + "127.0.0.1:8000" + "/ws")

WS.onclose = function() {
  alert("WebSocket connection closed.")
}
WS.onerror = function() {
  alert("WebSocket connection error.")
}

export function WSSend(msg){
  WSWait(WS, function(){
      WS.send(msg);
  });
}
function WSWait(socket, callback){
  setTimeout(
      function () {
          if (socket.readyState === WebSocket.OPEN) {
              if (callback != null){
                  callback();
              }
          } else {
              WSWait(socket, callback);
          }

      }, 5);
}

export class Queue {
  constructor() {
    this.elements = {};
    this.head = 0;
    this.tail = 0;
  }
  enqueue(element) {
    this.elements[this.tail] = element;
    this.tail++;
  }
  dequeue() {
    const item = this.elements[this.head];
    delete this.elements[this.head];
    this.head++;
    return item;
  }
  peek() {
    return this.elements[this.head];
  }
  get length() {
    return this.tail - this.head;
  }
  get isEmpty() {
    return this.length === 0;
  }
}

const pinia = createPinia()

var app = createApp(App)
app.use(pinia)
app.use(VueGtag, {config: { id: "G-C3WQH98SZB" }})
app.mount('#app')

