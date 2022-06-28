const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true
})

module.exports = {
  devServer: {
    host: "localhost",
    port: 8000,
    proxy: {
      "/": {
        target: "http://localhost:8000",
        secure: false
      }
    }
  }
};