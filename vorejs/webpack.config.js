const path = require('path');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');

const config = {
  target: 'web',
  entry: {
    index: './src/libvore.ts',
  },
  output: {
    path: path.resolve(__dirname, './dist'),
    filename: 'libvore.js',
    library: 'LibVoreJS',
    libraryTarget: 'umd',
    globalObject: 'this',
    umdNamedDefine: true,
  },
  watchOptions: {
    aggregateTimeout: 600,
    ignored: /node_modules/,
  },
  plugins: [
    new CleanWebpackPlugin({
      cleanStaleWebpackAssets: false,
      cleanOnceBeforeBuildPatterns: [path.resolve(__dirname, './dist')],
    }),
  ],
  module: {
    rules: [
      {
        test: /\.ts(x?)$/,
        exclude: [/node_modules/, /test/],
        use: [
          {
            loader: 'ts-loader',
          },
        ],
      },
      {
        test: /\.go$/,
        use: [
          { loader: path.resolve('go-loader.js') }
        ]
      }
    ],
  },
  resolve: {
    extensions: ['.go', '.tsx', '.ts', '.js'],
  },
  resolveLoader: {
    modules: ['node_modules', path.resolve(__dirname)],
  },
};

module.exports = (env, argv) => {
  if (argv.mode === 'development') {
  } else if (argv.mode === 'production') {
  } else {
    throw new Error('Specify env');
  }

  return config;
};