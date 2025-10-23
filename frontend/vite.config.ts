import { defineConfig } from 'vite';
import preact from '@preact/preset-vite';

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [preact()],
	server: {
		proxy: {
			'/process': {
				target: 'http://localhost:8080',
				changeOrigin: true,
			}
		}
	}
});
