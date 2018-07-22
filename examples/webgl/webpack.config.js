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
        }]
    },
}
