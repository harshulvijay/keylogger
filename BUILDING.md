# Building

A short guide on how to build the project.

## App

### Prerequisites

  1. Golang (must be in `PATH`)
  2. Mingw-w64 (on Windows)

### Building

#### On Windows

Run `app/scripts/build.ps1`. Output is generated in `app/.out`.

#### On Linux

TODO

## Server

An API written using Express and Firebase and deployed on Vercel.

### Prerequisites

  1. [pnpm](https://pnpm.io/installation)
  2. Vercel
  3. Firebase credentials

### Building

First, go to `server` and run `pnpm install`.
Then, run `pnpm build`. Output is stored in `server/api`.

### Deploying

Configure Vercel and Vercel CLI, then go to `server` and run:

```sh
vercel --prod
```

It should now deploy your API to Vercel.

Don't forget to populate the necessary environment variables. See `server/.dev.env` for what needs to be populated.
