import { defineStore } from 'pinia'

// useStore could be anything like useUser, useCart
// the first argument is a unique id of the store across your application
export const userStore = defineStore({
    id:'main',
    state: () => {
        return {
            username: '',
            color: '',
            points: 0,
            resets: 0,
            bet: "",
        }
    },
    getters: {
        getUser() {
            return {
                username: this.username,
                color: this.color,
                points: this.points,
                resets: this.resets,
            }
        },
        ready() {
            if (this.username != '') {
                return true
            } else {
                return false
            }
        },
        getBet() {
            return this.bet
        },
        getUsername() {
            return this.username
        }
    },
    actions: {
        updateUser(usr) {
            this.$state = {
                username: usr.username,
                color: usr.color,
                points: usr.points,
                resets: usr.resets,
            }
        },
        addPoints(amt) {
            this.$patch({points: this.points + amt})
        },
        subtractPoints(amt) {
            this.$patch({points: this.points - amt})
        },
        toggle() {
            this.$patch({t:!this.t})
        },
        setBet(bet) {
            this.$patch({bet: bet})
        }
    }
})
