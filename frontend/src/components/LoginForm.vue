<template>
<div id="LoginBackground">
    <div id="close" @click="$emit('display', 'game')">X</div>
    <div id="LoginContainer">
        <form action="none">
            <label>Username</label>
            <input type="text" :class="{ badInput: unHasError }" v-model="form.username" required="yes" maxlength="24" />
            <label>Password</label>
            <input type="password" :class="{ badInput: pHasError }" v-model="form.password" required="yes" maxlength="255">
            <label v-show="tCreate">Confirm Password</label>
            <input type="password" :class="{ badInput: pHasError }" v-model="form.confirm" required="yes" id="confPass" v-show="tCreate">
            <div id="midRow">
                <div>
                    <input type="checkbox" v-model="form.remember" id="remember">
                    <label id="rememberLabel">Remember me</label>
                </div>
                <div id="noAcc" @click="toggleCreate">Don't have an account?</div>
            </div>
            <button class="submit" v-if="!tCreate" @click="login">Login</button>

            <button class="submit" v-if="tCreate" @click="createAccount">Create</button>
        </form>
        <p id="serverResponse"> {{ serverResponse }}</p>
    </div>
</div>
</template>

<script setup>
import { defineEmits, ref, reactive } from 'vue'
import { useReCaptcha } from 'vue-recaptcha-v3'
const { executeRecaptcha, recaptchaLoaded } = useReCaptcha()
defineEmits(['display'])
const form = reactive({
    username: '',
    password: '',
    confirm: '',
    remember: false
})
const tCreate = ref(false)
const unHasError = ref(false)
const pHasError = ref(false)
const serverResponse = ref('')
function toggleCreate() {
    tCreate.value = !tCreate.value
    if (tCreate.value) {
        document.getElementById("noAcc").innerText = "Already have an account?"
    } else {
        document.getElementById("noAcc").innerText = "Don't have an account?"
    }
}
function loginFieldsReady() {
    let ret = true
    let UnameReg = /^[a-zA-Z0-9-_]{3,24}$/
    if (form.username.search(UnameReg) === -1) {
        unHasError.value = true
        serverResponse.value = "Username must be 3-24 characters, and only contain letters, numbers, - and _"
        ret = false
    } else {
        unHasError.value = false
    }
    if (form.password.length < 8) {
        pHasError.value = true
        serverResponse.value = "Password must be at least 8 characters"
        ret = false
    } else if (form.password.length > 255){
        pHasError.value = true
        serverResponse.value = "Password must be no larger than 255 characters"
        ret = false
    } else {
        pHasError.value = false
    }
    return ret
}
async function createAccount(e) {
    e.preventDefault()
    if (!loginFieldsReady()) {
        return
    }
    if (form.password != form.confirm) {
        pHasError.value = true
        serverResponse.value = "passwords do not match"
        return
    }
    serverResponse.value = ""
    await recaptchaLoaded()
    const token = await executeRecaptcha('login')
    let req = new XMLHttpRequest()
    req.open("POST", "/api/createaccount")
    req.withCredentials = true
    req.send(JSON.stringify({
        username: form.username,
        password: form.password,
        remember: form.remember,
        token: token
    }))
    req.onreadystatechange = () => {
        if (req.readyState == XMLHttpRequest.DONE) {
            let resp = JSON.parse(req.response)
            let err = resp["error"]
            if (err == undefined) {
                location.reload()
            } else {
                document.getElementById("serverResponse").innerText = `Error: ${err}`;
            }
        }
    }
}
function login(e) {
    e.preventDefault()
    if (!loginFieldsReady()) {
        return
    }
    let req = new XMLHttpRequest()
    req.open("POST", "/api/login")
    req.withCredentials = true
    req.send(JSON.stringify({
        username: form.username,
        password: form.password,
        remember: form.remember,
    }))
    req.onreadystatechange = () => {
        if (req.readyState == XMLHttpRequest.DONE) {
            let resp = JSON.parse(req.response)
            let err = resp["error"]
            if (err == undefined) {
                location.reload()
            } else {
                serverResponse.value = `Error: ${err}`;
            }
        }
    }
}

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
    #LoginContainer {
        width: 40%;
        padding: 30px;
        margin: auto;
        margin-top: 5em;
        border-radius: 30px;
        background-color: rgba(0, 30, 61, 1);
    }
    #noAcc {
        font-family: 'Concert One', cursive;
        color:#417fc2;
        font-size: 14px;
    }
    #noAcc:hover {
        cursor: pointer;
        color:#2b4d72;
    }
    #rememberLabel {
        padding-left: 3px;
    }
    #midRow {
        width: 100%;
        display: flex;
        justify-content: space-around;
        align-items:flex-end;
    }
    form {
        display: flex;
        flex-flow: column;
        align-items: flex-start;
    }
    form label {
        font-family: 'Concert One', cursive;
        padding-left: 1em;
    }
    input[type=text], input[type=password] {
        width: 100%;
        padding: 6px;
        margin: 8px 0;
        border-radius: 15px;
        font-family: 'Concert One', cursive;
        font-size: 16px;
    }
    .submit {
        margin-top: 1em;
        border: none;
        width: 30%;
        background: none;
        color: #f3f9f8;
        font-family: 'Fredoka One', cursive;
        font-size: 24px;
        align-self: flex-end;
    }
    .submit:hover {
        cursor: pointer;
        color: #979b9a;
    }
    #serverResponse {
        font-family: 'Times New Roman', Times, serif;
        color: rgb(212, 97, 117);
    }
    .badInput {
        border: 2px solid red;
    }
    @media (max-aspect-ratio:13/16) {
        #LoginContainer {
            width: 80%;
        }
    }
</style>
