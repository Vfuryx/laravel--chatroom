
const chat = {
    state: {
        username: '',
        avatar: '',
        userList: [],
        record: {},
        groupList: [],
        groups: {},
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
        ADD_MY_RECOR: (state, data) => {
            state.record[data.to].push(data)
        },
        SET_GROUPLIST: (state, data) => {
            console.log(data)
            for (let group of data) {
                state.groupList.push(group.name)
                state.groups[group.name] = []
            }
        },
        SET_GROUP: (state, data) => {
            state.groupList.push(data.Content)
            state.groups[data.Content] = []
        },
        SET_ADDGROUP: (state, data) => {
            console.log(data)
            state.groups[data.to].push(data)
        },
    },
    getters: {
        name: state => {
            return state.username
        },
        chatUserList: state => {
            return state.userList
        },
        group: state => {
            return state.group
        },
        groupList: state => {
            return state.groupList
        },
        getRecordBy: (state) => (name) => {
            return state.record[name]
        },
        getGroupBy: (state) => (name) => {
            return state.groups[name]
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
        setGroupList({
            commit
        }, data) {
            commit('SET_GROUPLIST', data)
        },
        creatGroup({
            commit
        }, data) {
            commit('SET_GROUP', data)
        },
        removeUser({
            commit
        }, data) {
            commit('SET_RENOVEUSER', data)
        },
        giveMe({
            commit
        }, data) {
            commit('SET_ADDRECOR', data)
        },
        addMyRecor({
            commit
        }, data) {
            commit('ADD_MY_RECOR', data)
        },
        group({
            commit
        }, data) {
            commit('SET_ADDGROUP', data)
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
        },
    }
}

export default chat