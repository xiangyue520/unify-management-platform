import Api from '@/common/Api';

export const mockApi = {
    // 获取权限列表
    list: Api.create("/mock-datas", 'get'),
    deatil: Api.create("/mock-datas/detail", 'get'),
    serviceList: Api.create("/mock-datas/services", 'get'),
    create: Api.create("/mock-datas", 'post'),
    update: Api.create("/mock-datas", 'put'),
    delete: Api.create("/mock-datas/{id}", 'delete'),
}