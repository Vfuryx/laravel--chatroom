const websocket = {
    state: {
        ws: '',
        name: '',
        pwd: '',
        url: "ws://" + location.host + "/ws",
        timeoutObj: '',
        login: false,
        lockReconnect: false
    },
    mutations: {
        SET_WS: (state, ws) => {
            state.ws = ws
        },
        SET_NAME: (state, name) => {
            state.name = name
        },
        SET_PWD: (state, pwd) => {
            state.pwd = pwd
        },
        SET_TIMEOUTOBJ: (state, timeoutObj) => {
            state.timeoutObj = timeoutObj
        },
        SET_LOGIN: (state, is_login) => {
            state.login = is_login
        }
    },
    getters: {
        isLogin: state => {
            return state.login
        },
        ws: state => {
            return state.ws
        },
        timeoutObj: state => {
            return state.timeoutObj
        },
        user: state => {
            return state.name
        }
    },
    actions: {
        InitWebSocket({
            commit,
            state
        }, data) {
            commit('SET_NAME', data.username)
            commit('SET_PWD', data.password)
            commit('SET_LOGIN', true)
            let websock = new WebSocket(
                state.url + "?name=" + data.username + "&password=" + data.password
            );
            commit('SET_WS', websock)
        },
        SendMsg(data) {
            this.state.websocket.ws.send(data)
        },
        SetTimeout({
            commit
        }, data) {
            commit('SET_TIMEOUTOBJ', data)
        }
    }
}

export default websocket