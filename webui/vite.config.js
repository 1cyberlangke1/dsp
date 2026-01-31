import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
    plugins: [
        react(),
    ],
    server: {
        port: 5173,
        proxy: {
            '/admin': {
                target: 'http://localhost:5001',
                changeOrigin: true,
            },
            '/v1': {
                target: 'http://localhost:5001',
                changeOrigin: true,
            },
        },
    },
    build: {
        outDir: '../static/admin',
        emptyOutDir: true,
    },
})
