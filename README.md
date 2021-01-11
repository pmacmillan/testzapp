How to use:

1. Copy `ssl.dev-coursera.*` from `~/base/coursera/web/config/serve`
2. Configure Oauth ZAPP in dev marketplace: `https://mock.dev-coursera.org:9090/redirect`
3. Configure ZAPP in marketplace to connect to: `https://mock.dev-coursera.org:9090/`
4. Copy `config.sample.json` to `config.json` and enter required values from marketplace
5. run `./testzapp`

With the dev version of zoom running, you should then be able to use the Install button on the app page to install the app. It will oath connect to https://mock.dev-coursera.org:9090/redirect and then the zoom app should open and your zapp will be installed.

