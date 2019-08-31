const chat = {
    state: {
        username: '',
        avatar: '',
        userList: [],
        record: {}
    },
    mutations: {
        SET_USERNAME: (state, username) => {
            state.username = username
        },
        SET_AVATAR: (state, avatar) => {
            state.avatar = avatar
        },
        SET_USERLIST: (state, userList) => {
            state.userList = userList
            for (let user of userList) {
                state.record[user.name] = []
            }
        },
        SET_ADDUSER: (state, user) => {
            state.userList.push(user)
        },
        SET_RENOVEUSER: (state, name) => {
            let k = 0
            for(let index in state.userList){
                if(state.userList[index].name == name){
                    k = index
                    break
                }
            }

            state.userList.splice(k, 1)
        },
        SET_RECOR: (state, data) => {
            state.record[data] = []
        },
        SET_ADDRECOR: (state, data) => {
            state.record[data.from].push(data)
        },
    },
    getters: {
        name: state => {
            return state.username
        },
        chatUserList: state => {
            return state.userList
        },
        getRecordBy: (state) => (name) => {
            return state.record[name]
        },
        getRecord: (state) => {
            return state.record
        },

    },
    actions: {
        setList({
            commit
        }, data) {
            commit('SET_USERLIST', data)
        },
        removeUser({
            commit
        }, data) {
            commit('SET_RENOVEUSER', data)
        },
        giveMe({
            commit
        }, data) {

            console.log(4545444545445454)
            console.log(data)
            commit('SET_ADDRECOR', data)
        },
        addUser({
            commit
        }, data) {
            commit('SET_ADDUSER', data)
        },
        addRecord({
            commit
        }, data) {
            commit('SET_RECOR', data)

        }
    }
}

export default chat