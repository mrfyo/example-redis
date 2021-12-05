Vue.component("article-component", {
    props: {
        user: {
            required: true
        }
    },
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
            personal: false
        }
    },
    created() {
        this.fetchArticles()
    },

    watch: {
        user: function(newVal, oldVal) {
            if(newVal === undefined && !this.personal) {
                this.fetchArticles()
            }
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

        handleSelectMode() {
            const personal = !this.personal
            const user = this.user
            if(personal) {
                this.tableData = this.tableData.filter(item => item.poster === user.username)
            }else {
                this.fetchArticles()
            }
            this.personal = personal
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

        handlePublish(row) {
            const { id } = row
            const user = this.user
            this.$confirm('此操作将发布该文章, 是否继续?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                request.post(`api/articles/publish`, {
                    userId: user.id,
                    articleId: id
                }).then(resp => {
                    const result = resp.data;
                    if (result.code === 0) {
                        this.$message.success("发布成功");
                        this.fetchArticles()
                    } else {
                        this.$message.error(result.message);
                    }
                })
            })
        },

        handleLike(row) {
            const { id } = row
            const user = this.user
            this.$confirm(`Hi ${user.nickname}, 是否为该文章投票?`, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                request.post(`api/articles/vote`, {
                    userId: user.id,
                    articleId: id
                }).then(resp => {
                    const result = resp.data;
                    if (result.code === 0) {
                        this.$message.success("投票成功");
                        this.fetchArticles()
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
