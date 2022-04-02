<template>
    <div>
        <v-md-editor
            :mode="mode"
            :disabled-menus="[]"
            @save="setDocInfo"
            @upload-image="uploadImage"
            v-model="doc.content"
            :height="height + 'px'"
        ></v-md-editor>
        <el-dialog title="简介" v-model="docInfoDialog.visible" :show-close="false" width="400px">
            <el-form :model="doc" label-width="55px">
                <el-form-item prop="title" label="标题:" required>
                    <el-input v-model.trim="doc.title" placeholder="请输入文档标题" auto-complete="off"></el-input>
                </el-form-item>
                <el-form-item prop="category" label="分类:" required>
                    <el-select style="width: 100%" v-model="doc.category" filterable allow-create placeholder="请选择文档分类">
                        <el-option v-for="item in category" :key="item.category" :label="item.category" :value="item.category"> </el-option>
                    </el-select>
                </el-form-item>
            </el-form>

            <template #footer>
                <div class="dialog-footer">
                    <el-button type="primary" @click="btnOk">确 定</el-button>
                    <el-button @click="docInfoDialog.visible = false">取 消</el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>

<script lang="ts">
import { onMounted, toRefs, reactive, defineComponent } from 'vue';
import { ElMessage } from 'element-plus';
import { docApi } from './api';
import { useRoute } from 'vue-router';
import { isTrue } from '@/common/assert';

export default defineComponent({
    name: 'DocEdit',
    props: {},
    setup() {
        const route = useRoute();
        const state = reactive({
            mode: 'editable',
            doc: {
                id: 0,
                title: '',
                category: '',
                content: '',
            },
            docInfoDialog: {
                visible: false,
            },
            category: [],
            height: 700,
        });

        const setDocInfo = () => {
            state.docInfoDialog.visible = true;
        };

        const uploadImage = async (event: any, insertImage: any, files: any) => {
            console.log(event);
            console.log(insertImage);
            console.log(files[0]);
            isTrue(files.length == 1, '只支持单个文件上传');
            const file = files[0];

            let param = new FormData(); // 创建form对象
            param.append('file', file); // 通过append向form对象添加数据
            const path = await docApi.upload.requestWithHeaders(param, { 'Content-Type': 'multipart/form-data' });
            ElMessage.success('文件上传成功');
            insertImage({
                url: `${(window as any).globalConfig.BaseApiUrl}${path}`,
                desc: file.name,
            });
        };

        onMounted(async () => {
            state.height = window.innerHeight + 5;
            state.mode = route.query.mode as string;
            const docId = Number.parseInt(route.query.id as string);
            state.category = await docApi.docCategory.request();
            if (!docId) {
                return;
            }
            state.doc = await docApi.docDeatil.request({ id: docId });
        });

        const btnOk = async () => {
            await docApi.saveDoc.request(state.doc);
            ElMessage.success('保存成功');
            state.docInfoDialog.visible = false;
        };

        return {
            ...toRefs(state),
            setDocInfo,
            uploadImage,
            btnOk,
        };
    },
});
</script>
<style lang="scss">
</style>
