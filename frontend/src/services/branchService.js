import api from './api'

export const branchService = {
  async getAll(companyId, page = 1, pageSize = 50) {
    const res = await api.get(`/companies/${companyId}/branches`, { params: { page, pageSize } })
    return res.data
  },

  async create(companyId, data) {
    const res = await api.post(`/companies/${companyId}/branches`, data)
    return res.data
  },

  async update(companyId, branchId, data) {
    const res = await api.put(`/companies/${companyId}/branches/${branchId}`, data)
    return res.data
  },

  async remove(companyId, branchId) {
    await api.delete(`/companies/${companyId}/branches/${branchId}`)
  },
}
