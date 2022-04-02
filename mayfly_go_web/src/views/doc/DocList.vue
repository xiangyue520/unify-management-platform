<template>
    <div class="db-list">
        <el-card>
            <el-button type="primary" icon="plus" @click="editDoc(true)">新建文档</el-button>
            <el-button :disabled="chooseData == null" @click="editDoc(false)" type="primary" icon="edit">编辑</el-button>
            <el-button :disabled="chooseData == null" @click="deleteDoc(chooseId)" type="danger" icon="delete">删除</el-button>
            <div style="float: right">
                <el-select v-model="query.category" placeholder="请选择分类" filterable clearable>
                    <el-option v-for="item in category" :key="item.category" :label="item.category" :value="item.category"> </el-option>
                </el-select>
                <el-input class="ml5" v-model="query.title" style="width: 200px" placeholder="请输入标题" auto-complete="off" clearable></el-input>
                <el-button class="ml5" v-waves type="success" icon="search" @click="search()"></el-button>
            </div>
            <el-table :data="datas" ref="table" @current-change="choose" show-overflow-tooltip>
                <el-table-column label="选择" width="60px">
                    <template #default="scope">
                        <el-radio v-model="chooseId" :label="scope.row.id">
                            <i></i>
                        </el-radio>
                    </template>
                </el-table-column>
                <el-table-column prop="title" label="标题" min-width="240">
                    <template #default="scope">
                        <el-link @click="previewDoc(scope.row)" type="primary" target="_blank">{{ scope.row.title }}</el-link>
                    </template>
                </el-table-column>
                <el-table-column prop="category" label="分类" min-width="100">
                    <template #default="scope">
                        <el-tag color="#E4F5EB" size="small">{{ scope.row.category }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column min-width="115" prop="creator" label="创建账号"></el-table-column>
                <el-table-column min-width="160" prop="createTime" label="创建时间">
                    <template #default="scope">
                        {{ $filters.dateFormat(scope.row.createTime) }}
                    </template>
                </el-table-column>
                <el-table-column min-width="115" prop="modifier" label="修改账号"></el-table-column>
                <el-table-column min-width="160" prop="updateTime" label="最后修改时间">
                    <template #default="scope">
                        {{ $filters.dateFormat(scope.row.updateTime) }}
                    </template>
                </el-table-column>
            </el-table>
            <el-row style="margin-top: 20px" type="flex" justify="end">
                <el-pagination
                    style="text-align: right"
                    @current-change="handlePageChange"
                    :total="total"
                    layout="prev, pager, next, total, jumper"
                    v-model:current-page="query.pageNum"
                    :page-size="query.pageSize"
                ></el-pagination>
            </el-row>
        </el-card>
    </div>
</template>

<script lang='ts'>
import { toRefs, reactive, onMounted, defineComponent } from 'vue';
import { docApi } from './api';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useRouter } from 'vue-router';
export default defineComponent({
    name: 'DocList',
    components: {},
    setup() {
        const router = useRouter();
        const state = reactive({
            /**
             * 查询条件
             */
            query: {
                title: '',
                category: '',
                pageNum: 1,
                pageSize: 10,
            },
            datas: [],
            total: 0,
            category: [],
            chooseId: 0,
            chooseData: null as any,
        });

        onMounted(async () => {
            search();
        });

        const choose = (item: any) => {
            if (!item) {
                return;
            }
            state.chooseId = item.id;
            state.chooseData = item;
        };

        const search = async () => {
            let res: any = await docApi.docs.request(state.query);
            state.datas = res.list;
            state.total = res.total;
            state.category = await docApi.docCategory.request();
        };

        const handlePageChange = (curPage: number) => {
            state.query.pageNum = curPage;
            search();
        };

        const editDoc = (isAdd: boolean) => {
            const { href } = router.resolve({
                path: `/doc/edit`,
                query: {
                    mode: 'editable',
                    id: isAdd ? 0 : state.chooseId,
                    title: isAdd ? '新建文档' : state.chooseData.title,
                },
            });
            window.open(href, '_blank');
        };

        const previewDoc = (row: any) => {
            const { href } = router.resolve({
                path: `/doc/edit`,
                query: {
                    mode: 'preview',
                    id: row.id,
                    title: row.title,
                },
            });
            window.open(href, '_blank');
        };

        const deleteDoc = async (id: number) => {
            ElMessageBox.confirm(`确定删除 [${state.chooseData.title}] 该文档?`, '删除提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning',
            })
                .then(async () => {
                    await docApi.deleteDoc.request({ id });
                    ElMessage.success('删除成功');
                    state.chooseId = 0;
                    state.chooseData = null;
                    search();
                })
                .catch(() => {});
        };

        const valChange = () => {
            search();
        };

        return {
            ...toRefs(state),
            search,
            choose,
            handlePageChange,
            editDoc,
            previewDoc,
            deleteDoc,
            valChange,
        };
    },
});
</script>
<style lang="scss">
</style>
