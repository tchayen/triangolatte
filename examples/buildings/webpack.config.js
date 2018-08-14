const webpack = require('webpack')

module.exports = {
  output: {
    path: `${__dirname}`,
    filename: 'bundle.js',
  },
  module: {
    rules: [{
      test: /\.js$/,
      exclude: /node_modules/,
      use: { loader: 'babel-loader' }
    }, {
      test: /\.scss$/,
      exclude: /node_modules/,
      use: [
        { loader: 'style-loader' },
        { loader: 'css-loader' },
        { loader: 'sass-loader' },
      ]
    }, {
      test: /\.glsl$/,
      use: { loader: 'webpack-glsl-loader' },
    }],
  },
  plugins: [
    new webpack.DefinePlugin({
      __SERVER__: JSON.stringify(process.env.NODE_ENV === 'production'
        ? 'http://178.62.254.13'
        : 'http://localhost:3010'
      )
    })
  ]
}
