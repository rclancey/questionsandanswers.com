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

## API Calls

### List questsions & answers

```
GET /api/questions?<PARAMS>
```

Query parameters (URL encoded):

* `count` _integer_

  Maximum number of results to return per page

  Default: 10

  Max: 100

* `page` _integer_

  Page number in result set (where first page is zero)

  Default: 0

* `sort` _string_

  Sort order for results. Can be one of the following:

  * `questionDate` ascending date question was originally asked
  * `-questionDate` descending date question was originally asked
  * `answerDate` ascending date question was most recently answered
  * `-answerDate` descending date question was most recently answered

  Default: `-questionDate`

Response structure:

```json
{
    "meta": {
        "pageNum": 3,
        "pageSize": 100,
        "resultCount": 100,
        "sortOrder": "-questionDate",
        "prevPage": "/api/questions?count=100&page=2&sort=-questionDate",
        "thisPage": "/api/questions?count=100&page=3&sort=-questionDate",
        "nextPage": "/api/questions?count=100&page=4&sort=-questionDate",
        "totalPages": 10,
        "totalResults": 950
    },
    "results": [
        {
            "id": "12345678-90ab-cdef-1234-567890abcdef",
            "question": "Why does it always rain on me?",
            "questionDate": 1600622848050,
            "answer": "Because you lied when you were seventeen.",
            "answerDate": 1600622876227
        },
        {
            "id": "abcdef12-3456-7890-abcd-ef1234567890",
            "question": "Who will save your soul?",
            "questionDate": 1600622818050
        }
    ]
}
```

### Get question by ID

```
GET /api/question/12345678-90ab-cdef-1234-567890abcdef
```

Response structure:

```json
{
    "id": "12345678-90ab-cdef-1234-567890abcdef",
    "question": "Why does it always rain on me?",
    "questionDate": 1600622848050,
    "answer": "Because you lied when you were seventeen.",
    "answerDate": 1600622876227
}
```

### Create a question

```
POST /api/question
```

The request body should be the plain text of the question (`text/plain`).

Response structure:

```json
{
    "id": "abcdef12-3456-7890-abcd-ef1234567890",
    "question": "Who will save your soul?",
    "questionDate": 1600622818050
}
```

### Answer a question

```
PUT /api/question/12345678-90ab-cdef-1234-567890abcdef
```

The request body should be the plain text of the question (`text/plain`).

Response structure:

```json
{
    "id": "12345678-90ab-cdef-1234-567890abcdef",
    "question": "Why does it always rain on me?",
    "questionDate": 1600622848050,
    "answer": "Because you lied when you were seventeen.",
    "answerDate": 1600622876227
}
```
