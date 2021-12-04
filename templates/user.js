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
    };
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

    createUser() {
      axios.post("users", this.formData).then((resp) => {
        const result = resp.data;
        if (result.code === 0) {
          this.$message.success("创建成功");
          this.close();
        } else {
          this.$message.error(result.message);
        }
      });
    },
  },
  template: "#user-template",
});
