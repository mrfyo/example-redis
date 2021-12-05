Vue.component("article-component", {
    data() {
        return {
            dialogVisible: false,
            formData: {
                title: undefined,
                content: undefined,
            },
            rules: {
                title: [{
                    required: true,
                    message: '请输入标题',
                    trigger: 'blur'
                }],
                content: [{
                    required: true,
                    message: '请输入内容',
                    trigger: 'blur'
                }],
            },
        }
    },
    methods: {
        onOpen() {
        },
        onClose() {
            this.$refs['elForm'].resetFields()
        },
        close() {
            this.dialogVisible = false
        },
        handleOpen() {
            this.dialogVisible = true
        },
        handelConfirm() {
            this.$refs['elForm'].validate(valid => {
                if (!valid) return

                this.createArticle()
            })
        },

        createArticle() {
            axios.post("articles", this.formData).then(resp => {
                const result = resp.data
                if (result.code === 0) {
                    this.$message.success("创建成功")
                    this.close()
                } else {
                    this.$message.error(result.message)
                }
            })
        }
    },
    template: '#article-template',
});
