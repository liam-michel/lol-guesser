import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import dotenv from 'dotenv'
import path from 'path'
import { fileURLToPath } from 'url'
import {dirname} from 'path'
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

dotenv.config({path: path.resolve(__dirname, '../.env')})

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server:{
    port: parseInt(process.env.REACT_PORT || '3000')
  }
})
