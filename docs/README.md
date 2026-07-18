# EC v7 documentation

This directory contains the Hugo source for the EC v7 documentation site at
<https://ec.pbbgaming.com/docs/>. It uses the Hextra theme and organizes
content according to the Diátaxis framework.

## Preview locally

Hugo Extended and Go are required.

```sh
cd docs
hugo server --buildDrafts --disableFastRender
```

## Build

```sh
cd docs
hugo --gc --minify
```

The generated site is written to `docs/public/`.

## Deploy

Once the remote site has been provisioned, deploy it from the repository root:

```sh
ssh ec.pbbgaming.com /opt/ecv7/deploy-docs.sh
```
