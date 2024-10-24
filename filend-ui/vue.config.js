const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  devServer:{
    port: 9071,
  },
  transpileDependencies: true,
})
