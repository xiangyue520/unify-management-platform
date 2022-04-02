import Api from '@/common/Api';

export const docApi = {
    docs: Api.create("/docs", 'get'),
    docCategory: Api.create("/docs/categories", 'get'),
    upload: Api.create("/docs/upload", 'post'),
    docDeatil: Api.create("/docs/{id}", 'get'),
    saveDoc: Api.create("/docs", 'post'),
    deleteDoc: Api.create("/docs/{id}", 'delete'),
}