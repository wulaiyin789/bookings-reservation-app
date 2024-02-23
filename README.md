# Bookings and Reservations

Basic bookings and reservations application built by GoLang

- Built in Go version 1.15
- Uses the [chi router](github.com/go-chi/chi)
- Uses [alex edwards scs session management](github.com/alexedwards/scs)
- Uses [nosurf](github.com/justinas/nosurf)

## Run the application
- Start the application by using following code

```
./run.sh
```

## Testing
- Go to main directory and run the following code

```
go test 0v ./...
```

- Check the coverage testing by running the following code
```
go coverage
```

## Setup .tmpl emmet in VScode settings.json
```
{
  ...
  "emmet.preferences": {},
  "emmet.includeLanguages": {
      "golang": "tmpl"
  },
  "emmet.showSuggestionsAsSnippets": true,
  "files.associations": {
      "*html": "html",
      "*njk": "html",
      "*.tmpl": "html"
  },
  ...
}
```