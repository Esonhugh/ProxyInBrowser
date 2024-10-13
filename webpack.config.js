const path = require('path');
const TerserPlugin = require("terser-webpack-plugin");
const WebpackObfuscator = require('webpack-obfuscator');

module.exports = {
  mode: "development",
  devtool: "inline-source-map",
  entry: {
    main: "./src/payload.ts",
  },
  output: {
    path: path.resolve(__dirname, './build'),
    filename: "bundle.js" // <--- Will be compiled to this single file
  },
  resolve: {
    extensions: [".ts", ".tsx", ".js"],
  },
  module: {
    rules: [
      { 
        test: /\.tsx?$/,
        loader: "ts-loader"
      },
    ]
  },
  plugins: [
    new WebpackObfuscator({rotateStringArray: true, reservedStrings: [ '\s*' ]}, [])
  ],
  optimization: {
    minimize: true,
    minimizer: [
        new TerserPlugin(),
    ],
  }
};