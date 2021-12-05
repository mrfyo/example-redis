Vue.component("user-component", {
  data() {
    return {
      dialogVisible: false,
      formData: {
        username: undefined,
        password: undefined,
        nickname: undefined,
      },
      rules: {
        username: [
          {
            required: true,
            message: "请输入用户名",
            trigger: "blur",
          },
        ],
        password: [
          {
            required: true,
            message: "请输入密码",
            trigger: "blur",
          },
        ],
        nickname: [
          {
            required: true,
            message: "请输入昵称",
            trigger: "blur",
          },
        ],
      },
      tableData: [],
      loading: false
    };
  },
  created() {
    this.fetchUsers()
  },
  methods: {
    onOpen() {},
    onClose() {
      this.$refs["elForm"].resetFields();
    },
    close() {
      this.dialogVisible = false;
    },
    handleOpen() {
      this.dialogVisible = true;
    },
    handelConfirm() {
      this.$refs["elForm"].validate((valid) => {
        if (!valid) return;

        this.createUser();
      });
    },

    handleRemove(row) {
      const {id} = row
      this.$confirm('此操作将永久删除该用户, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        request.delete(`api/users/${id}`).then(resp => {
          const result = resp.data;
          if (result.code === 0) {
            this.$message.success("删除成功");
            this.tableData = this.tableData.filter(item => item.id != id)
          } else {
            this.$message.error(result.message);
          }
        })
      })
    },

    createUser() {
      request.post("api/users", this.formData).then((resp) => {
        const result = resp.data;
        if (result.code === 0) {
          this.$message.success("创建成功");
          this.fetchUsers()
          this.close();
        } else {
          this.$message.error(result.message);
        }
      });
    },

    fetchUsers() {
      this.loading = true
      request.get("api/users", {
        params: {
          offset: 0,
          limit: 10
         }
      }).then(resp => {
        const result = resp.data
        if(result.code === 0) {
          const {items} = result.data
          this.tableData = items
        }
      }).finally(() => {
        this.loading = false
      })
    }
  },
  template: "#user-template",
});
