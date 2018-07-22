module.exports = {
    output: {
        path: `${__dirname}/public`,
        filename: 'bundle.js',
    },
    devServer: {
        contentBase: `${__dirname}/public`,
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
    resolve: {
        alias: {
            src: `${__dirname}/src`,
            math: `${__dirname}/src/math`,
        },
    },
}
