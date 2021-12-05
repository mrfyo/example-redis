const vm = new Vue({
    el: "#app",
    data() {
      return {
        user: undefined
      }
    },
    methods: {
      updateUser(user) {
          this.user = user
      }
    }
  });