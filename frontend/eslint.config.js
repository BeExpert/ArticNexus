import pluginVue from 'eslint-plugin-vue'

export default [
  {
    ignores: ['dist/**', 'public/**']
  },
  ...pluginVue.configs['flat/vue3-essential']
]
