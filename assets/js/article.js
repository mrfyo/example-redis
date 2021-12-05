Vue.component("article-component", {
    data() {
        return {
            dialogVisible: false,
            formData: {
                poster: undefined,
                title: undefined,
                content: undefined,
            },
            rules: {
                poster: [{
                    required: true,
                    message: '请输入标题',
                    trigger: 'blur'
                }],
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
            loading: false,
            tableData: [],
        }
    },
    created() {
        this.fetchArticles()
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

        handleRemove(row) {
            const { id } = row
            this.$confirm('此操作将永久删除该文章, 是否继续?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                request.delete(`api/articles/${id}`).then(resp => {
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

        createArticle() {
            request.post("api/articles", this.formData).then(resp => {
                const result = resp.data
                if (result.code === 0) {
                    this.$message.success("创建成功")
                    this.fetchArticles()
                    this.close()
                } else {
                    this.$message.error(result.message)
                }
            })
        },

        fetchArticles() {
            this.loading = true
            request.get("api/articles", {
                params: {
                    offset: 0,
                    limit: 10
                }
            }).then(resp => {
                const result = resp.data
                if (result.code === 0) {
                    const { items } = result.data
                    this.tableData = items
                } else {
                    this.$message.error(result.message)
                }
            }).finally(() => {
                this.loading = false
            })
        }
    },
    template: '#article-template',
});
