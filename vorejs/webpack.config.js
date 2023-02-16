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
  experiments: {
    asyncWebAssembly: true,
    syncWebAssembly: true,
  },
  module: {
    rules: [
      {
        test: /\.go$/,
        use: [
          { 
            loader: path.resolve('go-loader.js'),
            options: {
              name: '[name].[ext]'
            }
          }
        ],
      },
      {
        test: /\.wasm$/,
        type: 'javascript/auto',
        loader: 'file-loader',
        options: {
          name: '[name].[ext]',
        },
      },
      {
        test: /\.ts(x?)$/,
        exclude: [/node_modules/, /test/],
        use: [
          {
            loader: 'ts-loader',
          },
        ],
      },
    ],
  },
  resolve: {
    alias: {
      libvorejs$: '/dist/main.wasm'
    },
    extensions: ['.go', '.wasm', '.tsx', '.ts', '.js'],
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