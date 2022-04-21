<template>
    <div class="db-list">
        <el-card>
            <el-button type="primary" icon="plus" @click="editMockData(true)">新建</el-button>
            <el-button :disabled="chooseData == null" @click="editMockData(false)" type="primary" icon="edit">编辑</el-button>
            <el-button :disabled="chooseData == null" @click="deleteMock(chooseId)" type="danger" icon="delete">删除</el-button>
            <div style="float: right">
                <el-select v-model="query.service" placeholder="请选择服务" filterable clearable>
                    <el-option v-for="item in services" :key="item.service" :label="item.service" :value="item.service"> </el-option>
                </el-select>
                <el-input v-model.trim="query.path" placeholder="请求路径" style="width: 200px" class="ml5" clearable></el-input>
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
                <el-table-column prop="service" label="服务" min-width="120"> </el-table-column>
                <el-table-column prop="method" label="method" min-width="80"> </el-table-column>
                <el-table-column prop="path" label="path" min-width="200" show-overflow-tooltip> </el-table-column>
                <el-table-column prop="description" label="描述" min-width="160" show-overflow-tooltip> </el-table-column>
                <el-table-column prop="status" label="状态" min-width="80">
                    <template #default="scope">
                        <el-tooltip :content="scope.row.status == 1 ? '启用' : '禁用'" placement="top">
                            <el-switch
                                v-model="scope.row.status"
                                :active-value="1"
                                :inactive-value="-1"
                                active-color="#13ce66"
                                inactive-color="#ff4949"
                                @change="changeStatus(scope.row)"
                            ></el-switch>
                        </el-tooltip>
                    </template>
                </el-table-column>
                <el-table-column min-width="115" prop="creator" label="创建账号"></el-table-column>
                <el-table-column min-width="170" prop="createTime" label="创建时间">
                    <template #default="scope">
                        {{ $filters.dateFormat(scope.row.createTime) }}
                    </template>
                </el-table-column>
                <el-table-column min-width="115" prop="modifier" label="修改账号"></el-table-column>
                <el-table-column min-width="170" prop="updateTime" label="最后修改时间">
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

        <mock-data-edit @submitSuccess="search" :services="services" v-model:visible="editDialog.visible" v-model:data="editDialog.data" />
    </div>
</template>

<script lang='ts'>
import { toRefs, reactive, onMounted, defineComponent } from 'vue';
import { mockApi } from './api';
import { ElMessage, ElMessageBox } from 'element-plus';
import MockDataEdit from './MockDataEdit.vue';
export default defineComponent({
    name: 'MockDataList',
    components: { MockDataEdit },
    setup() {
        const state = reactive({
            /**
             * 查询条件
             */
            query: {
                path: '',
                service: '',
                pageNum: 1,
                pageSize: 10,
            },
            datas: [],
            total: 0,
            services: [],
            chooseId: 0,
            chooseData: null as any,
            editDialog: {
                visible: false,
                data: null,
            },
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
            state.chooseId = 0;
            state.chooseData = null;
            let res: any = await mockApi.list.request(state.query);
            state.datas = res.list;
            state.total = res.total;
            state.services = await mockApi.serviceList.request();
        };

        const handlePageChange = (curPage: number) => {
            state.query.pageNum = curPage;
            search();
        };

        const editMockData = (isAdd: boolean) => {
            state.editDialog.data = isAdd ? null : state.chooseData;
            state.editDialog.visible = true;
        };

        const changeStatus = async (row: any) => {
            console.log(row);
            await mockApi.create.request({ id: row.id, status: row.status });
        };

        const deleteMock = async (id: number) => {
            ElMessageBox.confirm(`确定删除 [${state.chooseData.method}] 该mock数据?`, '删除提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning',
            })
                .then(async () => {
                    await mockApi.delete.request({ id });
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
            changeStatus,
            handlePageChange,
            editMockData,
            deleteMock,
            valChange,
        };
    },
});
</script>
<style lang="scss">
</style>