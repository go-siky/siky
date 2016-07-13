/**
 * Created by liwei on 2016/6/24.
 */
var proxy = require('http-proxy-middleware');

var rApiProxy = proxy('/v2', {
    target: 'http://10.3.15.35:5000',
    changeOrigin: true   // for vhosted sites
});
var gApiProxy = proxy('/api/', {
    target: 'http://localhost:8080',
    changeOrigin: true   // for vhosted sites
});

module.exports = {
    injectChanges: true,
    files: ["./**/*.{html, htm, css, js, ts}", "*.html"],
    watchOptions: { "ignored": "node_modules"},
    server: {
        baseDir: ".",
        middleware: {
            1: rApiProxy,
            2: gApiProxy,
            3: require('connect-history-api-fallback')({index: '/index.html', verbose: false})
        }
    }

};