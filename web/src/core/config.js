/**
 * 网站配置文件
 */
import packageInfo from '../../package.json'

const greenText = (text) => `\x1b[32m${text}\x1b[0m`

export const config = {
  appName: 'Gin-Vue-Admin',
  showViteLogo: true,
  KeepAliveTabs: true,
  logs: []
}

export const viteLogo = (env) => {
  if (config.showViteLogo) {
   
  }
}

export default config
