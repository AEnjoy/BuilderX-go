import API from "./axiosInstance.js";

export function DirList(data) {
    return API.post('/api/v1/file/list', data)
}
export function FileSelect(data) {
    return API.post('/api/v1/file/select', data)
}