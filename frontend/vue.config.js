const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true
})

module.exports = {
  devServer: {
    host: "localhost",
    proxy: {
      "/": {
        target: "http://localhost:8000",
        secure: false,
        ws: true
      }
    }
  }
};