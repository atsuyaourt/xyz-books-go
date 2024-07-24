import { defineConfig, presetIcons, presetUno } from 'unocss'

export default defineConfig({
  content: {
    filesystem: ['internal/views/**/*.templ'],
  },
  presets: [presetUno(), presetIcons()],
})
