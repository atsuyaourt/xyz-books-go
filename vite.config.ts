import { fileURLToPath, URL } from 'url'
import { defineConfig } from 'vite'

import Vue from '@vitejs/plugin-vue'

import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import VueRouter from 'unplugin-vue-router/vite'
import Layouts from 'vite-plugin-vue-layouts'

import Unocss from 'unocss/vite'

import { VueRouterAutoImports } from 'unplugin-vue-router'

export default defineConfig({
  root: 'internal/front',
  plugins: [
    VueRouter({ routesFolder: 'internal/front/src/pages', dts: 'internal/front/src/typed-router.d.ts' }),
    Vue(),
    Layouts({
      defaultLayout: 'default',
      layoutsDirs: 'src/layouts',
    }),
    Unocss({
      configFile: 'internal/front/uno.config.ts',
    }),
    AutoImport({
      imports: [
        'vue',
        VueRouterAutoImports,
        {
          axios: [['default', 'axios']],
          '@tanstack/vue-query': ['useQuery', 'useMutation', 'useInfiniteQuery'],
          pinia: ['defineStore', 'storeToRefs'],
        },
      ],
      dts: 'src/auto-imports.d.ts',
      dirs: ['src/composables', 'src/stores'],
      vueTemplate: true,
    }),
    Components({
      dts: 'src/components.d.ts',
      directoryAsNamespace: true,
    }),
  ],
  server: {
    proxy: {
      '/api/v1': 'http://localhost:3000/',
    },
  },
  resolve: {
    alias: [{ find: '@', replacement: fileURLToPath(new URL('./internal/front/src', import.meta.url)) }],
  },
})
