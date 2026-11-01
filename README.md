# HyenaExternalUfaPlugin

A hello-world **external** plug-in for the United Federation of Agents. It is
modelled on `cats-ufa-plugin` in
[AI-evo1-spiral](https://github.com/unitedfederationofagents/AI-evo1-spiral)
and meets the contract in that repo's `docs/PluginSubApps.md`:

1. says hello over representable after every (re)connect (`connectLoop`);
2. serves `/ui-manifest.json` with ETag/304 and an icon at `/plugin-icon.svg`;
3. serves one page using relative URLs only, so it works under
   `/apps/<appId>/` and `/host/<id>/apps/<appId>/`.

As an external plug-in it runs in the **sandboxed** tier
(`docs/PluginTrust.md`): opaque origin, no cookies or storage, no host APIs.
The page therefore gets its data server-side (`withHello`) instead of fetching.

`third_party/representable` is a copy (MIT, stdlib only) of AI-evo1-spiral's
`representable` package, so this repo builds without it. Keep it in step with
upstream.

## Run

```sh
make test
make build
./hyena-ufa-plugin --port 8091 --lr-port 8082   # beside a running local-representative
```

Then, in local-representative's **plugins** tab, "Hyena" shows as `pending`;
click **approve hyena-ufa-plugin**. The page is framed at
`/apps/hyena-ufa-plugin/`. Port 8091 avoids the example plug-in's 8090.

Flags: `--port`, `--name` (app id; must match the representable client name),
`--lr-host`, `--lr-port`, `--dev-mode` (must match LR's).

LR does not launch or restart external plug-ins (lifecycle is out of scope in
version 1), so run this yourself.
