import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  server: {
    host: '0.0.0.0',
    //port: 8088,
    proxy: {
      '/api': {
        // 后台地址
        target: 'http://localhost:18088/',
        changeOrigin: true,
        secure: false,
        ws: true,
        followRedirects: true,
        rewrite: path => path.replace(/^\/ /, ''),
        onProxyReq(proxyReq, req, res) {
          originHost = req.headers['x-forwarded-for']
          const cookie = req.headers['cookie']
          if (cookie) {
            proxyReq.setHeader('cookie', cookie)
          }
        },
        onProxyRes(proxyRes, req, res) {
          if (proxyRes.headers['set-cookie']) {
            // 域名信息与实际业务相关
            proxyRes.headers['Access-Control-Allow-Credentials'] = 'true';
          }
        }

      },
    }
  },

})
