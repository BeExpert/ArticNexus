import api from './api'

export const companyService = {
  async getAll(page = 1, pageSize = 200) {
    const res = await api.get('/companies', { params: { page, pageSize } })
    return res.data
  },

  // Returns only companies assigned to the current authenticated user.
  async getMyCompanies() {
    const res = await api.get('/auth/me/companies')
    return res.data
  },

  async getById(id) {
    const res = await api.get(`/companies/${id}`)
    return res.data
  },

  async create(data) {
    const res = await api.post('/companies', data)
    return res.data
  },

  async update(id, data) {
    const res = await api.put(`/companies/${id}`, data)
    return res.data
  },

  async remove(id) {
    await api.delete(`/companies/${id}`)
  },
}
