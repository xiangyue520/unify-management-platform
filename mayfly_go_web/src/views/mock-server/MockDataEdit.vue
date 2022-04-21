<template>
    <div class="mock-data-dialog">
        <el-dialog
            :title="isAdd ? '新建mock数据' : '修改mock数据'"
            v-model="dialogVisible"
            :close-on-click-modal="false"
            :before-close="cancel"
            :show-close="true"
            :destroy-on-close="true"
            width="900px"
        >
            <el-form :model="form" ref="mockDataForm" :rules="rules" label-width="90px">
                <el-form-item prop="service" label="服务:" required>
                    <el-select style="width: 100%" v-model="form.service" filterable allow-create placeholder="请选择或输入新建服务">
                        <el-option v-for="item in services" :key="item.service" :label="item.service" :value="item.service"> </el-option>
                    </el-select>
                </el-form-item>

                <el-form-item prop="realUrl" label="目标地址:">
                    <el-input v-model.trim="form.realUrl" placeholder="请输入目标地址,禁用时转发至目标地址" clearable></el-input>
                </el-form-item>

                <el-form-item prop="method" label="method:" required>
                    <el-select style="width: 100%" v-model="form.method" placeholder="请选择请求方法">
                        <el-option label="GET" value="GET"> </el-option>
                        <el-option label="POST" value="POST"> </el-option>
                        <el-option label="PUT" value="PUT"> </el-option>
                        <el-option label="DELETE" value="DELETE"> </el-option>
                    </el-select>
                </el-form-item>

                <el-form-item prop="path" label="path:">
                    <el-input v-model.trim="form.path" placeholder="请输入请求路径"></el-input>
                </el-form-item>

                <el-form-item prop="description" label="描述:">
                    <el-input v-model.trim="form.description" placeholder="请输入描述"></el-input>
                </el-form-item>

                <el-form-item prop="data" label="数据" id="content" required>
                    <codemirror :canChangeMode="false" ref="cmEditor" v-model="form.data" language="json" width="700px" />
                </el-form-item>
            </el-form>

            <template #footer>
                <div class="dialog-footer">
                    <el-button type="primary" :loading="btnLoading" @click="btnOk" size="small" :disabled="submitDisabled">保 存</el-button>
                    <el-button @click="cancel()" :disabled="submitDisabled" size="small">关 闭</el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>

<script lang="ts">
import { ref, toRefs, reactive, watch, defineComponent } from 'vue';
import { ElMessage } from 'element-plus';
import { mockApi } from './api';

import { codemirror } from '@/components/codemirror';

export default defineComponent({
    name: 'MockDataEdit',
    components: {
        codemirror,
    },
    props: {
        visible: {
            type: Boolean,
        },
        data: {
            type: Object,
        },
        services: {
            type: Array,
        },
    },
    setup(props: any, { emit }) {
        const { services } = toRefs(props);
        const mockDataForm: any = ref(null);

        const state = reactive({
            dialogVisible: false,
            submitDisabled: false,
            title: '新增数据',
            isAdd: true,
            services: [] as any,
            form: {
                id: null,
                service: '',
                method: '',
                description: '',
                path: '',
                data: '',
                status: 1,
            },
            btnLoading: false,
            rules: {
                service: [
                    {
                        required: true,
                        message: '请选择或输入新建服务',
                        trigger: ['change', 'blur'],
                    },
                ],
                method: [
                    {
                        required: true,
                        message: '请选择method',
                        trigger: ['change', 'blur'],
                    },
                ],
                path: [
                    {
                        required: true,
                        message: '请输入请求路径',
                        trigger: ['change', 'blur'],
                    },
                ],
                status: [
                    {
                        required: true,
                        message: '请选择状态',
                        trigger: ['change', 'blur'],
                    },
                ],
            },
        });

        watch(props, async (newValue) => {
            if (newValue.data) {
                state.form = await mockApi.deatil.request({ id: newValue.data.id });
            } else {
                state.form = {} as any;
                state.form.data = '';
                state.form.status = 1;
            }
            if (state.form.id) {
                state.isAdd = false;
            } else {
                state.isAdd = true;
            }
            state.services = services;
            state.dialogVisible = newValue.visible;
        });

        const btnOk = () => {
            mockDataForm.value.validate((valid: any) => {
                if (valid) {
                    mockApi.create.request(state.form).then(
                        () => {
                            ElMessage.success('保存成功');
                            emit('submitSuccess');
                            state.submitDisabled = false;
                            cancel();
                        },
                        () => {
                            state.submitDisabled = false;
                        }
                    );
                } else {
                    return false;
                }
            });
        };

        const cancel = () => {
            emit('update:visible', false);
            emit('cancel');
        };

        return {
            ...toRefs(state),
            mockDataForm,
            btnOk,
            cancel,
        };
    },
});
</script>
<style lang="scss">
// 	.m-dialog {
// 		.el-cascader {
// 			width: 100%;
// 		}
// 	}
#content {
    .CodeMirror {
        height: 400px !important;
    }
}
</style>
