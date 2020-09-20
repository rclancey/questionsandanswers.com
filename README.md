# questionsandanwsers.com

Coding exercise for NuORDER

This is a simple web application, with an API backend written in [go (golang)](https://golang.org/) and a frontend written in [React](https://reactjs.org/).  The frontend was bootstrapped using [create-react-app](https://github.com/facebook/create-react-app).

## Setting up a development environment

First, build the API server:

```sh
$ cd go
$ make
```

Next, install JS dependencies for the front end:

```sh
$ cd ../js
$ yarn install
```

Then, in a spearate terminal window, run the dev JS server. Note: the purpose of this to take advantage of hot reloading while developing. You will not be running a node server in production.

```sh
$ cd js
$ yarn run dev
```

By default, this will start the dev server on port 3000. If you have something else running on that port, it will choose a different port, of which you should make note and use in the next step.

Finally, start up the API server:

```sh
$ cd ../go
$ ./questionsandanswers -port 8080 -default-proxy http://localhost:3000/
```

Now, you can point your browser to http://localhost:8080/ and go to town.

## Production deployment

Production deployment is a pretty bare bones operation for this exercise.  In the project's root directory, just run `make`. This will build everything, and package it up in a gzipped tarfile named `build/questionsandanswers-0.0.1.tar.gz`. Copy that tarfile to your production server and unpack it. Then, in the unpacked directory, run the following command:

```sh
$ ./bin/startup.sh
```

